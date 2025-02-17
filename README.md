# ErrBuilder - A simple error builder for Go

A lightweight custom error builder library for Go.

I originally made this for my personal projects. It uses the `Builder` pattern to create a custom error object with a custom message and a custom error code. As well as convenience methods and an `ErrMap` type. 

`ErrMap` is a string indexed map of `error` types.

The purpose of this library was to make a standard way of creating errors in my projects. 

I wanted to create a robust and standardized meaning for certain error codes, I largely borrowed the gRPC error code specification as inspiration.

This is a general purpose library that can be used in any Go project.

## Features

- **Error Builder**: Custom Error Builder for creating structured error messages to your requirements.
- **Error Codes**: A rich set of error codes, with an interface, for defining the type of error. Follows the gRPC error code specification.
- **Error Map**: An optional and dynamic map to contain error messages for complex control flows, perhaps even deferred error handling.
- **Error Details**: Custom ErrDetails type that allows providing extra data to be JSON (or other type) formatted into the error message.
- **Builtin Custom Error wrappers**: (errors.go)[/errors.go] contains 5 builtin error functions that demonstrate how to use the errbuilder, and are provided for common error usage requirements.

Works very well with my (assert-lib)[https://github.com/ZanzyTHEbar/assert-lib] library.

## Installation

```bash
go get github.com/ZanzyTHEbar/errbuilder-go
```

## Usage

```go
package main

import (
    "context"
    "github.com/ZanzyTHEbar/errorbuilder"
)

func main() {
    // Create new ErrorMap to hold our error messages
    var errs errsx.ErrorMap

    if len(user.Username) < 4 {
	    errs.Set("username", "Username must be at least 4 characters")
    }
    if len(user.Password) < 8 {
	    errs.Set("password", "Password must be at least 8 characters")
    }
    if len(user.Email) == 0 {
	     errs.Set("email", "Email is required")
    }

    // Create custom error to handle our error messages
    customError := errbuilder.NewErrBuilder().
                    WithCode(errbuilder.CodeInvalidArgument).
		    WithMsg("Bad Request").
		    WithDetails(errbuilder.NewErrDetails(errs))

    // Check if there were any errors
    if errs != nil {
	    // Return the errors as a JSON response
	    return customError
    }
}
```
