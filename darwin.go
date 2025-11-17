//go:build darwin

package systeminfo

import (
	"os/exec"
	"strings"
)

func parse(lines []string, key string) string {
	for _, l := range lines {
		if strings.HasPrefix(l, key) {
			parts := strings.SplitN(l, ":", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	return ""
}

func Get() (Info, error) {
	out, _ := exec.Command("system_profiler", "SPHardwareDataType").Output()
	lines := strings.Split(string(out), "")
	return Info{
		SerialNumber: parse(lines, "Serial Number"),
		UUID:         parse(lines, "Hardware UUID"),
		Manufacturer: "Apple",
		Model:        parse(lines, "Model Name"),
	}, nil
}
