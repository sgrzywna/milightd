# milightd

Simple web service utilizing [milight](https://github.com/sgrzywna/milight).

## Build

To build service run:

```bash
make
```

## Start service

To connect to the Mi-Light device running at `192.168.0.102` and port `5987` (the default Mi-Light port), and to listen to commands at port 8080:

```bash
./milightd -mihost 192.168.0.102 -miport 5987 -port 8080
```

To see all available command line switches run:

```bash
./milightd -h
```

## Control the light

Service accepts JSON data to control color, brightness and status of the light:

```json
{
  "color": "white",
  "brightness": 32,
  "switch": "on"
}
```

API is [documented](api/swagger.yaml) with Swagger specification.

All parameters are optional, for example to turn light off only `switch` parameter must be present.

## Examples

To turn white light on with brightness 64 (maximal brightness):

```go
light := milightdclient.Light{}
light.SetColor("white")
light.SetBrightness(64)
light.SetSwitch(true)

client := milightdclient.NewClient("http://127.0.0.1:8080")
err := client.SetLight(light)
```

cURL call:

```bash
curl -X POST "http://127.0.0.1:8080/api/v1/light" -H "accept: application/json" -H "Content-Type: application/json" -d "{ \"color\": \"white\", \"brightness\": 64, \"switch\": \"on\"}"
```

To turn nice green light on with brightness 32:

```go
light := milightdclient.Light{}
light.SetColor("green")
light.SetBrightness(32)
light.SetSwitch(true)

client := milightdclient.NewClient("http://127.0.0.1:8080")
err := client.SetLight(light)
```

cURL call:

```bash
curl -X POST "http://127.0.0.1:8080/api/v1/light" -H "accept: application/json" -H "Content-Type: application/json" -d "{ \"color\": \"green\", \"brightness\": 32, \"switch\": \"on\"}"
```

To turn light off:

```go
light := milightdclient.Light{}
light.SetSwitch(false)

client := milightdclient.NewClient("http://127.0.0.1:8080")
err := client.SetLight(light)
```

cURL call:
```bash
curl -X POST "http://127.0.0.1:8080/api/v1/light" -H "accept: application/json" -H "Content-Type: application/json" -d "{ \"switch\": \"off\"}"
```

## Use case

The project [statuslight](https://github.com/sgrzywna/statuslight) implements service that utilizes `milightd` to control the light.
