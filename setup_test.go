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
	}
	for _, test := range tests {
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
		g := &TorProxy{}

		if err := g.UnmarshalCaddyfile(d); err != nil {
			t.Errorf("Couldn't parse the config: %s", err.Error())
		}

		for from, to := range g.Config.To {
			if test.config.To[from] != to {
				t.Errorf("Expected %+v, Got %+v", test.config.To, g.Config.To)
			}
		}
	}
}
