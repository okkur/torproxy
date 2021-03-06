# Torproxy

An easy way to proxy your http requests through the onion network

 [![state](https://img.shields.io/badge/state-beta-blue.svg)]() [![release](https://img.shields.io/github/release/okkur/torproxy.svg)](https://torproxy.okkur.org/releases/) [![license](https://img.shields.io/github/license/okkur/torproxy.svg)](LICENSE)

**NOTE: This is a beta release, we do not consider it completely production ready yet. Use at your own risk.**

Route your http requests through the onion network without Tor browser

## Using Torproxy
Note: The `master` branch is using [Caddy v2](https://caddyserver.com/), if you
want to use Torproxy with previous caddy versions, check the `caddy-v1` branch.

If you want Torproxy to start a fresh instance of Tor, you need to install Tor
on your machine. Take a look at the [Tor download instructions](https://www.torproject.org/download/)

You can get torproxy as a plugin on Caddy's `v1` [build server](https://caddyserver.com/download).

Or install torproxy's `v2` version using `go get`:
```
go get go.okkur.org/torproxy/cmd/torproxy
```

Create a config file like the example below. For more information about the
available config options, check the [Configuration](/docs/configuration.md) page.
```
example.test {
  torproxy example.test somewhereonthe.onion 
}
```

You can run torproxy with Caddyfile config adapter using this command:
```
torproxy start -config torproxy.config -adapter caddyfile
```
Take a look at our full [documentation](https://torproxy.okkur.org/docs).

## Support
For detailed information on support options see our [support guide](/SUPPORT.md).

## Helping out
Best place to start is our [contribution guide](/CONTRIBUTING.md).

----

*Code is licensed under the [Apache License, Version 2.0](/LICENSE).*  
*Documentation/examples are licensed under [Creative Commons BY-SA 4.0](/docs/LICENSE).*  
*Illustrations, trademarks and third-party resources are owned by their respective party and are subject to different licensing.*

---

Copyright 2019 - The Torproxy authors
