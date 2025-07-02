package errbuilder

import (
	"encoding/json"
	"errors"
	"fmt"
)

type ErrBuilder struct {
	Code    ErrCode    `json:"code"`
	Msg     string     `json:"message"`
	Cause   error      `json:"Cause"`
	Label   string     `json:"label"`
	Details ErrDetails `json:"details"`
}

// NewErrBuilder is a constructor for ErrBuilder
func New() *ErrBuilder {
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
