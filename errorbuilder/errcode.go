package errsx

import (
	"fmt"
	"strconv"
	"strings"
)

type ErrCode uint32

const (

	// CodeCanceled indicates that the operation was canceled, typically by the
	// caller.
	CodeCanceled ErrCode = 1

	// CodeUnknown indicates that the operation failed for an unknown reason.
	CodeUnknown ErrCode = 2

	// CodeInvalidArgument indicates that client supplied an invalid argument.
	CodeInvalidArgument ErrCode = 3

	// CodeDeadlineExceeded indicates that deadline expired before the operation
	// could complete.
	CodeDeadlineExceeded ErrCode = 4

	// CodeNotFound indicates that some requested entity (for example, a file or
	// directory) was not found.
	CodeNotFound ErrCode = 5

	// CodeAlreadyExists indicates that client attempted to create an entity (for
	// example, a file or directory) that already exists.
	CodeAlreadyExists ErrCode = 6

	// CodePermissionDenied indicates that the caller doesn't have permission to
	// execute the specified operation.
	CodePermissionDenied ErrCode = 7

	// CodeResourceExhausted indicates that some resource has been exhausted. For
	// example, a per-user quota may be exhausted or the entire file system may
	// be full.
	CodeResourceExhausted ErrCode = 8

	// CodeFailedPrecondition indicates that the system is not in a state
	// required for the operation's execution.
	CodeFailedPrecondition ErrCode = 9

	// CodeAborted indicates that operation was aborted by the system, usually
	// because of a concurrency issue such as a sequencer check failure or
	// transaction abort.
	CodeAborted ErrCode = 10

	// CodeOutOfRange indicates that the operation was attempted past the valid
	// range (for example, seeking past end-of-file).
	CodeOutOfRange ErrCode = 11

	// CodeUnimplemented indicates that the operation isn't implemented,
	// supported, or enabled in this service.
	CodeUnimplemented ErrCode = 12

	// CodeInternal indicates that some invariants expected by the underlying
	// system have been broken. This code is reserved for serious errors.
	CodeInternal ErrCode = 13

	// CodeUnavailable indicates that the service is currently unavailable. This
	// is usually temporary, so clients can back off and retry idempotent
	// operations.
	CodeUnavailable ErrCode = 14

	// CodeDataLoss indicates that the operation has resulted in unrecoverable
	// data loss or corruption.
	CodeDataLoss ErrCode = 15

	// CodeUnauthenticated indicates that the request does not have valid
	// authentication credentials for the operation.
	CodeUnauthenticated ErrCode = 16

	minCode = CodeCanceled
	maxCode = CodeUnauthenticated
)

func (c ErrCode) String() string {
	switch c {
	case CodeCanceled:
		return "canceled"
	case CodeUnknown:
		return "unknown"
	case CodeInvalidArgument:
		return "invalid_argument"
	case CodeDeadlineExceeded:
		return "deadline_exceeded"
	case CodeNotFound:
		return "not_found"
	case CodeAlreadyExists:
		return "already_exists"
	case CodePermissionDenied:
		return "permission_denied"
	case CodeResourceExhausted:
		return "resource_exhausted"
	case CodeFailedPrecondition:
		return "failed_precondition"
	case CodeAborted:
		return "aborted"
	case CodeOutOfRange:
		return "out_of_range"
	case CodeUnimplemented:
		return "unimplemented"
	case CodeInternal:
		return "internal"
	case CodeUnavailable:
		return "unavailable"
	case CodeDataLoss:
		return "data_loss"
	case CodeUnauthenticated:
		return "unauthenticated"
	}
	return fmt.Sprintf("code_%d", c)
}

// MarshalText implements [encoding.TextMarshaler].
func (c ErrCode) MarshalText() ([]byte, error) {
	return []byte(c.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler].
func (c *ErrCode) UnmarshalText(data []byte) error {
	dataStr := string(data)
	switch dataStr {
	case "canceled":
		*c = CodeCanceled
		return nil
	case "unknown":
		*c = CodeUnknown
		return nil
	case "invalid_argument":
		*c = CodeInvalidArgument
		return nil
	case "deadline_exceeded":
		*c = CodeDeadlineExceeded
		return nil
	case "not_found":
		*c = CodeNotFound
		return nil
	case "already_exists":
		*c = CodeAlreadyExists
		return nil
	case "permission_denied":
		*c = CodePermissionDenied
		return nil
	case "resource_exhausted":
		*c = CodeResourceExhausted
		return nil
	case "failed_precondition":
		*c = CodeFailedPrecondition
		return nil
	case "aborted":
		*c = CodeAborted
		return nil
	case "out_of_range":
		*c = CodeOutOfRange
		return nil
	case "unimplemented":
		*c = CodeUnimplemented
		return nil
	case "internal":
		*c = CodeInternal
		return nil
	case "unavailable":
		*c = CodeUnavailable
		return nil
	case "data_loss":
		*c = CodeDataLoss
		return nil
	case "unauthenticated":
		*c = CodeUnauthenticated
		return nil
	}
	// Ensure that non-canonical codes round-trip through MarshalText and
	// UnmarshalText.
	if strings.HasPrefix(dataStr, "code_") {
		dataStr = strings.TrimPrefix(dataStr, "code_")
		code, err := strconv.ParseInt(dataStr, 10 /* base */, 64 /* bitsize */)
		if err == nil && (code < int64(minCode) || code > int64(maxCode)) {
			*c = ErrCode(code)
			return nil
		}
	}
	return fmt.Errorf("invalid code %q", dataStr)
}

// CodeOf returns the error's status code if it is or wraps an [*ErrBuilder] and
// [CodeUnknown] otherwise.
func CodeOf(err error) ErrCode {
	if errBuilder, ok := asErrorBuilder(err); ok {
		return errBuilder.ErrCode()
	}
	return CodeUnknown
}
