package connector

// ConnectOptions describes the options for Connector.Connect.
type ConnectOptions struct {
	// Addr      string
	// Timeout   time.Duration
	// User      *url.Userinfo
	// Selector  gosocks5.Selector
	// UserAgent string
	// NoTLS     bool
}

// ConnectOption allows a common way to set ConnectOptions.
type ConnectOption func(opts *ConnectOptions)
