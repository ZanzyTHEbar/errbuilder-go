package errbuilder

import "errors"

type ErrDetails struct {
	Errors ErrorMap `json:"errors"`
}

// NewErrDetails is a constructor for ErrDetails
func NewErrDetails(errors ErrorMap) ErrDetails {
	return ErrDetails{Errors: errors}
}

// UnWrap is a method to return the error details as a map of errors.
func (err *ErrDetails) UnWrap() (ErrorMap, error) {
	if err.Errors == nil {
		return nil, errors.New("no error details found")
	}
	return err.Errors, nil
}
