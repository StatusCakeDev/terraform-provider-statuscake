package adapters_test

import (
	"strings"
	"testing"
	"time"

	"github.com/StatusCakeDev/terraform-provider-statuscake/v2/internal/adapters"
)

func newClient(cfg adapters.Config) (*adapters.Client, error) {
	return adapters.New(cfg)
}

func valid() adapters.Config {
	return adapters.Config{
		Token:      "tok",
		RPS:        4,
		Retries:    3,
		MinBackoff: time.Second,
		MaxBackoff: 30 * time.Second,
	}
}

func TestNewAcceptsAValidConfig(t *testing.T) {
	t.Parallel()

	client, err := newClient(valid())
	if err != nil || client == nil {
		t.Fatalf("New() = %v, %v; want a client", client, err)
	}
}

func TestNewRejectsInvalidConfig(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		mutate func(*adapters.Config)
		want   string
	}{
		"no token": {
			mutate: func(c *adapters.Config) {
				c.Token = ""
			},
			want: "API token is required",
		},
		"zero rate": {
			mutate: func(c *adapters.Config) {
				c.RPS = 0
			},
			want: "rps must be greater than zero",
		},
		"negative retries": {
			mutate: func(c *adapters.Config) {
				c.Retries = -1
			},
			want: "retries must not be negative",
		},
		"negative min backoff": {
			mutate: func(c *adapters.Config) {
				c.MinBackoff = -1
			},
			want: "min backoff must not be negative",
		},
		"max below min": {
			mutate: func(c *adapters.Config) {
				c.MaxBackoff = time.Millisecond
			},
			want: "must not be less than min backoff",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cfg := valid()
			tt.mutate(&cfg)

			_, err := newClient(cfg)
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("New() error = %v, want one containing %q", err, tt.want)
			}
		})
	}
}

func TestNewReportsEveryProblemTogether(t *testing.T) {
	t.Parallel()

	_, err := newClient(adapters.Config{Retries: -1, MinBackoff: -1})
	if err == nil {
		t.Fatal("New() succeeded")
	}
	for _, want := range []string{"API token", "rps", "retries", "min backoff"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q does not mention %q", err, want)
		}
	}
}
