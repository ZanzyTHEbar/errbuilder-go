package errbuilder

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type ErrDetails struct {
	Errors ErrorMap `json:"errors"`
}

type ErrBuilder struct {
	Code    ErrCode    `json:"code"`
	Msg     string     `json:"message"`
	Cause   error      `json:"Cause"`
	Label   string     `json:"label"`
	Details ErrDetails `json:"details"`
}

// NewErrDetails is a constructor for ErrDetails
func NewErrDetails(errors ErrorMap) ErrDetails {
	return ErrDetails{Errors: errors}
}

// NewErrBuilder is a constructor for ErrBuilder
func NewErrBuilder() *ErrBuilder {
	return &ErrBuilder{}
}

// MarshalJSON implements the json.Marshaler interface.
func (builder *ErrBuilder) MarshalJSON() ([]byte, error) {
	// use json.Marshal to convert the error message to a JSON byte slice
	byteBuffer, err := json.Marshal(map[string]interface{}{
		"code":    builder.Code,
		"message": builder.Msg,
		"Cause":   builder.Cause.Error(),
		"label":   builder.Label,
		"details": builder.Details,
	})
	if err != nil {
		return nil, err
	}

	return byteBuffer, nil
}

// Error is a method to return an error, this is an implementation of the error interface.
func (builder *ErrBuilder) Error() string {

	// validate the error instance, if it is nil, return nil
	if builder.Code == 0 || builder.Msg == "" {
		builder.Code = CodeInternal
		builder.Msg = "Internal Server Error"
	}

	// if the Cause is nil, set the Cause to the error message
	if builder.Cause == nil {
		builder.Cause = errors.New(builder.Msg)
	}

	// if the label is empty, set the label to the String of the error code
	if builder.Label == "" {
		builder.Label = builder.Code.String()
	}

	// convert the builder instance to a formatted error message and return it
	return fmt.Sprintf("code: %d, label: %s, message: %s, Cause: %s, details: %v",
		builder.Code, builder.Label, builder.Msg, builder.Cause.Error(), builder.Details)
}

// WithCode is a method to set the error code.
func (builder *ErrBuilder) WithCode(code ErrCode) *ErrBuilder {
	builder.Code = code
	return builder
}

// WithLabel is a method to set the error label.
func (builder *ErrBuilder) WithLabel(label string) *ErrBuilder {
	builder.Label = label
	return builder
}

// WithMsg is a method to set the error message.
func (builder *ErrBuilder) WithMsg(msg string) *ErrBuilder {
	builder.Msg = msg
	return builder
}

// WithCause is a method to set the error Cause.
func (builder *ErrBuilder) WithCause(Cause error) *ErrBuilder {
	builder.Cause = Cause
	return builder
}

// WithDetails is a method to set the error details.
func (builder *ErrBuilder) WithDetails(details ErrDetails) *ErrBuilder {
	builder.Details = details
	return builder
}

// Code returns the error's status code.
func (err *ErrBuilder) ErrCode() ErrCode {
	return err.Code
}

// Unwrap allows [errors.Is] and [errors.As] access to the underlying error.
func (err *ErrBuilder) Unwrap() error {
	return err.Cause
}

// UnWrap is a method to return the error details as a map of errors.
func (err *ErrDetails) UnWrap() (ErrorMap, error) {
	if err.Errors == nil {
		return nil, errors.New("no error details found")
	}
	return err.Errors, nil
}

// asErrorBuilder converts the given error to an ErrBuilder.
func asErrorBuilder(err error) (*ErrBuilder, bool) {
	if err == nil {
		return nil, false
	}

	var errBuilder *ErrBuilder

	if ok := errors.As(err, &errBuilder); !ok {
		return nil, false
	}

	return errBuilder, true
}

// errorf calls fmt.Errorf with the supplied template and arguments, then wraps
// the resulting error.
func errorf(c ErrCode, template string, args ...any) *ErrBuilder {
	return NewErrBuilder().
		WithCause(fmt.Errorf(template, args...)).
		WithCode(c)
}

// wrapIfUncoded ensures that all errors are wrapped. It leaves already-wrapped
// errors unchanged, uses wrapIfContextError to apply codes to context.Canceled
// and context.DeadlineExceeded, and falls back to wrapping other errors with
// CodeUnknown.
func WrapIfUncoded(err error) error {
	if err == nil {
		return nil
	}
	maybeCodedErr := WrapIfContextError(err)
	if _, ok := asErrorBuilder(maybeCodedErr); ok {
		return maybeCodedErr
	}
	return NewErrBuilder().WithCause(maybeCodedErr).WithCode(CodeUnknown)
}

// wrapIfContextError applies CodeCanceled or CodeDeadlineExceeded to Go's
// context.Canceled and context.DeadlineExceeded errors, but only if they
// haven't already been wrapped.
func WrapIfContextError(err error) error {
	if err == nil {
		return nil
	}
	if _, ok := asErrorBuilder(err); ok {
		return err
	}
	if errors.Is(err, context.Canceled) {
		return NewErrBuilder().WithCause(err).WithCode(CodeCanceled)
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return NewErrBuilder().WithCause(err).WithCode(CodeDeadlineExceeded)
	}
	// Ick, some dial errors can be returned as os.ErrDeadlineExceeded
	// instead of context.DeadlineExceeded :(
	// https://github.com/golang/go/issues/64449
	if errors.Is(err, os.ErrDeadlineExceeded) {
		return NewErrBuilder().WithCause(err).WithCode(CodeDeadlineExceeded)
	}
	return err
}

