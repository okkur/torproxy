package torproxy

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/cretz/bine/tor"
	"github.com/caddyserver/caddy"
)

// TorProxy config
type Config struct {
	To     map[string]string
	Client *Tor
}

// Tor instance config struct
type Tor struct {
	// Socks5 proxy port
	Port      int
	DataDir   string
	Torrc     string
	DebugMode bool
	LogFile   string

	instance        *tor.Tor
	contextCanceler context.CancelFunc
	onion           *tor.OnionService
}

// DefaultOnionServicePort is the port used to serve the onion service on
const DefaultOnionServicePort = 4242

// TODO: Discuss these values
const (
	torProxyKeepalive = 30000000
	torFallbackDelay  = 30000000 * time.Millisecond
	torProxyTimeout   = 30000000 * time.Second
)

// ParseTor parses advanced config for Tor client
func (t *Tor) ParseTor(c *caddy.Controller) error {
	switch c.Val() {
	case "port":
		value, err := strconv.Atoi(c.RemainingArgs()[0])
		if err != nil {
			return fmt.Errorf("The given value for port field is not standard. It should be an integer")
		}
		t.Port = value

	case "datadir":
		t.DataDir = c.RemainingArgs()[0]

	case "torrc":
		t.Torrc = c.RemainingArgs()[0]

	case "debug_mode":
		value, err := strconv.ParseBool(c.RemainingArgs()[0])
		if err != nil {
			return fmt.Errorf("The given value for debug_mode field is not standard. It should be a boolean")
		}
		t.DebugMode = value

	case "logfile":
		t.LogFile = c.RemainingArgs()[0]

	default:
		return c.ArgErr() // unhandled option for tor
	}
	return nil
}

// SetDefaults sets the default values for prometheus config
// if the fields are empty
func (t *Tor) SetDefaults() {
	if t.Port == 0 {
		t.Port = DefaultOnionServicePort
	}
}
