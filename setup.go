package torproxy

import (
	"fmt"
	"net/http"

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

func parse(c *caddy.Controller) (Config, error) { return Config{}, nil }

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

	c.OnShutdown(func() error {
		return nil
	})

	return nil
}

type TorProxy struct {
	Next   httpserver.Handler
	Config Config
}

func (rd TorProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	if err := fmt.Errorf("Hello World"); err != nil {
		if err.Error() == "option disabled" {
			return rd.Next.ServeHTTP(w, r)
		}
		return http.StatusInternalServerError, err
	}

	return 0, nil
}
