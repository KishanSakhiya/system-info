# systeminfo
Cross-platform Go module to retrieve Serial Number, UUID, Manufacturer, and Model.

## Usage
```go
import "github.com/KishanSakhiya/system-info"

func main() {
    info, _ := systeminfo.Get()
    println(info.SerialNumber)
}
```
