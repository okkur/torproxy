package torproxy

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	cproxy "github.com/mholt/caddy/caddyhttp/proxy"
	"golang.org/x/net/proxy"
)

// Proxy redirects the request to the local onion serivce and the actual proxying
// happens inside onion service's http handler
func (c Config) Proxy(w http.ResponseWriter, r *http.Request) error {
	u, err := url.Parse(c.To[r.URL.Host])
	if err != nil {
		return err
	}

	// Use proxied Host
	r.Host = u.Host

	// Create a socks5 dialer
	dialer, err := proxy.SOCKS5("tcp", fmt.Sprintf("127.0.0.1:%d", c.Client.Port), nil, proxy.Direct)
	if err != nil {
		log.Fatal(err)
	}

	reverseProxy := cproxy.NewSingleHostReverseProxy(u, "", torProxyKeepalive, torProxyTimeout, torFallbackDelay)
	reverseProxy.Transport = &http.Transport{
		Dial: dialer.Dial,
	}

	tmpResponse := TorResponse{headers: make(http.Header)}
	if err := reverseProxy.ServeHTTP(&tmpResponse, r, nil); err != nil {
		return fmt.Errorf("[txtdirect]: Coudln't proxy the request to the background onion service. %s", err.Error())
	}

	// Decompress the body based on "Content-Encoding" header and write to a writer buffer
	if err := tmpResponse.WriteBody(); err != nil {
		return fmt.Errorf("[txtdirect]: Couldn't write the response body: %s", err.Error())
	}

	// Replace the URL hosts with the request's host
	if err := tmpResponse.ReplaceBody(u.Scheme, u.Host, r.Host); err != nil {
		return fmt.Errorf("[txtdirect]: Couldn't replace urls inside the response body: %s", err.Error())
	}

	copyHeader(w.Header(), tmpResponse.Header())

	// Write the status from the temporary ResponseWriter to the main ResponseWriter
	w.WriteHeader(tmpResponse.status)

	// Write the final response from the temporary ResponseWriter to the main ResponseWriter
	if _, err := w.Write(tmpResponse.Body()); err != nil {
		return fmt.Errorf("[txtdirect]: Couldn't write the temporary response to main response body: %s", err.Error())
	}

	return nil
}
