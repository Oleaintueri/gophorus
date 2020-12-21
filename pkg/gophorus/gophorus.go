package gophorus

import (
	"github.com/Oleaintueri/gophorus/internal/pkg/ports"
	"github.com/Oleaintueri/gophorus/internal/pkg/utility"
)

type Gophorus interface {
	Scan() ([]*utility.GenericDevice, error)
}

func NewPortScanner(ip string, opts ...ports.OptionPortScanner) (Gophorus, error) {
	return ports.NewPortScanner(ip, opts...)
}
