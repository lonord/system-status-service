# system-status-service
Simple system status service, including cpu, mem and disk info.

## Install

#### Download directly

[release](https://github.com/lonord/system-status-service/releases)

#### Build locally (go environment is required)

Clone this repo

```bash
git clone https://github.com/lonord/system-status-service
```

Build

```bash
./build.sh

# Show more build options
./build.sh -h
```

## Usage

```bash
Usage of system-status-service:
  -host string
        service listen host (default "0.0.0.0")
  -port int
        service listen port (default 2020)
  -v    show version
```

Run this binary and request the service

```bash
curl http://localhost:2020/system
```

JSON data received

```json
{
	"cpu": {
		"usage": [
			2.2471910109557927,
			0.9999999993742675,
			0,
			0
		],
		"temp": 51002
	},
	"memory": {
		"memoryUsedPercent": 16.411380306023695,
		"swapUsedPercent": 0
	},
	"disk": {
		"total": 30711148544,
		"used": 9860124672
	}
}
```

## License

MIT
