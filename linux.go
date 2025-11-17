//go:build linux

package systeminfo

import (
	"os"
	"strings"
)

func read(path string) string {
	b, _ := os.ReadFile(path)
	return strings.TrimSpace(string(b))
}

func Get() (Info, error) {
	return Info{
		SerialNumber: read("/sys/class/dmi/id/product_serial"),
		UUID:         read("/sys/class/dmi/id/product_uuid"),
		Manufacturer: read("/sys/class/dmi/id/sys_vendor"),
		Model:        read("/sys/class/dmi/id/product_name"),
	}, nil
}
