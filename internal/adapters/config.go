package adapters

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/time/rate"

	"github.com/StatusCakeDev/statuscake-go"
	"github.com/StatusCakeDev/statuscake-go/backoff"
	"github.com/StatusCakeDev/statuscake-go/credentials"
	"github.com/StatusCakeDev/statuscake-go/throttle"
)

// Config is everything the adapter needs to reach StatusCake. It is owned by
// this package, so the composition root can build a [Client] without knowing
// which SDK is behind it.
type Config struct {
	Token string

	// RPS caps requests per second, Retries is how many times a failed request
	// is tried again, and the backoffs bound the wait between tries.
	RPS        int
	Retries    int
	MinBackoff time.Duration
	MaxBackoff time.Duration

	// Endpoint replaces the StatusCake API host. Release builds never set it,
	// since the provider offers no such setting (see the dev build tag in
	// internal/provider); it exists so that tests can aim the adapter at a
	// fake API.
	Endpoint string

	// UserAgent identifies this build in StatusCake's logs.
	UserAgent string
}

// validate rejects settings that would not fail until the first request, or
// that would quietly misbehave: a rate of zero never lets a second request
// through, and a negative retry count is meaningless.
func (c Config) validate() error {
	var errs []error
	if c.Token == "" {
		errs = append(errs, errors.New("API token is required"))
	}
	if c.RPS <= 0 {
		errs = append(errs, fmt.Errorf("rps must be greater than zero, got %d", c.RPS))
	}
	if c.Retries < 0 {
		errs = append(errs, fmt.Errorf("retries must not be negative, got %d", c.Retries))
	}
	if c.MinBackoff < 0 {
		errs = append(errs, fmt.Errorf("min backoff must not be negative, got %s", c.MinBackoff))
	}
	if c.MaxBackoff < c.MinBackoff {
		errs = append(errs, fmt.Errorf("max backoff (%s) must not be less than min backoff (%s)", c.MaxBackoff, c.MinBackoff))
	}
	return errors.Join(errs...)
}

// newClient returns a StatusCake API client configured by cfg. Requests are rate
// limited and, when they fail, retried with exponential backoff.
func newClient(cfg Config) (*statuscake.Client, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	bearer := credentials.NewBearerWithStaticToken(cfg.Token)
	opts := []statuscake.Option{
		statuscake.WithBackoff(backoff.Exponential{
			BaseDelay:  cfg.MinBackoff,
			Multiplier: 2.0,
			Jitter:     0.2,
			MaxDelay:   cfg.MaxBackoff,
		}),
		statuscake.WithHTTPClient(&http.Client{
			Transport: throttle.NewWithDefaultTransport(
				rate.NewLimiter(rate.Limit(cfg.RPS), 1),
			),
		}),
		statuscake.WithMaxRetries(cfg.Retries),
		statuscake.WithRequestCredentials(bearer),
		statuscake.WithUserAgent(cfg.UserAgent),
	}

	if cfg.Endpoint != "" {
		opts = append(opts, statuscake.WithHost(cfg.Endpoint))
	}

	return statuscake.NewClient(opts...), nil
}
