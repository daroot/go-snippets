package httptools

import (
	"net"
	"net/http"
	"time"
	// "golang.org/x/net/http2"
)

// NewClient creates an http.Client with reasonable defaults
// for timeouts and connection pooling.
// See these for some reasons why you don't want to directly use &http.Client{}
//   - https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779
//   - https://simon-frey.com/blog/go-as-in-golang-standard-net-http-config-will-break-your-production/
//   - https://blog.cloudflare.com/exposing-go-on-the-internet/
//
//nolint:gomnd // the point of this code is to setup all the magic defaults
func NewClient() *http.Client {
	tr := &http.Transport{
		DialContext: (&net.Dialer{
			KeepAlive: time.Second * 30,
			DualStack: true,
			Timeout:   time.Millisecond * 500,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       time.Second * 300,
		TLSHandshakeTimeout:   time.Second,
		ExpectContinueTimeout: time.Second,
		ResponseHeaderTimeout: time.Second * 5,
		Proxy:                 http.ProxyFromEnvironment,
	}
	// _ = http2.ConfigureTransport(tr)
	return &http.Client{Transport: tr}
}
