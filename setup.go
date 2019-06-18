package torproxy

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddy/caddymain"
	"github.com/mholt/caddy/caddyhttp/httpserver"
)

func main() {
	caddymain.EnableTelemetry = false
	caddymain.Run()
}

func init() {
	caddy.RegisterPlugin("torproxy", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
	// TODO: hardcode directive after stable release into Caddy
	httpserver.RegisterDevDirective("torproxy", "")
}

func parse(c *caddy.Controller) (Config, error) {
	var config Config
	var client Tor
	to := make(map[string]string)

	for c.Next() {
		if c.Val() == "torproxy" {
			c.Next() // skip directive name
		}

		// Parse the Config.From and Config.To URIs
		fromURI, err := url.Parse(c.Val())
		if err != nil {
			return Config{}, fmt.Errorf("Couldn't parse the `from` URI %s", err.Error())
		}
		toURI, err := url.Parse(c.RemainingArgs()[0])
		if err != nil {
			return Config{}, fmt.Errorf("Couldn't parse the `from` URI: %s", err.Error())
		}

		if toURI.Scheme == "" {
			toURI.Scheme = "http://"
		}

		// Fill the config instance
		to[fromURI.String()] = toURI.String()
	}

	config.To = to
	config.Client = &client
	config.Client.SetDefaults()
	return config, nil
}

func setup(c *caddy.Controller) error {
	config, err := parse(c)
	if err != nil {
		return err
	}

	// Add handler to Caddy
	cfg := httpserver.GetConfig(c)
	mid := func(next httpserver.Handler) httpserver.Handler {
		return TorProxy{
			Next:   next,
			Config: config,
		}
	}
	cfg.AddMiddleware(mid)

	config.Client.Start(c)

	c.OnShutdown(func() error {
		return config.Client.Stop()
	})

	return nil
}

type TorProxy struct {
	Next   httpserver.Handler
	Config Config
}

func (rd TorProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	if err := rd.Config.Proxy(w, r); err != nil {
		if err.Error() == "option disabled" {
			return rd.Next.ServeHTTP(w, r)
		}
		return http.StatusInternalServerError, err
	}

	return 0, nil
}
