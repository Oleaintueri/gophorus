package ports

type Protocol uint

const (
	PROTOCOL_UDP = iota
	PROTOCOL_TCP
)

func (p Protocol) Value() string {
	return []string{"udp", "tcp"}[p]
}

type Scheme uint

const (
	SCHEME_HTTP = iota
	SCHEME_HTTPS
)

func (s Scheme) Value() string {
	return []string{"http", "https"}[s]
}
