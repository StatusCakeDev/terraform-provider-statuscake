package diag_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/StatusCakeDev/statuscake-go"

	intdiag "github.com/StatusCakeDev/terraform-provider-statuscake/v2/internal/diag"
)

func TestFromErr(t *testing.T) {
	t.Parallel()

	t.Run("no error gives no diagnostics", func(t *testing.T) {
		t.Parallel()
		if got := intdiag.WrapErr("failed", nil); got != nil {
			t.Errorf("got %v, want nil", got)
		}
	})

	t.Run("a plain error is prefixed with the context", func(t *testing.T) {
		t.Parallel()

		got := intdiag.WrapErr("failed to create check", errors.New("boom"))
		want := diag.Diagnostics{{Severity: diag.Error, Summary: "failed to create check: boom"}}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("(-want +got):\n%s", diff)
		}
	})

	t.Run("an API error without violations is a single diagnostic", func(t *testing.T) {
		t.Parallel()

		got := intdiag.WrapErr("failed", statuscake.APIError{Status: 500, Message: "Internal error"})
		if len(got) != 1 || got[0].Severity != diag.Error || got[0].Detail != "" {
			t.Errorf("unexpected diagnostics: %+v", got)
		}
	})

	t.Run("each field with violations has its own diagnostic, in order of field", func(t *testing.T) {
		t.Parallel()

		err := statuscake.APIError{
			Status:  400,
			Message: "Bad request",
			Errors: map[string][]string{
				"name":         {"is too long", "is not unique"},
				"check_period": {"is invalid"},
				"address":      {"is required"},
				"tags":         {"has too many items"},
				"user_agent":   {"is too long"},
			},
		}

		// The violations are held in a map, so an unsorted result would differ
		// between runs. Repeat to make that likely to be noticed.
		for range 20 {
			got := intdiag.WrapErr("failed", err)

			var fields []string
			details := map[string]string{}
			for _, d := range got {
				if d.Severity != diag.Error {
					t.Errorf("severity = %v, want error", d.Severity)
				}
				field := strings.TrimSuffix(d.Summary[strings.LastIndex(d.Summary, ": ")+2:], " contains violations")
				fields = append(fields, field)
				details[field] = d.Detail
			}

			wantFields := []string{"address", "check_period", "name", "tags", "user_agent"}
			if diff := cmp.Diff(wantFields, fields); diff != "" {
				t.Fatalf("fields differ (-want +got):\n%s", diff)
			}

			// The violations of a field keep the order the API gave them.
			if got, want := details["name"], "is too long; is not unique"; got != want {
				t.Errorf("name details = %q, want %q", got, want)
			}
		}
	})
}
