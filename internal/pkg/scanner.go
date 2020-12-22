package pkg

type GenericDevice struct {
	IP         string `json:"ip"`
	Port       int    `json:"port"`
	Open       bool   `json:"open"`
	DeviceType string `json:"deviceType"`
	DeviceName string `json:"deviceName"`
}
