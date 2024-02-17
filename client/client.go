package client

import (
	"context"
)

// WithTransport sets the transport for the client
func WithTransport(transport Transport) OptionFunc {
	return func(o any) error {
		c, ok := o.(*Client)
		if !ok {
			return nil
		}
		c.transport = transport
		return nil
	}
}

// Client is the client for for the service
type Client struct {
	transport Transport
}

// NewClient creates a new Client
func NewClient(opts ...Option) (*Client, error) {
	c := &Client{}
	if err := ApplyOptions(c, opts...); err != nil {
		return nil, err
	}
	return c, nil
}

// Execute executes the operation and returns the output or an error (unknown or problem) if any
func (c *Client) Execute(ctx context.Context, op *Operation, input Struct) (any, error) {
	res, err := c.transport.Execute(ctx, &Request{
		Operation: op,
		Input:     input,
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}
