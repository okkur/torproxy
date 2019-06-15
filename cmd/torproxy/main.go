package main

import (
	"github.com/mholt/caddy/caddy/caddymain"
	_ "github.com/okkur/torproxy"
)

func main() {
	caddymain.EnableTelemetry = false
	caddymain.Run()
}
