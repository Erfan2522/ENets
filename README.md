# ENets by Erfan

ENets is an open-source Magisk module for managing a rooted Android Wi-Fi hotspot through a lightweight local web dashboard.

> **Status: v0.2.0-alpha integration preview.** The integration branch packages a pinned VirtualAP engine and provides experimental local start/stop controls. Device and firmware compatibility, especially internet sharing, still require real-device testing.

## Current preview features

- Start and stop the access point.
- Configure SSID and WPA2, WPA3, WPA2/WPA3 transition, or open security.
- Select 2.4 GHz or 5 GHz, with optional channel and channel width.
- Hide the SSID.
- Read hotspot status and a read-only connected-client list.
- Local responsive dashboard served only on `127.0.0.1:8080`.
- Magisk service stages the engine and its hostapd/dnsmasq dependencies in an ENets-owned data directory.

The current backend does not expose client blocking, firewall configuration, traffic accounting, VPN routing, or a maximum-client limit. 6 GHz is not supported by this integration.

## Compatibility

- Magisk only in the first development phase.
- ARM64 (`arm64-v8a`) preview package.
- Android 8.0+ is an intended target, not a verified minimum for this backend.
- Device support depends on the Wi-Fi driver, AP interface, regulatory domain, kernel networking features, and firmware.
- The engine currently expects a usable physical Wi-Fi interface and creates/uses `ap0`; other interface layouts may need a device adapter.
- A successful engine start does not by itself prove that clients can connect or access the internet.

ENets does not force a Wi-Fi country code. The integration patch removes the upstream hardcoded US regulatory-domain lines so the device/kernel regulatory domain remains authoritative. Channel availability must still be tested on the device.

## Installable preview

The **Build ENets ARM64 Magisk preview** GitHub Actions workflow creates an installable ZIP artifact and a separate corresponding-source archive. Open the [integration branch Actions](https://github.com/Erfan2522/ENets/actions/workflows/build-magisk.yml) and download the latest successful artifact.

Install the ZIP from Magisk and reboot. Open `http://127.0.0.1:8080` locally on the device. Do not expose this server to the LAN or hotspot clients.

This is an experimental build. Test start, stop, reconnect, reboot, and client internet access before relying on it. Keep a way to disable/remove the Magisk module if the driver or firmware behaves unexpectedly.

## Backend and upstream attribution

The integration is based on [VirtualAP](https://github.com/ravindu644/VirtualAP), pinned at commit `19718afbf5df48f12d82a8a30bf205cf7bfa6c6e`. The build applies two ENets-specific patches:

1. Moves the runtime data directory from `/data/local/virtualap` to `/data/local/enets` to avoid colliding with the original VirtualAP app.
2. Removes the hardcoded US country setting and 802.11d line so the device/kernel regulatory domain is not overridden.

The engine source is GPL-3.0-or-later. ENets is GPL-3.0-only. Upstream copyright and license notices are retained; the workflow publishes the corresponding source, including the pinned hostapd and dnsmasq submodules. See [docs/UPSTREAM_VIRTUALAP.md](docs/UPSTREAM_VIRTUALAP.md).

## Development

The Go server exposes:
- `GET /api/v1/health`
- `GET /api/v1/capabilities`
- `GET /api/v1/hotspot/status`
- `POST /api/v1/hotspot/start`
- `POST /api/v1/hotspot/stop`
- `GET /api/v1/hotspot/clients`

State-changing requests require a same-origin request and a root-running service. The server binds to loopback only, rejects non-loopback peers and untrusted Host headers, and does not return the stored hotspot passphrase.

Run tests and build the server locally:

```sh
go test ./...
go build -trimpath -o enetsd ./server
```

## Security principles

- Never ship a default hotspot password.
- Never log or return hotspot passwords.
- Store the engine configuration with root-only permissions.
- Keep the dashboard loopback-only; do not enable remote access without authentication and CSRF protections.
- Do not force a regulatory country or blindly edit vendor configuration.
- Preserve upstream notices and provide corresponding source for distributed binaries.

## Roadmap

See [docs/ROADMAP.md](docs/ROADMAP.md), [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md), [docs/COMPATIBILITY.md](docs/COMPATIBILITY.md), and [docs/INSTALLATION.md](docs/INSTALLATION.md).

## License

ENets is GPL-3.0-only. See [LICENSE](LICENSE).
