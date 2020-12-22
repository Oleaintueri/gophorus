package test

import (
	"github.com/Oleaintueri/gophorus/internal/pkg/ports"
	"testing"
)

func Test_ScanSingleIp(t *testing.T) {
	portScanner, err := ports.NewPortScanner("127.0.0.1",
		ports.WithPorts([]int{
		80,
		443,
		8000,
	}), ports.WithTimeout(2000))

	if err != nil {
		t.Error(err)
	}

	devices, err := portScanner.Scan()

	if err != nil {
		t.Error(err)
	}

	for i := range devices {
		t.Logf("%v", devices[i])
	}
}

func Test_ScanEntireCidr(t *testing.T) {
	portScanner, err := ports.NewPortScanner("127.0.0.1/24",
		ports.WithEntireCidr(true),
		ports.WithPorts([]int{
		80,
		443,
		8000,
	}),
	ports.WithReturnOnlyOpen(true))

	if err != nil {
		t.Error(err)
	}

	devices, err := portScanner.Scan()

	if err != nil {
		t.Error(err)
	}

	for i := range devices {
		t.Logf("%v", devices[i])
	}
}
