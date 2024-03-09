package model

type DeviceStatus struct {
	UDID       string `json:"udid"`
	MacAddress string `json:"mac_address"`
	Online     bool   `json:"online"`
}
