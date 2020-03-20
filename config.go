package torproxy

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/cretz/bine/tor"
	"gopkg.in/natefinch/lumberjack.v2"
)

// TorProxy config
type Config struct {
	To     map[string]string
	Client *Tor
}

// Tor instance config struct
type Tor struct {
	// Socks5 proxy port
	Host      string
	Port      int
	DataDir   string
	Torrc     string
	DebugMode bool
	LogFile   string

	debugger        io.Writer
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
func (t *Tor) ParseTor(d *caddyfile.Dispenser) error {
	switch d.Val() {
	case "host":
		t.Host = d.RemainingArgs()[0]

	case "port":
		value, err := strconv.Atoi(d.RemainingArgs()[0])
		if err != nil {
			return fmt.Errorf("The given value for port field is not standard. It should be an integer")
		}
		t.Port = value

	case "datadir":
		t.DataDir = d.RemainingArgs()[0]

	case "torrc":
		t.Torrc = d.RemainingArgs()[0]

	case "debug_mode":
		value, err := strconv.ParseBool(d.RemainingArgs()[0])
		if err != nil {
			return fmt.Errorf("The given value for debug_mode field is not standard. It should be a boolean")
		}
		t.DebugMode = value

	case "logfile":
		t.LogFile = d.RemainingArgs()[0]

	default:
		return d.ArgErr() // unhandled option for tor
	}

	return nil
}

// SetDefaults sets the default values for prometheus config
// if the fields are empty
func (t *Tor) SetDefaults() {
	if t.DebugMode {
		if t.LogFile != "" {
			t.debugger = &lumberjack.Logger{
				Filename:   t.LogFile,
				MaxSize:    100,
				MaxAge:     14,
				MaxBackups: 10,
			}
		}
		t.debugger = os.Stdout
	}

	if t.Port == 0 {
		t.Port = DefaultOnionServicePort
	}
}

// TorConstructor return a new instance of Tor client struct.
// Used to manage the Tor client's life cycle
func TorConstructor() (caddy.Destructor, error) {
	return &Tor{}, nil
}

// Destruct stops the Tor client
func (t *Tor) Destruct() error {
	return t.Stop()
}

// IsInstalled checks the Tor client using the `tor --version` command
func (t *Tor) IsInstalled() error {
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
