package systeminfo

// Info holds hardware identification fields.
type Info struct {
	SerialNumber string `json:"serial_number"`
	UUID         string `json:"uuid"`
	Manufacturer string `json:"manufacturer"`
	Model        string `json:"model"`
}

// WMI structs
type winBIOS struct {
	SerialNumber string
}
type winCSP struct {
	IdentifyingNumber string
	UUID              string
	Vendor            string
	Name              string
}
