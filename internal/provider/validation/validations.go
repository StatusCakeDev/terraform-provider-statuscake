package validation

import (
	"fmt"
	"net/mail"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// IsEmailAddress is a SchemaValidateFunc that tests if the provided value is
// of type string and matches the format of a valid email address.
func IsEmailAddress(i interface{}, k string) ([]string, []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	if v == "" {
		return nil, []error{fmt.Errorf("expected %q email address to not be empty, got %q", k, i)}
	}

	if _, err := mail.ParseAddress(v); err != nil {
		return nil, []error{fmt.Errorf("expected %q to be a valid email address, got %q: %+v", k, v, err)}
	}

	return nil, nil
}

// StringIsNumerical is a SchemaValidateFunc that tests if the provided value
// is of type string and represents a numerical value.
func StringIsNumerical(i interface{}, k string) ([]string, []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	if v == "" {
		return nil, []error{fmt.Errorf("expected %q number to not be empty, got %q", k, i)}
	}

	if _, err := strconv.Atoi(v); err != nil {
		return nil, []error{fmt.Errorf("expected %q to be a valid number, got %q: %+v", k, v, err)}
	}

	return nil, nil
}

// Int32InSlice returns a SchemaValidateFunc that tests if the provided value
// is of type int32 and matches the value of an element in the valid slice.
func Int32InSlice(valid []int32) schema.SchemaValidateFunc {
	return func(i interface{}, k string) ([]string, []error) {
		v, ok := i.(int)
		if !ok {
			return nil, []error{fmt.Errorf("expected type of %q to be integer", k)}
		}

		for _, validInt := range valid {
			if int32(v) == validInt {
				return nil, nil
			}
		}

		return nil, []error{fmt.Errorf("expected %q to be one of %+v, got %d", k, valid, v)}
	}
}
