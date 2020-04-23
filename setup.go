package torproxy

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

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

config_handler:
	for d.Next() {
		if d.Val() == "torproxy" {
			d.Next() // skip directive name
		}

		// Parse the Tor client's config
		if d.Val() == "{" {
			for d.Next() {
				if d.Val() == "}" {
					continue config_handler
				}
				if err := client.ParseTor(d); err != nil {
					return fmt.Errorf("Couldn't parse the Tor client config: %s", err.Error())
				}
			}
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
	t.Config = config

	if t.Config.Client.Host == "" && !t.testing {
		t.Config.Client.Start()
	}

	return nil
}

// CaddyModule returns the Caddy module information.
func (TorProxy) CaddyModule() caddy.ModuleInfo {
	t := &TorProxy{Config: Config{Client: &Tor{}}}
	pool := caddy.NewUsagePool()
	pool.LoadOrNew("torclient", TorConstructor)

	if t.Config.Client.Host == "" && !t.testing {
		if err := t.Config.Client.IsInstalled(); err != nil {
			log.Fatalf(err.Error())
		}
	}

	return caddy.ModuleInfo{
		Name: "http.handlers.torproxy",
		New: func() caddy.Module {
			return t
		},
	}
}

type TorProxy struct {
	Config Config
	// Set "testing" to true in tests to skip Tor client's auto start
	testing bool
}

func (t TorProxy) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	return t.Config.Proxy(w, r)
}
