# Architecture

## Components

- **Magisk module:** lifecycle scripts and module metadata.
- **Go server:** local HTTP server and future authenticated API.
- **Web dashboard:** responsive HTML, CSS, and JavaScript.
- **Android adapter layer (planned):** capability detection and device-specific hotspot operations.

## Safety boundaries

The current server serves static assets only and listens on 127.0.0.1:8080 by default. The Magisk service script intentionally does not launch it. The installer intentionally aborts because this alpha is not install-ready.

The future adapter must be capability-driven. Generic dashboard choices must not be treated as proof that the device supports a feature.

## Planned API areas

- `GET /api/status` — service and device capability summary.
- `GET /api/hotspot` — current hotspot state, with secrets omitted.
- `PUT /api/hotspot` — validated settings request, requiring authentication and CSRF protection.
- `GET /api/clients` — read-only connected-client information.

These endpoints are design notes only; they are not implemented in this scaffold.
