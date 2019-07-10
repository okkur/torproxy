package torproxy

import (
	"testing"

	"github.com/caddyserver/caddy"
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
		c := caddy.NewTestController("http", test.configFile)
		config, err := parse(c)
		if err != nil {
			t.Error(err)
		}

		for from, to := range config.To {
			if test.config.To[from] != to {
				t.Errorf("Expected %+v, Got %+v", test.config.To, config.To)
			}
		}
	}
}
