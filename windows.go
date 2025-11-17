//go:build windows

package systeminfo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"unicode/utf16"

	"github.com/yusufpapurcu/wmi"
)

type winResult struct {
	IdentifyingNumber string
	UUID              string
	Vendor            string
	Name              string
}

func getInfoWindows() (Info, error) {
	// 1) PowerShell CIM
	if info, err := getByPowerShellCIM(); err == nil && hasMeaningful(info) {
		return info, nil
	}
	// 2) StackExchange/wmi (COM)
	if info, err := getByWMI(); err == nil && hasMeaningful(info) {
		return info, nil
	}
	// 3) wmic cmd if present
	if out, err := run("cmd", "/C", "wmic csproduct get Vendor,Name,IdentifyingNumber,UUID /value"); err == nil {
		info := parseKeyValueOutput(out)
		if hasMeaningful(info) {
			return info, nil
		}
	}
	// 4) registry fallback
	info := getByRegistry()
	if hasMeaningful(info) {
		return info, nil
	}
	return Info{}, errors.New("could not determine system info on windows")
}

func getByPowerShellCIM() (Info, error) {
	ps := `Try {
  $o = Get-CimInstance Win32_ComputerSystemProduct -ErrorAction Stop | Select-Object IdentifyingNumber,UUID,Vendor,Name
  $o | ConvertTo-Json -Compress
} Catch {
  $o = Get-WmiObject Win32_ComputerSystemProduct -ErrorAction Stop | Select-Object IdentifyingNumber,UUID,Vendor,Name
  $o | ConvertTo-Json -Compress
}`
	out, err := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", ps).Output()
	if err != nil {
		return Info{}, err
	}
	s := decodePossibleUtf16LE(out)
	var m map[string]interface{}
	if json.Unmarshal([]byte(s), &m) == nil {
		return Info{
			SerialNumber: strFromMap(m, "IdentifyingNumber"),
			UUID:         strFromMap(m, "UUID"),
			Manufacturer: strFromMap(m, "Vendor"),
			Model:        strFromMap(m, "Name"),
		}, nil
	}
	return parseKeyValueOutput(s), nil
}

func decodePossibleUtf16LE(b []byte) string {
	if len(b) >= 2 && b[0] == 0xFF && b[1] == 0xFE {
		return utf16LEToString(b[2:])
	}
	nulls := bytes.Count(b, []byte{0})
	if nulls*4 > len(b) {
		return utf16LEToString(b)
	}
	return string(b)
}

func utf16LEToString(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	if len(b)%2 != 0 {
		b = b[:len(b)-1]
	}
	u16 := make([]uint16, 0, len(b)/2)
	for i := 0; i < len(b); i += 2 {
		u := uint16(b[i]) | uint16(b[i+1])<<8
		u16 = append(u16, u)
	}
	r := utf16.Decode(u16)
	return string(r)
}

func strFromMap(m map[string]interface{}, k string) string {
	if v, ok := m[k]; ok {
		return strings.TrimSpace(fmt.Sprintf("%v", v))
	}
	if v, ok := m[strings.ToLower(k)]; ok {
		return strings.TrimSpace(fmt.Sprintf("%v", v))
	}
	return ""
}

func getByWMI() (Info, error) {
	var bios []winBIOS
	if err := wmi.Query("SELECT SerialNumber FROM Win32_BIOS", &bios); err == nil && len(bios) > 0 {
		return Info{SerialNumber: strings.TrimSpace(bios[0].SerialNumber)}, nil
	}
	var prods []winCSP
	if err := wmi.Query("SELECT IdentifyingNumber,UUID,Vendor,Name FROM Win32_ComputerSystemProduct", &prods); err == nil && len(prods) > 0 {
		p := prods[0]
		return Info{
			SerialNumber: strings.TrimSpace(p.IdentifyingNumber),
			UUID:         strings.TrimSpace(p.UUID),
			Manufacturer: strings.TrimSpace(p.Vendor),
			Model:        strings.TrimSpace(p.Name),
		}, nil
	}
	return Info{}, errors.New("wmi query failed or returned nothing")
}

func getByRegistry() Info {
	var info Info
	if out, err := run("cmd", "/C", `reg query "HKLM\\HARDWARE\\DESCRIPTION\\System\\BIOS" /v SystemSerialNumber`); err == nil {
		if v := parseRegValue(out); v != "" {
			info.SerialNumber = v
		}
	}
	if out, err := run("cmd", "/C", `reg query "HKLM\\HARDWARE\\DESCRIPTION\\System\\BIOS" /v SystemManufacturer`); err == nil {
		if v := parseRegValue(out); v != "" {
			info.Manufacturer = v
		}
	}
	if out, err := run("cmd", "/C", `reg query "HKLM\\HARDWARE\\DESCRIPTION\\System\\BIOS" /v SystemProductName`); err == nil {
		if v := parseRegValue(out); v != "" {
			info.Model = v
		}
	}
	return info
}

func parseRegValue(out string) string {
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		f := strings.Fields(line)
		if len(f) >= 3 && strings.HasPrefix(f[1], "REG_") {
			return strings.Join(f[2:], " ")
		}
	}
	return ""
}
