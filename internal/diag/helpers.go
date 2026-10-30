// Package diag provides helpers for converting errors into Terraform
// Diagnostics.
package diag

import (
	"maps"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/StatusCakeDev/statuscake-go"
)

// WrapErr will convert an error into a Diagnostics. Each Diagnostic entry will
// have the summary line prefixed with a contextual message.
func WrapErr(message string, err error) diag.Diagnostics {
	if err == nil {
		return nil
	}
	return diagnostics(message, err)
}

func diagnostics(message string, err error) diag.Diagnostics {
	errs := statuscake.Errors(err)
	if len(errs) == 0 {
		return fromErr(message, err)
	}
	return violations(message, err, errs)
}

func fromErr(message string, err error) diag.Diagnostics {
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Error,
			Summary:  message + ": " + err.Error(),
		},
	}
}

func violations(message string, err error, errs map[string][]string) diag.Diagnostics {
	var diags diag.Diagnostics

	// Report the fields in a stable order, since the violations are held in a
	// map and would otherwise come out differently on each run.
	for _, field := range slices.Sorted(maps.Keys(errs)) {
		violations := errs[field]
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  message + ": " + err.Error() + ": " + field + " contains violations",

			// TODO: Use AttributePath to indicate validation errors.
			Detail: strings.Join(violations, "; "),
		})
	}
	return diags
}
