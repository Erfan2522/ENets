# Architecture

## Components

- **Magisk module:** lifecycle scripts and module metadata.
- **Go server:** local HTTP server and future authenticated API.
- **Web dashboard:** responsive HTML, CSS, and JavaScript.
- **Android adapter layer (planned):** capability detection and device-specific hotspot operations.
- **Capability registry (planned):** independent status and explanation for every feature.

## Cross-device design

The dashboard and core server must remain device-independent. Device-specific behavior belongs in adapters selected using verified platform and device information (for example, Android release, vendor, firmware family, and available system interfaces).

Adapters must expose the same capability-oriented contract to the core. Adding support for a new manufacturer or firmware family should not require rewriting the dashboard or core.

## Graceful feature-level compatibility

An unsupported feature must not disable the whole application. ENets disables only that feature's control and explains why; unrelated supported features remain available.

Each capability has its own status:
- `supported`: detected and validated on this device/firmware combination.
- `unsupported`: confirmed unavailable or unsafe to control.
- `unknown`: insufficient evidence; do not offer an operation yet.
- `error`: capability check or operation failed.

Do not treat radio/driver detection as proof that Android hotspot configuration supports the same feature. Do not report a setting as applied until the platform confirms it.

Adapter failures must be isolated. A failed operation must not corrupt unrelated settings or take down the core dashboard. Critical core failures may stop the affected service safely, but must not trigger broad system changes.

## Safety boundaries

The current server serves static assets only and listens on 127.0.0.1:8080 by default. The Magisk service script intentionally does not launch it. The installer intentionally aborts because this alpha is not install-ready.

The future adapter must be capability-driven. Generic dashboard choices must not be treated as proof that the device supports a feature.

## Planned API areas

- `GET /api/status` — service and device capability summary.
- `GET /api/hotspot` — current hotspot state, with secrets omitted.
- `PUT /api/hotspot` — validated settings request, requiring authentication and CSRF protection.
- `GET /api/clients` — read-only connected-client information.

These endpoints are design notes only; they are not implemented in this scaffold.
