// Package validation provides custom validation functions for Terraform schema
// attributes.
package validation

import (
	"net/mail"
	"strconv"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// IsEmailAddress is a SchemaValidateDiagFunc that tests if the provided value
// is of type string and matches the format of a valid email address.
func IsEmailAddress(i interface{}, p cty.Path) diag.Diagnostics {
	v, ok := i.(string)
	if !ok {
		return diag.Errorf("expected a string value, got %T", i)
	}

	if v == "" {
		return diag.Errorf("expected a non-empty email address")
	}

	a, err := mail.ParseAddress(v)
	if err != nil || a.Address != v {
		return diag.Errorf("expected a valid email address, got %q", v)
	}

	return nil
}

// StringIsNumerical is a SchemaValidateDiagFunc that tests if the provided
// value is of type string and represents a numerical value.
func StringIsNumerical(i interface{}, p cty.Path) diag.Diagnostics {
	v, ok := i.(string)
	if !ok {
		return diag.Errorf("expected a string value, got %T", i)
	}

	if v == "" {
		return diag.Errorf("expected a non-empty string representation of a number")
	}

	if _, err := strconv.Atoi(v); err != nil {
		return diag.Errorf("expected a string representation of a number, got %q", v)
	}

	return nil
}

// Int32InSlice returns a SchemaValidateDiagFunc that tests if the provided
// value is an integer that matches an element of the valid slice.
func Int32InSlice(valid []int32) schema.SchemaValidateDiagFunc {
	return func(i interface{}, p cty.Path) diag.Diagnostics {
		v, ok := i.(int)
		if !ok {
			return diag.Errorf("expected an int value, got %T", i)
		}

		// Cast the int to int64 to avoid issues with integer wrapping.
		for _, validInt := range valid {
			if int64(v) == int64(validInt) {
				return nil
			}
		}

		return diag.Errorf("expected value to be one of %v", valid)
	}
}
