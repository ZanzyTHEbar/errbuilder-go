# ErrBuilder - A simple error builder for Go

This is a simple error builder that I made for my personal projects. It uses the `Builder` pattern to create a custom error object with a custom message and a custom error code. As well as convience methods and an `ErrMap` type with is a string indexed map of `error` types.

The purpose of this library was to make a standard way of creating errors in my projects that has a standardized meaning for certain error codes, I largely borrowed the gRPC error code specification.

I wanted to have a way to create errors with a custom metadata. I also wanted to have a way to create a map of errors that I could use to store errors for building validation logic.

This is a general purpose library that can be used in any Go project. Please check out the [`examples`](/examples) directory for more information.

