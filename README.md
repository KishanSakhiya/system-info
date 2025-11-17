# systeminfo

Cross-platform Go module to retrieve Serial Number, UUID, Manufacturer, and Model.

Features:
- Windows: PowerShell CIM → StackExchange/wmi COM → Registry fallback
- Linux: /sys/class/dmi/id → dmidecode fallback
- macOS: system_profiler → ioreg fallback
- CLI `systeminfo` under `cmd/systeminfo`
- Unit tests and GitHub Actions CI

Usage:

```
import "github.com/kishansakhiya/systeminfo/pkg/systeminfo"

info, err := systeminfo.Get()
if err != nil { ... }
fmt.Println(info)
```
