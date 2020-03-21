## Configuration

Torproxy provides two methods for configuration, one way to configure Torproxy is
to only provide the source and target URIs like this example:

```
example.test {
  torproxy example.test somewhereonthe.onion
}
```

In this case, Torproxy will use the default values to start a Tor instance on port 4242.
Also note that you can specify more than one URI to proxy, like this example:

```
first.test, second.test, third.test {
  torproxy first.test somewhereonthe.onion
  torproxy second.test somewhereonthe.onion
  torproxy third.test somewhereonthe.onion
}
```

If you want to use a custom Tor instance or customize the Tor instance that Torproxy starts,
use the following options in the config:

- `host`: Tor daemon's host (Default: `127.0.0.1`)
- `port`: Tor daemon's port (Default: `4242`)
- `datadir`: DataDir is the parent directory that a temporary data directory will be created under for use by Tor.
- `torrc`: Tor's configuration file. If empty, a blank torrc is created in the data directory and is used instead.
- `debug_mode`: If enabled, debug logs would be written to `stdout`.
- `logfile`: Path to a file for storing logs. `debug_mode` should be enabled for this option.

**Note: If `host` field isn't provided, Torproxy will start a fresh Tor instance itself.**

Example:

```
from2.com {
  torproxy from2.com to2.onion {
    host 172.168.1.1
    port 4200
    datadir /data/dir
    torrc /etc/tor/torrc
    debug_mode true
    logfile /var/logs/stdout
  }
}
```
