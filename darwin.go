//go:build darwin

package systeminfo

import (
	"errors"
	"os/exec"
	"strings"
)

func getInfoDarwin() (Info, error) {
	var info Info
	if out, err := exec.Command("system_profiler", "SPHardwareDataType").CombinedOutput(); err == nil {
		for _, line := range strings.Split(string(out), "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "Serial Number") {
				if parts := strings.SplitN(line, ":", 2); len(parts) == 2 {
					info.SerialNumber = strings.TrimSpace(parts[1])
				}
			}
			if strings.HasPrefix(line, "Hardware UUID") {
				if parts := strings.SplitN(line, ":", 2); len(parts) == 2 {
					info.UUID = strings.TrimSpace(parts[1])
				}
			}
			if strings.HasPrefix(line, "Model Name") {
				if parts := strings.SplitN(line, ":", 2); len(parts) == 2 {
					info.Model = strings.TrimSpace(parts[1])
				}
			}
			if strings.HasPrefix(line, "Model Identifier") && info.Model == "" {
				if parts := strings.SplitN(line, ":", 2); len(parts) == 2 {
					info.Model = strings.TrimSpace(parts[1])
				}
			}
		}
		if hasMeaningful(info) {
			return info, nil
		}
	}
	// fallback ioreg
	if out, err := exec.Command("ioreg", "-c", "IOPlatformExpertDevice", "-r", "-d", "1").CombinedOutput(); err == nil {
		for _, line := range strings.Split(string(out), "\n") {
			line = strings.TrimSpace(line)
			if strings.Contains(line, "IOPlatformSerialNumber") {
				if parts := strings.Split(line, "="); len(parts) >= 2 {
					info.SerialNumber = strings.Trim(strings.TrimSpace(parts[1]), `"'`)
				}
			}
			if strings.Contains(line, "IOPlatformUUID") {
				if parts := strings.Split(line, "="); len(parts) >= 2 {
					info.UUID = strings.Trim(strings.TrimSpace(parts[1]), `"'`)
				}
			}
		}
	}
	if hasMeaningful(info) {
		return info, nil
	}
	return Info{}, errors.New("could not determine system info on darwin")
}
