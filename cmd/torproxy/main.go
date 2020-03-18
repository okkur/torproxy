package main

import (
	caddycmd "github.com/caddyserver/caddy/v2/cmd"
	_ "go.okkur.org/torproxy"
)

func main() {
	caddycmd.Main()
}
