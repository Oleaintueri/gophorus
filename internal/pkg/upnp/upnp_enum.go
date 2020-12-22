package upnp

type URN uint

const (
	ROOT_DEVICE URN = iota
)

func (u URN) String() string {
	return []string{"upnp:rootdevice"}[u]
}
