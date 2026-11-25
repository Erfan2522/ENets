# Capability detection

The first Android integration step is deliberately read-only.

## Endpoint

When the Go server is running locally, request:

```http
GET /api/v1/capabilities
```

The response reports Android release/API level, primary ABI, device codename, Wi-Fi interface names, and independently detected Wi-Fi/tethering shell command help.

## Status meanings

- `supported`: this exact adapter/feature combination has passed device and firmware validation.
- `unsupported`: the feature is confirmed unavailable or unsafe on this device.
- `unknown`: evidence is incomplete; ENets must not offer a write operation.
- `error`: the capability probe or operation failed.

A shell command appearing in `cmd wifi help` is only command-discovery evidence. It does not prove the command is permitted, that the Soft AP starts, or that clients receive internet access.

## Separate capabilities

ENets tracks Soft AP startup and internet tethering separately. AOSP documents that `cmd wifi start-softap` does not activate internet tethering. Therefore ENets must never report a fully working hotspot based only on Soft AP command availability.

Official reference: https://android.googlesource.com/platform/packages/modules/Wifi/+/refs/heads/android14-release/service/java/com/android/server/wifi/WifiShellCommand.java

## Detection limits

The detector reads Android properties, standard sysfs paths, command help output, and (if available) `iw list`. Missing tools or permissions produce unknown results rather than guessed support. Frequency listings describe radio/driver capabilities; they do not establish regulatory approval, AP-mode support, tethering-framework support, or a working vendor implementation.

The endpoint never reads hotspot passwords, never changes settings, and never edits hostapd or vendor configuration files.

## Safety boundary

The server remains bound to `127.0.0.1:8080` by default. Do not expose it to hotspot clients until authentication, authorization, CSRF protection, and interface restrictions are implemented and reviewed.