// wrapIfContextDone wraps errors with CodeCanceled or CodeDeadlineExceeded
// if the context is done. It leaves already-wrapped errors unchanged.
func WrapIfContextDone(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}
	err = WrapIfContextError(err)
	if _, ok := asErrorBuilder(err); ok {
		return err
	}
	ctxErr := ctx.Err()
	if errors.Is(ctxErr, context.Canceled) {
		return NewErrBuilder().WithCause(err).WithCode(CodeCanceled)
	} else if errors.Is(ctxErr, context.DeadlineExceeded) {
		return NewErrBuilder().WithCause(err).WithCode(CodeDeadlineExceeded)
	}
	return err
}

// wrapIfLikelyH2CNotConfiguredError adds a wrapping error that has a message
// telling the caller that they likely need to use h2c but are using a raw http.Client{}.
//
// This happens when running a gRPC-only server.
// This is fragile and may break over time, and this should be considered a best-effort.
func WrapIfLikelyH2CNotConfiguredError(request *http.Request, err error) error {
	if err == nil {
		return nil
	}
	if _, ok := asErrorBuilder(err); ok {
		return err
	}
	if url := request.URL; url != nil && url.Scheme != "http" {
		// If the scheme is not http, we definitely do not have an h2c error, so just return.
		return err
	}
	// net/http code has been investigated and there is no typing of any of these errors
	// they are all created with fmt.Errorf
	// grpc-go returns the first error 2/3-3/4 of the time, and the second error 1/4-1/3 of the time
	if errString := err.Error(); strings.HasPrefix(errString, `Post "`) &&
		(strings.Contains(errString, `net/http: HTTP/1.x transport connection broken: malformed HTTP response`) ||
			strings.HasSuffix(errString, `write: broken pipe`)) {
		return fmt.Errorf("possible h2c configuration issue when talking to gRPC server, see: %w", err)
	}
	return err
}

// wrapIfLikelyWithGRPCNotUsedError adds a wrapping error that has a message
// telling the caller that they likely forgot to use WithGRPC().
//
// This happens when running a gRPC-only server.
// This is fragile and may break over time, and this should be considered a best-effort.
func WrapIfLikelyWithGRPCNotUsedError(err error) error {
	if err == nil {
		return nil
	}
	if _, ok := asErrorBuilder(err); ok {
		return err
	}
	// golang.org/x/net code has been investigated and there is no typing of this error
	// it is created with fmt.Errorf
	// http2/transport.go:573:	return nil, fmt.Errorf("http2: Transport: cannot retry err [%v] after Request.Body was written; define Request.GetBody to avoid this error", err)
	if errString := err.Error(); strings.HasPrefix(errString, `Post "`) &&
		strings.Contains(errString, `http2: Transport: cannot retry err`) &&
		strings.HasSuffix(errString, `after Request.Body was written; define Request.GetBody to avoid this error`) {
		return fmt.Errorf("possible missing WithGPRC() client option when talking to gRPC server, see: %w", err)
	}
	return err
}

// HTTP/2 has its own set of error codes, which it sends in RST_STREAM frames.
// When the server sends one of these errors, we should map it back into our
// RPC error codes following
// https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md#http2-transport-mapping.
//
// This would be vastly simpler if we were using x/net/http2 directly, since
// the StreamError type is exported. When x/net/http2 gets vendored into
// net/http, though, all these types become unexported...so we're left with
// string munging.
func WrapIfRSTError(err error) error {
	const (
		streamErrPrefix = "stream error: "
		fromPeerSuffix  = "; received from peer"
	)
	if err == nil {
		return nil
	}
	if _, ok := asErrorBuilder(err); ok {
		return err
	}
	if urlErr := new(url.Error); errors.As(err, &urlErr) {
		// If we get an RST_STREAM error from http.Client.Do, it's wrapped in a
		// *url.Error.
		err = urlErr.Unwrap()
	}
	msg := err.Error()
	if !strings.HasPrefix(msg, streamErrPrefix) {
		return err
	}
	if !strings.HasSuffix(msg, fromPeerSuffix) {
		return err
	}
	msg = strings.TrimSuffix(msg, fromPeerSuffix)
	i := strings.LastIndex(msg, ";")
	if i < 0 || i >= len(msg)-1 {
		return err
	}
	msg = msg[i+1:]
	msg = strings.TrimSpace(msg)
	switch msg {
	case "NO_ERROR", "PROTOCOL_ERROR", "INTERNAL_ERROR", "FLOW_CONTROL_ERROR",
		"SETTINGS_TIMEOUT", "FRAME_SIZE_ERROR", "COMPRESSION_ERROR":
		return InternalServerErr(fmt.Errorf("http2 error: %w", err))
	case "REFUSED_STREAM":
		return NewErrBuilder().WithCause(err).WithCode(CodeUnavailable)
	case "CANCEL":
		return NewErrBuilder().WithCause(err).WithCode(CodeCanceled)
	case "ENHANCE_YOUR_CALM":
		return NewErrBuilder().WithCause(fmt.Errorf("bandwidth exhausted: %w", err)).WithCode(CodeResourceExhausted)
	case "INADEQUATE_SECURITY":
		return NewErrBuilder().WithCause(fmt.Errorf("transport protocol insecure: %w", err)).WithCode(CodePermissionDenied)
	default:
		return err
	}
}

// wrapIfMaxBytesError wraps errors returned reading from a http.MaxBytesHandler
// whose limit has been exceeded.
func WrapIfMaxBytesError(err error, tmpl string, args ...any) error {
	if err == nil {
		return nil
	}
	if _, ok := asErrorBuilder(err); ok {
		return err
	}
	var maxBytesErr *http.MaxBytesError
	if ok := errors.As(err, &maxBytesErr); !ok {
		return err
	}
	prefix := fmt.Sprintf(tmpl, args...)
	return errorf(CodeResourceExhausted, "%s: exceeded %d byte http.MaxBytesReader limit", prefix, maxBytesErr.Limit)
}
