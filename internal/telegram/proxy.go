package telegram

import (
	"context"
	"net"
	"net/url"

	"github.com/gotd/td/telegram/dcs"
	"golang.org/x/net/proxy"
)

// This file is used to manually create a proxy with the arguments and system environment.
func createProxy(proxyURL string) (dcs.DialFunc, error) {
	if proxyURL != "" {
		u, err := url.Parse(proxyURL)
		if err != nil {
			return nil, err
		}

		dialer, err := proxy.FromURL(u, proxy.Direct)
		if err != nil {
			return nil, err
		}

		return func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialContext(ctx, dialer, network, addr)
		}, nil
	}

	// Fallback to default proxy with environment support.
	return proxy.Dial, nil
}

// Copied from golang.org/x/net/proxy/dial.go
func dialContext(ctx context.Context, d proxy.Dialer, network, address string) (net.Conn, error) {
	var (
		conn net.Conn
		done = make(chan struct{}, 1)
		err  error
	)
	go func() {
		conn, err = d.Dial(network, address)
		close(done)
		if conn != nil && ctx.Err() != nil {
			_ = conn.Close()
		}
	}()
	select {
	case <-ctx.Done():
		err = ctx.Err()
	case <-done:
	}
	return conn, err
}
