package upnp

import (
	"github.com/Oleaintueri/gophorus/internal/pkg"
	"github.com/Oleaintueri/gossdp/pkg/ssdp"
	"time"
)

type options struct {
	deviceName string
	timeout    time.Duration
}

type OptionUpnpScanner interface {
	apply(*options)
}

type deviceNameOption string

func (d deviceNameOption) apply(opts *options) {
	opts.deviceName = string(d)
}

type timeoutOption int

func (t timeoutOption) apply(opts *options) {
	opts.timeout = time.Duration(t) * time.Millisecond
}

func WithDeviceName(name string) OptionUpnpScanner {
	return deviceNameOption(name)
}

func WithTimeout(timeout int) OptionUpnpScanner {
	return timeoutOption(timeout)
}

type UpnpScanner struct {
	devices []*pkg.GenericDevice
	*options
}

func NewUpnp(opts ...OptionUpnpScanner) *UpnpScanner {
	options := &options{
		deviceName: "rootdevice",
		timeout:    1000,
	}

	for _, o := range opts {
		o.apply(options)
	}

	var devices []*pkg.GenericDevice

	return &UpnpScanner{
		devices: devices,
		options: options,
	}

}

func (u *UpnpScanner) Scan() ([]*pkg.GenericDevice, error) {
	ssdpClient := ssdp.NewSSDP(ssdp.WithTimeout(2000))

	devices, err := ssdpClient.SearchDevices("upnp:rootdevice")

	if err != nil {
		return nil, err
	}

	for i := range devices {
		u.devices = append(u.devices, &pkg.GenericDevice{
			IP:         "",
			Port:       0,
			Open:       false,
			DeviceType: devices[i].DeviceType,
			DeviceName: devices[i].FriendlyName,
		})

	}

	return u.devices, nil
}
