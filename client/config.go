package client

// Option is a function that can be used to configure a client.
type Option interface {
	Apply(c any) error
}

// Options is a list of options
type Options []Option

// OptionFunc is an adapter to allow the use of ordinary functions as options.
type OptionFunc func(c any) error

// Apply applies the option to the given client.
func (f OptionFunc) Apply(c any) error {
	return f(c)
}

// ApplyOptions applies the given options to the given client.
func ApplyOptions(c any, options ...Option) error {
	for _, option := range options {
		if err := option.Apply(c); err != nil {
			return err
		}
	}
	return nil
}
