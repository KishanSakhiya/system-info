package systeminfo

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os/exec"
	"runtime"
	"strings"
)

// Get returns best-effort hardware information for the current platform.
func Get() (Info, error) {
	switch runtime.GOOS {
	case "windows":
		return getInfoWindows()
	case "linux":
		return getInfoLinux()
	case "darwin":
		return getInfoDarwin()
	default:
		return Info{}, errors.New("unsupported OS: " + runtime.GOOS)
	}
}

// Fingerprint returns a sha256 hex string of the concatenated cleaned fields.
func Fingerprint(i Info) string {
	join := strings.ToLower(strings.TrimSpace(i.SerialNumber)) + "|" +
		strings.ToLower(strings.TrimSpace(i.UUID)) + "|" +
		strings.ToLower(strings.TrimSpace(i.Manufacturer)) + "|" +
		strings.ToLower(strings.TrimSpace(i.Model))
	h := sha256.Sum256([]byte(join))
	return hex.EncodeToString(h[:])
}

// helper: basic cleaning
func clean(s string) string {
	return strings.Trim(strings.TrimSpace(s), `"'\n\r `)
}

// helper to run commands
func run(cmd string, args ...string) (string, error) {
	out, err := exec.Command(cmd, args...).CombinedOutput()
	return string(out), err
}

// simple parser for key=value style output
func parseKeyValueOutput(s string) Info {
	var info Info
	for _, l := range strings.Split(s, "\n") {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}
		if strings.HasPrefix(strings.ToLower(l), "vendor=") {
			info.Manufacturer = clean(l[len("vendor="):])
		}
		if strings.HasPrefix(strings.ToLower(l), "name=") {
			info.Model = clean(l[len("name="):])
		}
		if strings.HasPrefix(strings.ToLower(l), "identifyingnumber=") {
			parts := strings.SplitN(l, "=", 2)
			if len(parts) == 2 {
				info.SerialNumber = clean(parts[1])
			}
		}
		if strings.HasPrefix(strings.ToLower(l), "uuid=") {
			info.UUID = clean(l[len("uuid="):])
		}
	}
	return info
}

func hasMeaningful(i Info) bool {
	return strings.TrimSpace(i.SerialNumber) != "" || strings.TrimSpace(i.Model) != ""
}
