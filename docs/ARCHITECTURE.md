# Architecture

## Components
- **Magisk module:** installer, lifecycle service, backend payload and opt-in boot Auto Start.
- **Go server:** loopback-only API, hotspot control, settings persistence, profiles and diagnostics.
- **Web dashboard:** responsive HTML, CSS and JavaScript; settings are saved as they change.
- **Capability registry:** read-only evidence; unknown capabilities are not presented as supported.

## Persistent data
All ENets runtime data is stored under `/data/local/enets`.
- `ap.conf`: VirtualAP configuration, including the hotspot passphrase; mode 0600.
- `settings.json`: saved dashboard settings and Auto Start preference; mode 0600.
- `profiles.json`: named profiles; mode 0600.
- `autostart`: one-line boot preference; mode 0600.

Writes use a temporary file, fsync and atomic rename. Saved settings are returned only to the loopback dashboard so the form can be restored. Diagnostic exports redact common password/passphrase fields.

## Boot behavior
The Magisk service starts the local dashboard after boot. It starts the hotspot only when the user has explicitly enabled Auto Start. Missing preference means OFF. Auto Start makes one attempt; it does not loop indefinitely.

## Security and compatibility
- Bind only to `127.0.0.1:8080`; reject non-loopback peers and untrusted Host headers.
- State-changing API calls require same-origin requests and root.
- Keep configuration files private (0600); parent runtime directory is 0700.
- WPA3 is not declared supported from generic radio frequency data. Enable only after device/firmware SoftAP SAE capability is positively validated.
- Unknown capabilities remain unavailable rather than guessed.

## API
- `GET /api/v1/health`
- `GET /api/v1/capabilities`
- `GET /api/v1/hotspot/status`
- `POST /api/v1/hotspot/start`
- `POST /api/v1/hotspot/stop`
- `GET /api/v1/hotspot/clients`
- `GET|PUT /api/v1/settings`
- `GET|POST /api/v1/profiles`
- `GET /api/v1/diagnostics`
- `GET /api/v1/diagnostics/logs` — recent redacted log excerpt
- `GET /api/v1/version`

The service is intentionally local-only. Do not expose it to LAN or hotspot clients.
