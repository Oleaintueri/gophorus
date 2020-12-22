package upnp

import (
	"github.com/Oleaintueri/gophorus/internal/pkg"
	"github.com/Oleaintueri/gossdp/pkg/ssdp"
	"time"
)

type options struct {
	urn        string
	deviceName string
	timeout    time.Duration
}

type OptionUpnpScanner interface {
	apply(*options)
}

type urnOption string

func (u urnOption) apply(opts *options) {
	opts.urn = string(u)
}

type deviceNameOption string

func (d deviceNameOption) apply(opts *options) {
	opts.deviceName = string(d)
}

type timeoutOption int

func (t timeoutOption) apply(opts *options) {
	opts.timeout = time.Duration(t) * time.Millisecond
}

func WithUrn(urn string) OptionUpnpScanner {
	return urnOption(urn)
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
		urn:        ROOT_DEVICE.String(),
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
	ssdpClient := ssdp.NewSSDP(ssdp.WithTimeout(int(u.timeout)))

	devices, err := ssdpClient.SearchDevices(u.urn)

	if err != nil {
		return nil, err
	}

	for i := range devices {
		u.devices = append(u.devices, &pkg.GenericDevice{
			IP:         devices[i].ModelURL,
			Port:       0,
			Open:       true,
			DeviceType: devices[i].DeviceType,
			DeviceName: devices[i].FriendlyName,
		})

	}

	return u.devices, nil
}
