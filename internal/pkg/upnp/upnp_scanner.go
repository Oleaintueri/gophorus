package upnp

import "github.com/huin/goupnp"

func NewUpnp() {
	goupnp.NewServiceClients("urn:rootdevice")
}