package systeminfo

type Info struct {
	SerialNumber string `json:"serial_number"`
	UUID         string `json:"uuid"`
	Manufacturer string `json:"manufacturer"`
	Model        string `json:"model"`
}
