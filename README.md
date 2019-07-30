# Torproxy

An easy way to proxy your http requests through the onion network

 [![state](https://img.shields.io/badge/state-beta-blue.svg)]() [![release](https://img.shields.io/github/release/okkur/torproxy.svg)](https://github.com/okkur/torproxy/releases) [![license](https://img.shields.io/github/license/okkur/torproxy.svg)](LICENSE)

**NOTE: This is a beta release, we do not consider it completely production ready yet. Use at your own risk.**

Route your http requests through the onion network without Tor browser

## Using Torproxy

First you need to install the Tor on your machine. Check this page to download and learn how to install Tor: [Download Tor](https://www.torproject.org/download/)

You can get torproxy as a plugin on caddy's [build server](https://caddyserver.com/download).

Or install torproxy using `go get`.
```
go get -u go.okkur.org/torproxy/cmd/torproxy
```

Create a config file.
```
example.test {
  torproxy example.test somewhereonthe.onion 
}
```

Now you can run torproxy.
```
torproxy -conf tor.test
```
Take a look at our full [documentation](/docs).

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
