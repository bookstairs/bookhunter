package telegram

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"

	"github.com/gotd/td/telegram/dcs"
	"golang.org/x/net/proxy"

	"github.com/bookstairs/bookhunter/internal/log"
)

// Register Dialer Type for HTTP & HTTPS Proxy in golang.
func init() {
	proxy.RegisterDialerType("http", newHTTPProxy)
	proxy.RegisterDialerType("https", newHTTPProxy)
}

// This file is used to manually create a proxy with the arguments and system environment.

// CreateProxy is used to create a dcs.DialFunc for the telegram to send request.
// We don't support MTProxy now.
func CreateProxy(proxyURL string) (dcs.DialFunc, error) {
	if proxyURL != "" {
		log.Debugf("Try to manually create the proxy through %s", proxyURL)

		u, err := url.Parse(proxyURL)
		if err != nil {
			return nil, err
		}

		dialer, err := proxy.FromURL(u, Direct)
		if err != nil {
			return nil, err
		}

		return func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialContext(ctx, dialer, network, addr)
		}, nil
	}

	// Fallback to default proxy with environment support.
	dialer := proxy.FromEnvironment()
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		return dialContext(ctx, dialer, network, addr)
	}, nil
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

type direct struct{}

// Direct is a direct proxy: one that makes network connections directly.
var Direct = direct{}

func (direct) Dial(network, addr string) (net.Conn, error) {
	return net.Dial(network, addr)
}

// httpProxy is an HTTP / HTTPS connection proxy.
type httpProxy struct {
	host     string
	haveAuth bool
	username string
	password string
	forward  proxy.Dialer
}

func newHTTPProxy(uri *url.URL, forward proxy.Dialer) (proxy.Dialer, error) {
	s := new(httpProxy)
	s.host = uri.Host
	s.forward = forward
	if uri.User != nil {
		s.haveAuth = true
		s.username = uri.User.Username()
		s.password, _ = uri.User.Password()
	}

	return s, nil
}

func (s *httpProxy) Dial(_, addr string) (net.Conn, error) {
	// Dial and create the https client connection.
	c, err := s.forward.Dial("tcp", s.host)
	if err != nil {
		return nil, err
	}

	// HACK. http.ReadRequest also does this.
	reqURL, err := url.Parse("http://" + addr)
	if err != nil {
		_ = c.Close()
		return nil, err
	}
	reqURL.Scheme = ""

	req, err := http.NewRequestWithContext(context.Background(), "CONNECT", reqURL.String(), http.NoBody)
	if err != nil {
		_ = c.Close()
		return nil, err
	}
	req.Close = false
	if s.haveAuth {
		req.SetBasicAuth(s.username, s.password)
	}

	err = req.Write(c)
	if err != nil {
		_ = c.Close()
		return nil, err
	}

	resp, err := http.ReadResponse(bufio.NewReader(c), req)
	if err != nil {
		_ = resp.Body.Close()
		_ = c.Close()
		return nil, err
	}
	_ = resp.Body.Close()
	if resp.StatusCode != 200 {
		_ = c.Close()
		err = fmt.Errorf("connect server using proxy error, StatusCode [%d]", resp.StatusCode)
		return nil, err
	}

	return c, nil
}
