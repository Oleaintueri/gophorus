/*
Author: Alano Terblanche (Benehiko)

Written with the guidance of
Kent Gruber's Medium article (https://medium.com/@KentGruber/building-a-high-performance-port-scanner-with-golang-9976181ec39d)

*/
package ports

import (
	"context"
	"fmt"
	"github.com/Oleaintueri/gophorus/internal/pkg/utility"
	"github.com/yl2chen/cidranger"
	"golang.org/x/sync/semaphore"
	"net"
	"strings"
	"sync"
	"time"
)

// Rest options for testing an endpoint on an open port
type restOptions struct {
	params   string      // the query parameters that should be added "foo=1&bar=2
	endpoint string      // the endpoint that should be targeted "/api/v1/..."
	payload  interface{} // The data that should be sent
	protocol Protocol    // UDP / TCP
	scheme   Scheme      // HTTPS / HTTP
	method   string      // HTTP method, POST, GET etc.
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
	restOptions
}

type OptionPortScanner interface {
	apply(*options)
}

type protocolOption Protocol

func (p protocolOption) apply(opts *options) {
	opts.protocol = Protocol(p)
}

type restfulOption restOptions

func (r restfulOption) apply(opts *options) {
	opts.restOptions = restOptions(r)
}

type timeoutOption int

func (t timeoutOption) apply(opts *options) {
	opts.timeout = time.Duration(t) * time.Millisecond
}

type portsOption []int

func (p portsOption) apply(opts *options) {
	opts.ports = p
}

type endpointOption string

func (e endpointOption) apply(opts *options) {
	opts.endpoint = string(e)
}

type payloadOption struct {
	payload interface{}
}

func (p payloadOption) apply(opts *options) {
	opts.payload = p.payload
}

type entireCidrOption bool

func (e entireCidrOption) apply(opts *options) {
	opts.entireCIDR = bool(e)
}

func WithTimeout(timeout int) OptionPortScanner {
	return timeoutOption(timeout)
}

func WithPorts(ports []int) OptionPortScanner {
	return portsOption(ports)
}

func WithEndpoint(endpoint string) OptionPortScanner {
	return endpointOption(endpoint)
}

func WithEntireCidr(entireCidr bool) OptionPortScanner {
	return entireCidrOption(entireCidr)
}

func WithRestful(restOpts restOptions) OptionPortScanner {
	return restfulOption(restOpts)
}

func WithProtocol(protocol Protocol) OptionPortScanner {
	return protocolOption(protocol)
}

type PortScanner struct {
	devices []*utility.GenericDevice
	lock    *semaphore.Weighted
	options
}

// Create a new port scanner object
// Accepts an IP address in the IPv4 format such as eg. "192.168.0.1".
// Only the passed IP address with be scanned unless WithEntireCidr is specified in the opts field.
//
// Allow scanning of the entire CIDR with opts WithEntireCidr.
// Add ports with opts WithPorts.
// Add a specified timeout in milliseconds with opts WithTimeout.
// Add a verifying request on top of the port check specifying restful options with opts WithRest
func NewPortScanner(ipAddr string, opts ...OptionPortScanner) (*PortScanner, error) {

	// ranger helps with parsing the ip address cidr properties
	ranger := cidranger.NewPCTrieRanger()

	ip := net.ParseIP(ipAddr)

	options := options{
		protocol:    PROTOCOL_UDP,
		timeout:     10,
		ports:       []int{80},
		entireCIDR:  false,
		restOptions: nil,
	}

	for _, o := range opts {
		o.apply(&options)
	}

	var devices []*utility.GenericDevice

	if options.entireCIDR {
		containingNetworks, err := ranger.ContainingNetworks(ip)

		if err != nil {
			return nil, err
		}

		for i := range containingNetworks {
			for y := range options.ports {
				devices = append(devices, &utility.GenericDevice{
					IP:         containingNetworks[i].Network().IP.String(),
					Port:       options.ports[y],
					Open:       false,
					DeviceType: "",
				})
			}
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
			} else {
				println(err.Error())
			}
		}(ps.devices[i], timeout)

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
