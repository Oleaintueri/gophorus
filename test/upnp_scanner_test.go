package test

import (
	"github.com/Oleaintueri/gophorus/internal/pkg/upnp"
	"testing"
)

func Test_UpnpScanner(t *testing.T) {
	upnpClient := upnp.NewUpnp(upnp.WithUrn("sddp:all"))

	devices, err := upnpClient.Scan()

	if err != nil {
		t.Error(err)
	}

	for i := range devices {
		t.Logf("Device: %v", devices[i])
	}

}
