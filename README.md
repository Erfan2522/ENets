# ENets by Erfan

ENets is an open-source Magisk module project for managing an Android Wi-Fi hotspot through a lightweight local web dashboard.

> **Status: early development (v0.1.0-alpha).** An installable ARM64 read-only preview package is being built by GitHub Actions. The preview does not change Android hotspot settings and is not the completed hotspot manager.

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
- ARM64 is the first planned target. The preview package requires `arm64-v8a`.
- The goal is cross-device support through modular Android/vendor adapters, not a build limited to one phone.
- Hotspot controls vary by Android release, vendor, chipset, and firmware. Root access alone does not guarantee that a setting can be changed.

**Graceful compatibility rule:** if a device does not support a particular feature, ENets disables only that feature and explains why. The rest of the application and unrelated supported features continue working. Unknown capabilities are not treated as supported.

The implementation must detect device capabilities and refuse unsupported operations instead of blindly editing vendor configuration files. See [docs/COMPATIBILITY.md](docs/COMPATIBILITY.md) and [docs/ADAPTER_CONTRACT.md](docs/ADAPTER_CONTRACT.md).

## Preview package

The **Build Magisk preview** GitHub Actions workflow creates an installable ZIP artifact for ARM64. See [docs/INSTALLATION.md](docs/INSTALLATION.md) for instructions.

The preview runs a local read-only dashboard at `http://127.0.0.1:8080`. It does not apply SSID, password, security, band, channel, or hotspot on/off changes. It does not discover connected clients yet. It binds to loopback only and must not be exposed to hotspot clients or the LAN.

## Development

The Go server exposes a read-only `GET /api/v1/capabilities` endpoint and a health endpoint. Capability values remain unknown until validated by a device/firmware adapter. The server uses loopback and rejects non-loopback peers and untrusted Host headers.

### Run tests and build locally

Install Go, then run:

```sh
go test ./...
go build -trimpath -o enetsd ./server
```

The server expects the web assets in a `web/` directory relative to its working directory, or in the directory named by `ENETS_WEB_DIR`.

## Security principles

- Never ship a default admin password.
- Never log hotspot passwords or other secrets.
- Store sensitive settings with restrictive file permissions.
- Require authentication before enabling remote dashboard access.
- Do not silently update or overwrite user settings.
- Stop safely when a critical operation fails.

## Roadmap

See [docs/ROADMAP.md](docs/ROADMAP.md), [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md), [docs/COMPATIBILITY.md](docs/COMPATIBILITY.md), and [docs/INSTALLATION.md](docs/INSTALLATION.md).

## License

GPL-3.0-only. See [LICENSE](LICENSE).
