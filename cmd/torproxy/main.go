package main

import (
	"github.com/caddyserver/caddy/caddy/caddymain"
	_ "go.okkur.org/torproxy"
)

func main() {
	caddymain.EnableTelemetry = false
	caddymain.Run()
}
