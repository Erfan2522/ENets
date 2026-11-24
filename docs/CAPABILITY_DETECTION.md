# Capability detection

The first Android integration step is deliberately read-only.

## Endpoint

When the Go server is running locally, request:

```http
GET /api/v1/capabilities
```

The response reports Android release/API level, primary ABI, device codename, Wi-Fi interface names found through standard sysfs markers, and capability states.

## Status meanings

- `not_integrated`: ENets has no tested adapter for this operation yet. It cannot be applied.
- `unknown`: the environment did not expose enough information to decide.
- `hardware_detected_needs_platform_validation`: the Wi-Fi driver reported frequencies in that band. This is not confirmation that Android hotspot/AP mode supports the band.
- Future adapters may report `supported` only after a device/firmware-specific test confirms the operation.

## Detection limits

The detector reads Android properties, standard sysfs paths, and (if available) `iw list`. Missing tools or permissions produce unknown results rather than guessed support. Frequency listings describe radio/driver capabilities; they do not establish regulatory approval, AP-mode support, tethering-framework support, or a working vendor implementation.

The endpoint never reads hotspot passwords, never changes settings, and never edits hostapd or vendor configuration files.

## Safety boundary

The server remains bound to `127.0.0.1:8080` by default. Do not expose it to hotspot clients until authentication, authorization, CSRF protection, and interface restrictions are implemented and reviewed.
