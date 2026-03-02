/**
 * @author: dn-jinmin/dn-jinmin
 * @doc:
 */

package websocket

import (
	"time"
)

type ServerOptions func(opt *serverOption)

type serverOption struct {
	Authentication
	pattern string

	maxConnectionIdle time.Duration
}

func newServerOptions(opts ...ServerOptions) serverOption {
	o := serverOption{
		Authentication:    new(authentication),
		pattern:           "/ws",
		maxConnectionIdle: defaultMaxConnectionIdle,
	}

	for _, opt := range opts {
		opt(&o)
	}
	return o
}

func WithServerAuthentication(auth Authentication) ServerOptions {
	return func(opt *serverOption) {
		opt.Authentication = auth
	}
}

func WithServerPatten(pattern string) ServerOptions {
	return func(opt *serverOption) {
		opt.pattern = pattern
	}
}

func WithServerMaxConnectionIdle(maxConnectionIdle time.Duration) ServerOptions {
	return func(opt *serverOption) {
		if maxConnectionIdle > 0 {
			opt.maxConnectionIdle = maxConnectionIdle
		}
	}
}
