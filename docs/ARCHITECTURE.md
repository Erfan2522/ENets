# Architecture

## Components

- **Magisk module:** module metadata, installer, lifecycle service, and uninstall handling.
- **Go server:** loopback-only local HTTP server and read-only capability API.
- **Web dashboard:** responsive HTML, CSS, and JavaScript.
- **Android adapter layer (planned):** capability detection and device-specific hotspot operations.
- **Capability registry:** independent status and explanation for every feature.

## Cross-device design

The dashboard and core server must remain device-independent. Device-specific behavior belongs in adapters selected using verified platform and device information (for example, Android release, vendor, firmware family, and available system interfaces).

Adapters must expose the same capability-oriented contract to the core. Adding support for a new manufacturer or firmware family should not require rewriting the dashboard or core. See [ADAPTER_CONTRACT.md](ADAPTER_CONTRACT.md).

## Graceful feature-level compatibility

An unsupported feature must not disable the whole application. ENets disables only that feature's control and explains why; unrelated supported features remain available.

Each capability has its own status:
- `supported`: detected and validated on this device/firmware combination.
- `unsupported`: confirmed unavailable or unsafe to control.
- `unknown`: insufficient evidence; do not offer an operation yet.
- `error`: capability check or operation failed.

Do not treat radio/driver detection as proof that Android hotspot configuration supports the same feature. Do not report a setting as applied until the platform confirms it.

Adapter failures must be isolated. A failed operation must not corrupt unrelated settings or take down the core dashboard.

## Current preview boundaries

The installable preview:
- Starts the Go server through the Magisk service script on ARM64 devices.
- Binds only to `127.0.0.1:8080`.
- Accepts only loopback peers and localhost Host values.
- Serves the static dashboard and read-only capability endpoint.
- Does not change hotspot settings, firewall rules, routing, or vendor files.
- Does not yet implement authentication because there are no write/control endpoints and the server is strictly local/read-only. Authentication and authorization are mandatory before adding any state-changing or remote access feature.

The preview is not the completed hotspot manager.

## Planned API areas

- `GET /api/v1/health` — local server health.
- `GET /api/v1/capabilities` — read-only device capability report.
- Future `GET /api/v1/hotspot` — current hotspot state, with secrets omitted.
- Future `PUT /api/v1/hotspot` — validated settings request, requiring authentication and CSRF protection.
- Future `GET /api/v1/clients` — read-only connected-client information.

Future endpoints are design notes only; they are not implemented in this preview.
