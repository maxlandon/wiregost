package transporter

import (
	"crypto/tls"
	"net"
	"time"
)

// DialOptions describes the options for Transporter.Dial.
type DialOptions struct {
	Timeout time.Duration
	// Chain   *Chain
}

// DialOption allows a common way to set DialOptions.
type DialOption func(opts *DialOptions)

// HandshakeOptions describes the options for handshake.
type HandshakeOptions struct {
	Addr string
	Host string
	// User       *url.Userinfo
	Timeout   time.Duration
	Interval  time.Duration
	Retry     int
	TLSConfig *tls.Config
	// WSOptions  *WSOptions
	// KCPConfig  *KCPConfig
	// QUICConfig *QUICConfig
}

// HandshakeOption allows a common way to set HandshakeOptions.
type HandshakeOption func(opts *HandshakeOptions)

func Test() {
	cc := net.Conn
}
