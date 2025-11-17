//go:build windows

package systeminfo

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"strings"
)

type winResult struct {
	IdentifyingNumber string
	UUID              string
	Vendor            string
	Name              string
}

func Get() (Info, error) {
	// PowerShell CIM
	cmd := exec.Command("powershell", "-NoProfile", "-Command",
		`Get-CimInstance Win32_ComputerSystemProduct |
         Select-Object IdentifyingNumber,UUID,Vendor,Name |
         ConvertTo-Json`)
	out, err := cmd.Output()
	if err == nil {
		out = bytes.TrimPrefix(out, []byte{0xFF, 0xFE})
		var r winResult
		if json.Unmarshal(out, &r) == nil {
			return Info{
				SerialNumber: strings.TrimSpace(r.IdentifyingNumber),
				UUID:         strings.TrimSpace(r.UUID),
				Manufacturer: strings.TrimSpace(r.Vendor),
				Model:        strings.TrimSpace(r.Name),
			}, nil
		}
	}
	// fallback WMI COM skipped for brevity
	return Info{}, nil
}
