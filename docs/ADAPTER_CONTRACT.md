# Android adapter contract

This document defines how ENets will support different Android devices without coupling the dashboard or core server to one manufacturer.

## Design goals

- One dashboard and one core capability model for all devices.
- Vendor- and firmware-specific behavior isolated in adapters.
- Feature-level fallback: an unavailable feature never disables unrelated working features.
- Read-only probing before any configuration operation.
- No undocumented configuration-file edits as a generic fallback.

## Adapter lifecycle

1. **Identify** — read Android API/release, manufacturer, model, device codename, ABI, and available system services. Do not collect serial numbers, IMEI, MAC addresses, SSIDs, or hotspot passwords.
2. **Match** — choose the most specific adapter whose match rules are known. A manufacturer name alone is not enough to claim support.
3. **Probe** — query each feature independently without changing state.
4. **Expose** — return a capability record for every feature.
5. **Operate** — accept only validated requests for capabilities marked supported on the exact device/firmware combination.
6. **Verify** — read the setting back from Android and report success only when the platform confirms it.
7. **Recover** — if an operation fails, stop that operation, preserve unrelated settings, and keep other features running.

## Capability status values

- `supported`: this exact adapter/feature combination has passed validation.
- `unsupported`: the feature is confirmed unavailable or unsafe on this device.
- `unknown`: evidence is incomplete; do not offer a write operation.
- `error`: the capability probe or operation failed. Keep this distinct from unsupported.

A radio frequency listed by `iw` is evidence about the radio/driver only. It does not mark a hotspot band as supported. The same distinction applies to WPA3, channel selection, hidden SSID, and maximum clients.

## Required adapter interface

Every adapter must provide equivalent operations, even if its implementation differs:

- `Match(platform) -> confidence and adapter identity`
- `Probe() -> per-feature capability records`
- `ReadHotspot() -> state without secrets`
- `Apply(request) -> operation result`
- `Verify(request) -> confirmed applied values`

The concrete implementation may use Android framework services, documented shell interfaces, or a tested vendor API. It must not assume that root access makes a private API stable.

## Initial adapter families

The first development phase will evaluate:
- AOSP / Pixel-like Android
- Xiaomi / Redmi
- Samsung One UI
- Other OEMs through later adapters

These are planned adapter families, not compatibility claims. Each supported device/firmware combination must be recorded in the compatibility matrix.

## Safety requirements

- No firewall, iptables, nftables, routing, or client-blocking operations.
- Never log credentials or hotspot secrets.
- No automatic change to hotspot settings during detection.
- No silent fallback to a different band/security mode.
- Do not expose the control API to LAN/hotspot clients until authentication, authorization, CSRF protection, and interface restrictions have been implemented and reviewed.
