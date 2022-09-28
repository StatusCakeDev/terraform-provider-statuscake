package diag

import (
	"strings"

	"github.com/StatusCakeDev/statuscake-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

// FromErr will convert an error into a Diagnostics. Each Diagnostic entry will
// have the summary line prefixed with a contextual message.
func FromErr(message string, err error) diag.Diagnostics {
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
	for field, violations := range errs {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  message + ": " + err.Error() + ": " + field + " contains violations",

			// TODO: Use AttributePath to indicate validation errors.
			Detail: strings.Join(violations, "; "),
		})
	}
	return diags
}
