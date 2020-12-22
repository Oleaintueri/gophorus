/*
Author: Alano Terblanche (Benehiko)

Written with the guidance of
Kent Gruber's Medium article (https://medium.com/@KentGruber/building-a-high-performance-port-scanner-with-golang-9976181ec39d)
Source: https://gist.github.com/picatz/9c0028efd7b3ced3329f7322f41b16e1#file-port_scanner-go

*/
package ports

import (
	"context"
	"fmt"
	"github.com/Oleaintueri/gophorus/internal/pkg/utility"
	"github.com/projectdiscovery/mapcidr"
	"golang.org/x/sync/semaphore"
	"net"
	"strings"
	"sync"
	"time"
)

// Rest options for testing an endpoint on an open port
type restOptions struct {
	// the query parameters that should be added "foo=1&bar=2
	params string
	// the endpoint that should be targeted "/api/v1/..."
	endpoint string
	// The data that should be sent
	payload interface{}
	// UDP / TCP
	protocol Protocol
	// HTTPS / HTTP
	scheme Scheme
	// HTTP method, POST, GET etc.
	method string
}

// internal port scanner options
type options struct {
	protocol Protocol
	// the timeout of the requests
	timeout time.Duration
	// an array of ports
	ports []int
	// specify if the entire cidr should be scanned
	entireCIDR bool
	// return only open ports
	returnOnlyOpen bool

	*restOptions
}

type OptionPortScanner interface {
	apply(*options)
}

type returnOnlyOpenOption bool

func (r returnOnlyOpenOption) apply(opts *options) {
	opts.returnOnlyOpen = bool(r)
}

type protocolOption Protocol

func (p protocolOption) apply(opts *options) {
	opts.protocol = Protocol(p)
}

type restfulOption struct {
	restOptions *restOptions
}

func (r restfulOption) apply(opts *options) {
	opts.restOptions = r.restOptions
}

type timeoutOption int

func (t timeoutOption) apply(opts *options) {
	opts.timeout = time.Duration(t) * time.Millisecond
}

type portsOption []int

func (p portsOption) apply(opts *options) {
	opts.ports = p
}

type entireCidrOption bool

func (e entireCidrOption) apply(opts *options) {
	opts.entireCIDR = bool(e)
}

// Timeout in milliseconds
// Default is 1000
func WithTimeout(timeout int) OptionPortScanner {
	return timeoutOption(timeout)
}

// Custom list of ports to scan
// Default is 80
func WithPorts(ports []int) OptionPortScanner {
	return portsOption(ports)
}

// Set to true to scan the entire CIDR block
// Default is false
func WithEntireCidr(entireCidr bool) OptionPortScanner {
	return entireCidrOption(entireCidr)
}

// Add a restful endpoint to test when port is open
func WithRestful(restOpts restOptions) OptionPortScanner {
	return restfulOption{&restOpts}
}

// Set the protocol, either UDP or TCP
// Default is TCP
func WithProtocol(protocol Protocol) OptionPortScanner {
	return protocolOption(protocol)
}

// Set to true to only return open devices
// Default is false
func WithReturnOnlyOpen(onlyOpen bool) OptionPortScanner {
	return returnOnlyOpenOption(onlyOpen)
}

type PortScanner struct {
	devices []*utility.GenericDevice
	lock    *semaphore.Weighted
	*options
}

// Create a new port scanner object
// Accepts an IP address in the IPv4 format such as eg. "192.168.0.1" or when querying the whole CIDR "192.168.0.0/24"
// Only the passed IP address with be scanned unless WithEntireCidr is specified in the opts field.
//
// Allow scanning of the entire CIDR with opts WithEntireCidr.
// Add ports with opts WithPorts.
// Add a specified timeout in milliseconds with opts WithTimeout.
// Add a verifying request on top of the port check specifying restful options with opts WithRest
// Change return value of the Scan request by only returning open ports with opts WithReturnOnlyOpen
func NewPortScanner(ipAddr string, opts ...OptionPortScanner) (*PortScanner, error) {

	options := &options{
		protocol:   PROTOCOL_UDP,
		timeout:    1000,
		ports:      []int{80},
		entireCIDR: false,
		restOptions: &restOptions{
			params:   "",
			endpoint: "",
			payload:  nil,
			protocol: PROTOCOL_TCP,
			scheme:   SCHEME_HTTP,
			method:   "",
		},
	}

	for _, o := range opts {
		o.apply(options)
	}

	var devices []*utility.GenericDevice

	if options.entireCIDR {
		ips, err := mapcidr.IPAddresses(ipAddr)

		if err != nil {
			return nil, err
		}

		for i := range ips {
			for y := range options.ports {
				devices = append(devices, &utility.GenericDevice{
					IP:         ips[i],
					Port:       options.ports[y],
					Open:       false,
					DeviceType: "",
				})
			}
		}
	} else {
		for y := range options.ports {
			devices = append(devices, &utility.GenericDevice{
				IP:         ipAddr,
				Port:       options.ports[y],
				Open:       false,
				DeviceType: "",
			})
		}
	}

	return &PortScanner{
		devices: devices,
		lock:    semaphore.NewWeighted(utility.Ulimit()),
		options: options,
	}, nil
}

// Scan the network with the configurations set in NewPortScanner
func (ps *PortScanner) Scan() ([]*utility.GenericDevice, error) {
	wg := sync.WaitGroup{}

	timeout := ps.options.timeout
	protocol := ps.options.protocol.Value()

	for i := range ps.devices {
		wg.Add(1)
		// Lock the Semaphore
		err := ps.lock.Acquire(context.TODO(), 1)
		if err != nil {
			return nil, err
		}

		go func(device *utility.GenericDevice, timeout time.Duration) {
			// Once anonymous function is done executing the semaphore will release
			defer ps.lock.Release(1)

			isOpen, err := scanPort(device.IP, device.Port, protocol, timeout, &wg)
			if err == nil {
				if isOpen {
					device.Open = true
				}
			}

		}(ps.devices[i], timeout)

	}

	if ps.returnOnlyOpen {
		var onlyOpen []*utility.GenericDevice

		for i := range ps.devices {
			if ps.devices[i].Open {
				onlyOpen = append(onlyOpen, ps.devices[i])
			}
		}

		return onlyOpen, nil

	}

	return ps.devices, nil
}

// internal function for scanning the port
func scanPort(ip string, port int, protocol string, timeout time.Duration, wg *sync.WaitGroup) (isOpen bool,
	err error) {
	// wg will call Done at the end of the function's execution using defer
	defer wg.Done()

	// Check the port. If the connection throws an error, the port is closed.
	target := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout(protocol, target, timeout)

	if err != nil {
		// Wait a bit if the system complains about too many open files
		if strings.Contains(err.Error(), "too many open files") {
			time.Sleep(timeout)
			return scanPort(ip, port, protocol, timeout, wg)
		}
		return false, err
	}

	if conn != nil {
		err = conn.Close()
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}
