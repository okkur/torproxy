package main

import (
	"github.com/caddyserver/caddy/caddy/caddymain"
	_ "github.com/okkur/torproxy"
)

func main() {
	caddymain.EnableTelemetry = false
	caddymain.Run()
}
