<h1 align="center">Gophorus</h1>

<p align="center">
    <img alt="GitHub" src="https://img.shields.io/github/license/Oleaintueri/gophorus?style=flat-square">
    <img alt="GitHub go.mod Go version"
         src="https://img.shields.io/github/go-mod/go-version/Oleaintueri/gophorus?style=flat-square">
    <img alt="GitHub tag (latest SemVer)"
         src="https://img.shields.io/github/v/tag/Oleaintueri/gophorus?style=flat-square">
</p>


A really fast network utility library for scanning network devices.

Features:

- Port Scanning
- UPnP (ssdp)

## Getting started

### Installation

    GO111MODULE=on go get github.com/Oleaintueri/gophorus

### Usage

```go
package main

import (
	"fmt"
	"github.com/Oleaintueri/gophorus/internal/pkg/ports"
	"github.com/Oleaintueri/gophorus/pkg/gophorus"
)

func main() {
	// A singular IP to scan
	gp, err := gophorus.NewPortScanner("192.168.0.1",
		ports.WithReturnOnlyOpen(true),
		ports.WithPorts([]int{
			443,
			8000,
			9000,
		}))

	// Scan an entire cidr
	gp, err = gophorus.NewPortScanner("192.168.0.1",
		ports.WithEntireCidr(true),
		ports.WithReturnOnlyOpen(true),
		ports.WithTimeout(2000),
		ports.WithProtocol(ports.PROTOCOL_TCP),
		ports.WithPorts([]int{
			80,
			443,
			8000,
			9000,
		}))

	if err != nil {
		panic(err)
	}

	devices, err := gp.Scan()

	if err != nil {
		panic(err)
	}

	for i := range devices {
		fmt.Printf("Device: %v", devices[i])
	}
}
```

## Other Resources

- [Channels in Go](https://medium.com/rungo/anatomy-of-channels-in-go-concurrency-in-go-1ec336086adb)

- [Kent Gruber's Medium Article](https://medium.com/@KentGruber/building-a-high-performance-port-scanner-with-golang-9976181ec39d)

- [Kent Gruber's Source Code](https://gist.github.com/picatz/9c0028efd7b3ced3329f7322f41b16e1#file-port_scanner-go)