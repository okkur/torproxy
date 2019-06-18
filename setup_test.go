package torproxy

import (
	"testing"

	"github.com/mholt/caddy"
)

func Test_parse(t *testing.T) {
	type args struct {
		c *caddy.Controller
	}
	tests := []struct {
		configFile string
		config     Config
	}{
		{
			`
			torproxy from.com to.onion
			`,
			Config{
				From: []string{"from.com"},
				To:   []string{"to.onion"},
			},
		},
		{
			`
			torproxy from.com to.onion
			torproxy from2.com to2.onion
			`,
			Config{
				From: []string{"from.com", "from2.com"},
				To:   []string{"to.onion", "to2.onion"},
			},
		},
	}
	for _, test := range tests {
		c := caddy.NewTestController("http", test.configFile)
		config, err := parse(c)
		if err != nil {
			t.Error(err)
		}

		// Check Config.From
		for i, from := range test.config.From {
			if config.From[i] != from {
				t.Errorf("Expected %v, Got %v", test.config.From, config.From)
			}
		}

		// Check Config.To
		for i, to := range test.config.To {
			if config.To[i] != to {
				t.Errorf("Expected %v, Got %v", test.config.To, config.To)
			}
		}
	}
}
