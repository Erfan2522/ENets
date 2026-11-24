# ENets by Erfan

ENets is an open-source Magisk module project for managing an Android Wi-Fi hotspot through a lightweight local web dashboard.

> **Status: early development (v0.1.0-alpha scaffold).** This repository is not an install-ready module yet. The current dashboard is a UI prototype; it does not change Android hotspot settings.

## Planned features

- Configure hotspot SSID and password.
- Choose a supported security mode (WPA2/WPA3 where the device and firmware expose it).
- Choose a supported Wi-Fi band (Auto, 2.4 GHz, 5 GHz, and 6 GHz only when available).
- Read-only connected-client list.
- Minimal responsive web dashboard.
- Magisk service with start-on-boot support.
- Capability detection and clear unsupported states.

ENets will not manage firewall rules, block clients, or modify the system firewall. Traffic monitoring is not part of the initial release.

## Compatibility

- Magisk only in the first development phase.
- Android 8.0+ is the intended minimum; compatibility is not yet verified.
- ARM64 is the first planned target.
- Hotspot controls vary by Android release, vendor, chipset, and firmware. Root access alone does not guarantee that a setting can be changed.

The implementation must detect device capabilities and refuse unsupported operations instead of blindly editing vendor configuration files.

## Development

The project currently contains a web UI prototype and a Go server with a read-only capability detection endpoint. The server binds to `127.0.0.1:8080` by default for safety. Do not expose it to hotspot clients until authentication, authorization, CSRF protection, and network-interface restrictions are implemented and tested. See [docs/CAPABILITY_DETECTION.md](docs/CAPABILITY_DETECTION.md).

### Run the dashboard prototype

Open `web/index.html` in a browser. Applying settings is intentionally unavailable until a tested Android integration backend exists.

### Build the Go server

Install Go, then run:

```sh
go build -o enetsd ./server
```

The server expects the web assets in a `web/` directory relative to its working directory.

## Security principles

- Never ship a default admin password.
- Never log hotspot passwords or other secrets.
- Store sensitive settings with restrictive file permissions.
- Require authentication before enabling remote dashboard access.
- Do not silently update or overwrite user settings.
- Stop safely when a critical operation fails.

## Roadmap

See [docs/ROADMAP.md](docs/ROADMAP.md), [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md), and [docs/COMPATIBILITY.md](docs/COMPATIBILITY.md).

## License

GPL-3.0-only. See [LICENSE](LICENSE).
