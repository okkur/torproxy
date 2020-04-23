package torproxy

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"

	"golang.org/x/net/proxy"
)

// Proxy redirects the request to the local onion service and the actual proxying
// happens inside onion service's http handler
func (c Config) Proxy(w http.ResponseWriter, r *http.Request) error {
	u, err := url.Parse(c.To[r.Host])
	if err != nil {
		return err
	}

	if u.Port() == "" {
		toPort := "80"
		if u.Scheme == "https" {
			toPort = "443"
		}
		u.Host = net.JoinHostPort(u.Host, toPort)
	}

	// Use proxied Host
	r.Host = u.Host

	// Create a socks5 dialer
	dialer, err := proxy.SOCKS5("tcp", fmt.Sprintf("%s:%d", c.Client.Host, c.Client.Port), nil, proxy.Direct)
	if err != nil {
		return fmt.Errorf("Couldn't connect to socks proxy: %s", err.Error())
	}

	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.Transport = &http.Transport{
		Dial: dialer.Dial,
	}

	// Create a temporary response writer to save response's body and headers
	tmpResponse := TorResponse{
		headers: make(http.Header),
		dialer:  dialer,
		request: r,
	}

	// Proxy the request and write the response to the temporary response writer
	proxy.ServeHTTP(&tmpResponse, r)

	// If the received response is a redirect, proxy the request to the response's Location header
	// Do this until the final response isn't a redirect response
	for {
		if tmpResponse.status == http.StatusFound || tmpResponse.status == http.StatusMovedPermanently {
			if err = tmpResponse.Redirect(); err != nil {
				return fmt.Errorf("[torproxy]: Couldn't redirect the request to the response's \"Location\" header: %s", err)
			}
			continue
		}
		break
	}

	// Decompress the body based on "Content-Encoding" header and write to a writer buffer
	if err := tmpResponse.WriteBody(); err != nil {
		return fmt.Errorf("[torproxy]: Couldn't write the response body: %s", err.Error())
	}

	// Replace the URL hosts with the request's host
	if err := tmpResponse.ReplaceBody(u.Scheme, u.Host, r.Host); err != nil {
		return fmt.Errorf("[torproxy]: Couldn't replace urls inside the response body: %s", err.Error())
	}

	copyHeader(w.Header(), tmpResponse.Header())

	// Write the status from the temporary ResponseWriter to the main ResponseWriter
	w.WriteHeader(tmpResponse.status)

	// Write the final response from the temporary ResponseWriter to the main ResponseWriter
	if _, err := w.Write(tmpResponse.Body()); err != nil {
		return fmt.Errorf("[torproxy]: Couldn't write the temporary response to main response body: %s", err.Error())
	}

	return nil
}

// Redirect redirects the request to the previous response's Location header.
func (t *TorResponse) Redirect() error {
	u, err := url.Parse(t.Header().Get("Location"))
	if err != nil {
		return fmt.Errorf("[torproxy]: Couldn't parse the URI from Redirect response: %s", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.Transport = &http.Transport{
		Dial: t.dialer.Dial,
	}

	proxy.ServeHTTP(t, t.request)

	return nil
}
