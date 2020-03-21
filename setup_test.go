package torproxy

import (
	"bytes"
	"testing"

	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
)

func Test_parse(t *testing.T) {
	tests := []struct {
		configFile string
		config     Config
	}{
		{
			`
			torproxy from.com to.onion
			`,
			Config{
				To: map[string]string{"from.com": "http://to.onion"},
			},
		},
		{
			`
			torproxy from.com to.onion
			torproxy from2.com to2.onion
			`,
			Config{
				To: map[string]string{"from.com": "http://to.onion", "from2.com": "http://to2.onion"},
			},
		},
		{
			`
			torproxy from.com to.onion 
			torproxy from2.com to2.onion {
				port 4200
			}
			`,
			Config{
				To: map[string]string{"from.com": "http://to.onion", "from2.com": "http://to2.onion"},
				Client: &Tor{
					Host: "127.0.0.1",
					Port: 4200,
				},
			},
		},
		{
			`
			torproxy from.com to.onion 
			torproxy from2.com to2.onion {
				host 172.168.1.1
				port 4200
				datadir /data/dir
				torrc /etc/tor/torrc
				debug_mode true
				logfile /var/logs/stdout
			}
			`,
			Config{
				To: map[string]string{"from.com": "http://to.onion", "from2.com": "http://to2.onion"},
				Client: &Tor{
					Host:      "172.168.1.1",
					Port:      4200,
					DataDir:   "/data/dir",
					Torrc:     "/etc/tor/torrc",
					DebugMode: true,
					LogFile:   "/var/logs/stdout",
				},
			},
		},
	}
	for i, test := range tests {
		buf := bytes.NewBuffer([]byte(test.configFile))
		block, err := caddyfile.Parse("Caddyfile", buf)
		if err != nil {
			t.Errorf("Couldn't read the config: %s", err.Error())
		}

		// Extract the config tokens from the server blocks
		var tokens []caddyfile.Token
		for _, segment := range block[0].Segments {
			for _, token := range segment {
				tokens = append(tokens, token)
			}
		}

		d := caddyfile.NewDispenser(tokens)
		g := &TorProxy{Config: Config{Client: &Tor{}}, testing: true}

		if err := g.UnmarshalCaddyfile(d); err != nil {
			t.Errorf("Couldn't parse the config: %s", err.Error())
		}

		for from, to := range g.Config.To {
			if test.config.To[from] != to {
				t.Errorf("Expected %+v, Got %+v", test.config.To, g.Config.To)
			}
		}

		if test.config.Client != nil {
			if g.Config.Client.Host != test.config.Client.Host {
				t.Errorf("[%d]: Expected %s, but got %s", i, test.config.Client.Host, g.Config.Client.Host)
			}
			if g.Config.Client.Port != test.config.Client.Port {
				t.Errorf("[%d]: Expected %d, but got %d", i, test.config.Client.Port, g.Config.Client.Port)
			}
			if g.Config.Client.DataDir != test.config.Client.DataDir {
				t.Errorf("[%d]: Expected %s, but got %s", i, test.config.Client.DataDir, g.Config.Client.DataDir)
			}
			if g.Config.Client.Torrc != test.config.Client.Torrc {
				t.Errorf("[%d]: Expected %s, but got %s", i, test.config.Client.Torrc, g.Config.Client.Torrc)
			}
			if g.Config.Client.DebugMode != test.config.Client.DebugMode {
				t.Errorf("[%d]: Expected %t, but got %t", i, test.config.Client.DebugMode, g.Config.Client.DebugMode)
			}
			if g.Config.Client.LogFile != test.config.Client.LogFile {
				t.Errorf("[%d]: Expected %s, but got %s", i, test.config.Client.LogFile, g.Config.Client.LogFile)
			}
		}
	}
}
