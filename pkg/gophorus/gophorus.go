package gophorus

import (
	"github.com/Oleaintueri/gophorus/internal/pkg"
	"github.com/Oleaintueri/gophorus/internal/pkg/ports"
	"github.com/Oleaintueri/gophorus/internal/pkg/upnp"
)

type Gophorus interface {
	Scan() ([]*pkg.GenericDevice, error)
}

func NewPortScanner(ip string, opts ...ports.OptionPortScanner) (Gophorus, error) {
	return ports.NewPortScanner(ip, opts...)
}

func NewUpnpScanner() (Gophorus, error) {
	return upnp.NewUpnp(), nil
}
