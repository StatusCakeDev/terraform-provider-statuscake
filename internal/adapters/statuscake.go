package adapters

import "github.com/StatusCakeDev/statuscake-go"

// pageSize is how many items are requested at a time, and maxPages bounds how
// many pages are read for one listing so that an unexpectedly large account
// cannot turn one command into an unbounded number of requests.
const (
	pageSize = 100
	maxPages = 30
)

// Client reads monitoring data from the StatusCake API.
type Client struct {
	client *statuscake.Client
}

// New returns a Client that reaches the API as cfg describes. Nothing outside
// this package holds an SDK type, so replacing the SDK means replacing this
// package and nothing else.
func New(cfg Config) (*Client, error) {
	client, err := newClient(cfg)
	if err != nil {
		return nil, err
	}
	return &Client{client: client}, nil
}

// GetClient returns the underlying StatusCake SDK client. This is used by
// the provider to make API calls, and is not intended for use by other
// packages. This is necessary right now because the provider makes calls to the
// SDK directly, but in the future the provider should be refactored to use this
// package instead.
func (c *Client) GetClient() *statuscake.Client {
	return c.client
}
