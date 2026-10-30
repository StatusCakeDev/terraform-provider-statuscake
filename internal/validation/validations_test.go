package validation_test

import (
	"reflect"
	"testing"

	"github.com/StatusCakeDev/terraform-provider-statuscake/v2/internal/validation"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

var (
	path         = cty.IndexStringPath("test")
	validNumbers = []int32{1, 2, 3}
)

func TestIsEmailAddress(t *testing.T) {
	t.Run("returns no errors when the given value is a valid email address", func(t *testing.T) {
		diags := validation.IsEmailAddress("picard@starfleet.com", path)
		if len(diags) != 0 {
			t.Error("expected no errors but errors were returned")
		}
	})

	t.Run("returns an error when the value is not of type string", func(t *testing.T) {
		expected := []string{"expected a string value, got int"}

		diags := validation.IsEmailAddress(1701, path)
		if len(diags) == 0 {
			t.Error("expected errors but no errors were returned")
		}

		if !reflect.DeepEqual(collect(diags), expected) {
			t.Error("unexpected error message")
		}
	})

	t.Run("returns an error when the value is an empty string", func(t *testing.T) {
		expected := []string{"expected a non-empty email address"}

		diags := validation.IsEmailAddress("", path)
		if len(diags) == 0 {
			t.Error("expected errors but no errors were returned")
		}

		if !reflect.DeepEqual(collect(diags), expected) {
			t.Error("unexpected error message")
		}
	})

	t.Run("returns an error when the value is not a valid email address", func(t *testing.T) {
		expected := []string{`expected a valid email address, got "enterprise"`}

		diags := validation.IsEmailAddress("enterprise", path)
		if len(diags) == 0 {
			t.Error("expected errors but no errors were returned")
		}

		if !reflect.DeepEqual(collect(diags), expected) {
			t.Error("unexpected error message")
		}
	})
}

func TestStringIsNumerical(t *testing.T) {
	t.Run("returns no errors when the given value is a string representation of a numerical value", func(t *testing.T) {
		diags := validation.StringIsNumerical("123", path)
		if len(diags) != 0 {
			t.Error("expected no errors but errors were returned")
		}
	})

	t.Run("returns an error when the value is not of type string", func(t *testing.T) {
		expected := []string{"expected a string value, got int"}

		diags := validation.StringIsNumerical(1701, path)
		if len(diags) == 0 {
			t.Error("expected errors but no errors were returned")
		}

		if !reflect.DeepEqual(collect(diags), expected) {
			t.Error("unexpected error message")
		}
	})

	t.Run("returns an error when the value is an empty string", func(t *testing.T) {
		expected := []string{"expected a non-empty string representation of a number"}

		diags := validation.StringIsNumerical("", path)
		if len(diags) == 0 {
			t.Error("expected errors but no errors were returned")
		}

		if !reflect.DeepEqual(collect(diags), expected) {
			t.Error("unexpected error message")
		}
	})

	t.Run("returns an error when the value is not a string representation of a numerical value", func(t *testing.T) {
		expected := []string{`expected a string representation of a number, got "enterprise"`}

		diags := validation.StringIsNumerical("enterprise", path)
		if len(diags) == 0 {
			t.Error("expected errors but no errors were returned")
		}

		if !reflect.DeepEqual(collect(diags), expected) {
			t.Error("unexpected error message")
		}
	})
}

func TestInt32InSlice(t *testing.T) {
	t.Run("returns no errors when the given value is contained within the validation slice", func(t *testing.T) {
		diags := validation.Int32InSlice(validNumbers)(2, path)
		if len(diags) != 0 {
			t.Error("expected no errors but errors were returned")
		}
	})

	t.Run("returns an error when the value is not of type int32", func(t *testing.T) {
		expected := []string{"expected an int value, got string"}

		diags := validation.Int32InSlice(validNumbers)("1701", path)
		if len(diags) == 0 {
			t.Error("expected errors but no errors were returned")
		}

		if !reflect.DeepEqual(collect(diags), expected) {
			t.Error("unexpected error message")
		}
	})

	t.Run("returns an error when the value is not contained within the validation slice", func(t *testing.T) {
		expected := []string{"expected value to be one of [1 2 3]"}

		diags := validation.Int32InSlice(validNumbers)(4, path)
		if len(diags) == 0 {
			t.Error("expected errors but no errors were returned")
		}

		if !reflect.DeepEqual(collect(diags), expected) {
			t.Error("unexpected error message")
		}
	})
}

func collect(diags diag.Diagnostics) []string {
	strs := make([]string, len(diags))
	for i, diag := range diags {
		strs[i] = diag.Summary
	}
	return strs
}
