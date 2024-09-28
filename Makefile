OS = $(shell uname | tr A-Z a-z)
export PATH := $(abspath bin/):${PATH}

# Build Variables
export CGO_ENABLED ?= 0
export GOOS = $(shell go env GOOS)
ifeq (${VERBOSE}, 1)
ifeq ($(filter -v,${GOARGS}),)
	GOARGS += -v
endif
TEST_FORMAT = short-VERBOSE
endif

