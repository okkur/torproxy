package torproxy

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"

	"github.com/caddyserver/caddy"
	"github.com/caddyserver/caddy/caddy/caddymain"
	"github.com/caddyserver/caddy/caddyhttp/httpserver"
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
			toURI.Scheme = "http"
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
	if err := isTorInstalled(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

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

func isTorInstalled() error {
	// Setup and run the "tor --version" command
	cmd := exec.Command("tor", "--version")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
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
