package route

import (
	"time"

	"github.com/maxlandon/wiregost/server/transport/resolver"
)

// ChainOptions holds options for Chain.
type ChainOptions struct {
	Retries int
	Timeout time.Duration
	// Hosts    *Hosts
	Resolver resolver.Resolver
}

// ChainOption allows a common way to set chain options.
type ChainOption func(opts *ChainOptions)

// RetryChainOption specifies the times of retry used by Chain.Dial.
func RetryChainOption(retries int) ChainOption {
	return func(opts *ChainOptions) {
		opts.Retries = retries
	}
}

// TimeoutChainOption specifies the timeout used by Chain.Dial.
func TimeoutChainOption(timeout time.Duration) ChainOption {
	return func(opts *ChainOptions) {
		opts.Timeout = timeout
	}
}

// HostsChainOption specifies the hosts used by Chain.Dial.
// func HostsChainOption(hosts *Hosts) ChainOption {
//         return func(opts *ChainOptions) {
//                 opts.Hosts = hosts
//         }
// }

// ResolverChainOption specifies the Resolver used by Chain.Dial.
func ResolverChainOption(resolver resolver.Resolver) ChainOption {
	return func(opts *ChainOptions) {
		opts.Resolver = resolver
	}
}
