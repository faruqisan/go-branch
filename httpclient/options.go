package httpclient

import (
	"time"
)

// Option define an option for the client
type Option func(*Client)

// WithHTTPTimeout configure the client to have specified http timeout
func WithHTTPTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.httpTimeout = timeout
	}
}

// WithMaxConcurrentRequest sets max number of concurrent requests a client
// could make
func WithMaxConcurrentRequest(max int) Option {
	return func(c *Client) {
		c.maxConcurrentReq = max
	}
}

// WithRequestVolumeThreshold sets minimum number of requests in a rolling window
// that will trip the circuit.
func WithRequestVolumeThreshold(vol int) Option {
	return func(c *Client) {
		c.reqVolThreshold = vol
	}
}

// WithSleepWindow sets the amount of time, after tripping the circuit,to reject requests
// before allowing attempts again to determine if the circuit should again be closed.
func WithSleepWindow(sleepWindow time.Duration) Option {
	return func(c *Client) {
		c.sleepWindow = sleepWindow
	}
}

// WithErrorPercentThreshold sets the error percentage at or above which the circuit
// should trip open and start short-circuiting requests to fallback logic.
func WithErrorPercentThreshold(threshold int) Option {
	return func(c *Client) {
		c.errorPercentThreshold = threshold
	}
}
