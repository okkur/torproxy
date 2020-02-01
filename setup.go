package torproxy

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os/exec"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

func init() {
	caddy.RegisterModule(TorProxy{})
	httpcaddyfile.RegisterHandlerDirective("torproxy", parse)
}

func parse(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var t TorProxy
	err := t.UnmarshalCaddyfile(h.Dispenser)
	return t, err
}

// UnmarshalCaddyfile implements caddyfile.Unmarshaler.
func (t *TorProxy) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	var config Config
	var client Tor
	to := make(map[string]string)

	for d.Next() {
		if d.Val() == "torproxy" {
			d.Next() // skip directive name
		}

		// Parse the Config.From and Config.To URIs
		fromURI, err := url.Parse(d.Val())
		if err != nil {
			return fmt.Errorf("Couldn't parse the `from` URI %s", err.Error())
		}
		toURI, err := url.Parse(d.RemainingArgs()[0])
		if err != nil {
			return fmt.Errorf("Couldn't parse the `from` URI: %s", err.Error())
		}

		if toURI.Scheme == "" {
			toURI.Scheme = "http"
		}

		// Fill the config instance
		to[fromURI.String()] = toURI.String()
	}

	config.To = to
	config.Client = &client
	config.Client.SetDefaults()
	return nil
}

// CaddyModule returns the Caddy module information.
func (TorProxy) CaddyModule() caddy.ModuleInfo {
	if err := isTorInstalled(); err != nil {
		log.Fatalf(err.Error())
	}

	return caddy.ModuleInfo{
		Name: "http.handlers.torproxy",
		New:  func() caddy.Module { return new(TorProxy) },
	}
}

// func setup(c *caddy.Controller) error {
// 	if err := isTorInstalled(); err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}

// 	config, err := parse(c)
// 	if err != nil {
// 		return err
// 	}

// 	// Add handler to Caddy
// 	cfg := httpserver.GetConfig(c)
// 	mid := func(next httpserver.Handler) httpserver.Handler {
// 		return TorProxy{
// 			Next:   next,
// 			Config: config,
// 		}
// 	}
// 	cfg.AddMiddleware(mid)

// 	config.Client.Start(c)

// 	c.OnShutdown(func() error {
// 		return config.Client.Stop()
// 	})

// 	return nil
// }

type TorProxy struct {
	Config Config
}

func (t TorProxy) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	return t.Config.Proxy(w, r)
}

func isTorInstalled() error {
	// Setup and run the "tor --version" command
	cmd := exec.Command("tor", "--version")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}

	// Read the output into buffer
	buf := new(bytes.Buffer)
	buf.ReadFrom(stdout)

	// Check if the output contains Tor's version
	if buf.String()[0:3] != "Tor" {
		return fmt.Errorf("Tor is not installed on you machine.Please follow these instructions to install Tor: https://www.torproject.org/download/")
	}

	return nil
}
