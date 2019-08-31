package torproxy

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/cretz/bine/tor"
	"github.com/caddyserver/caddy"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

func (t *Tor) Start(c *caddy.Controller) {
	var debugger io.Writer
	if t.DebugMode {
		if t.LogFile != "" {
			debugger = &lumberjack.Logger{
				Filename:   t.LogFile,
				MaxSize:    100,
				MaxAge:     14,
				MaxBackups: 10,
			}
		}
		debugger = os.Stdout
	}

	torInstance, err := tor.Start(nil, t.starterConfig(debugger))
	if err != nil {
		log.Panicf("Unable to start Tor: %v", err)
	}

	listenCtx := context.Background()

	onion, err := torInstance.Listen(listenCtx, &tor.ListenConf{LocalPort: 8868, RemotePorts: []int{80}})
	if err != nil {
		log.Panicf("Unable to start onion service: %v", err)
	}

	t.onion = onion
	t.instance = torInstance
}

// Stop stops the tor instance, context listener and the onion service
func (t *Tor) Stop() error {
	if err := t.instance.Close(); err != nil {
		return fmt.Errorf("[torproxy]: Couldn't close the tor instance. %s", err.Error())
	}
	t.onion.Close()
	return nil
}

func (t *Tor) starterConfig(debugger io.Writer) *tor.StartConf {
	config := &tor.StartConf{
		NoAutoSocksPort: true,
		ExtraArgs:       []string{"--SocksPort", strconv.Itoa(t.Port)},
		TempDataDirBase: t.DataDir,
		TorrcFile:       t.Torrc,
	}

	if debugger != nil {
		config.DebugWriter = debugger
	}

	return config
}
