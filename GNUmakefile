BINARY_NAME := terraform-provider-statuscake

default: build

# Build provider binary
.PHONY: build
build:
	go build -o $(BINARY_NAME)

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m
