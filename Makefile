SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
.DEFAULT_GOAL := all
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules
GO ?= go
TESTARGS ?=
BINARY_NAME := terraform-provider-statuscake

ifeq ($(origin .RECIPEPREFIX), undefined)
  $(error This Make does not support .RECIPEPREFIX. Please use GNU Make 4.0 or later)
endif
.RECIPEPREFIX =

.PHONY: all
all: build testacc

.PHONY: help
help:
	@echo "Usage: make <target>"
	@echo
	@echo "Targets:"
	@echo "  all     Build and test the project (default)"
	@echo "  build   Build the project"
	@echo "  testacc Run acceptance tests"
	@echo "  clean   Clean the project"

.PHONY: build
build: $(BINARY_NAME)
	@echo "done"

$(BINARY_NAME):
	@echo "building provider"
	@go build -o $@

.PHONY: docs
docs:
	@go generate ./...

.PHONY: testacc
testacc:
	TF_ACC=1 $(GO) test ./... -v $(TESTARGS) -timeout 120m

.PHONY: clean
clean:
	@rm -f $(BINARY_NAME)
