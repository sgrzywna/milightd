# milight

Small go module to control Mi-Light iBox device (v6.0 protocol).

Inspired by the [node-milight-promise](https://github.com/mwittig/node-milight-promise). For protocol details see the [LimitlessLED-DevAPI](https://github.com/Fantasmos/LimitlessLED-DevAPI). Kudos to authors of these repositories.

Tested only with Mi-Light WiFi iBox Smart Light (Model No. iBox1), but should work with other compatible devices too.

## Usage

Minimal set of steps to initialize module and next turn on/off light with default (last set) color. For the full example see [examples/main.go](examples/main.go)

```go
import "github.com/sgrzywna/milight"

// Initialize module.
m, err := milight.NewMilight("192.168.0.102", 5987)
if err != nil {
    fmt.Printf("milight error: %s\n", err)
    os.Exit(1)
}
defer m.Close()

// Initialize communication session.
err = m.InitSession()
if err != nil {
    fmt.Printf("milight session error: %s\n", err)
    os.Exit(1)
}

err = m.On()
if err != nil {
    fmt.Printf("milight on error: %s\n", err)
    os.Exit(1)
}

time.Sleep(1 * time.Second)

err = m.Off()
if err != nil {
    fmt.Printf("milight off error: %s\n", err)
    os.Exit(1)
}
```
