// Package adapter implements the Terraform provider's operations in terms of
// the StatusCake SDK. It is the only package that imports the SDK: specs go in,
// domain models come out, and SDK errors are translated into the sentinel
// errors of the domain package.
package adapters
