//go:build linux

package systeminfo

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

func getInfoLinux() (Info, error) {
	var info Info
	tryRead := func(p string) string {
		b, _ := os.ReadFile(p)
		return strings.TrimSpace(string(b))
	}
	paths := map[string]*string{
		"/sys/class/dmi/id/product_serial":           &info.SerialNumber,
		"/sys/devices/virtual/dmi/id/product_serial": &info.SerialNumber,
		"/sys/class/dmi/id/product_uuid":             &info.UUID,
		"/sys/class/dmi/id/sys_vendor":               &info.Manufacturer,
		"/sys/class/dmi/id/product_name":             &info.Model,
	}
	for p, ptr := range paths {
		if _, err := os.Stat(p); err == nil {
			if v := tryRead(p); v != "" && !isFake(v) {
				*ptr = v
			}
		}
	}
	if !hasMeaningful(info) {
		if out, err := exec.Command("dmidecode", "-t", "1").CombinedOutput(); err == nil {
			s := string(out)
			if info.Manufacturer == "" {
				info.Manufacturer = firstMatchLinePrefix(s, "Manufacturer:")
			}
			if info.Model == "" {
				info.Model = firstMatchLinePrefix(s, "Product Name:")
			}
			if info.SerialNumber == "" {
				sn := firstMatchLinePrefix(s, "Serial Number:")
				if sn != "" && !isFake(sn) {
					info.SerialNumber = sn
				}
			}
			if info.UUID == "" {
				info.UUID = firstMatchLinePrefix(s, "UUID:")
			}
		}
	}
	if hasMeaningful(info) {
		return info, nil
	}
	return Info{}, errors.New("could not determine system info on linux")
}

func isFake(s string) bool {
	l := strings.ToLower(strings.TrimSpace(s))
	fakes := []string{"to be filled by o.e.m.", "none", "unknown", "system manufacturer", "system product name"}
	for _, f := range fakes {
		if strings.Contains(l, f) {
			return true
		}
	}
	return false
}

func firstMatchLinePrefix(s, prefix string) string {
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, prefix) {
			return strings.TrimSpace(strings.TrimPrefix(line, prefix))
		}
	}
	return ""
}
