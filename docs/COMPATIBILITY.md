# Compatibility policy

The intended minimum is Android 8.0+ with Magisk, and ARM64 is the first target. This is a target, not a verified compatibility claim.

Hotspot configuration is vendor- and release-specific. A rooted shell does not provide a universal supported interface for changing SSID, passphrase, security mode, channel, or band.

## Core rule: graceful feature-level compatibility

A device incompatibility must disable only the affected capability, not the entire ENets application. The dashboard, server, and unrelated supported features should continue working whenever their own requirements are met.

Each capability must have an independent status:
- `supported`: detected and validated on this device/firmware combination.
- `unsupported`: confirmed unavailable or not safely controllable.
- `unknown`: not enough evidence to make a safe claim.
- `error`: a check or operation failed; this is not equivalent to unsupported.

Unsupported or unknown controls must be disabled individually and accompanied by a clear explanation. ENets must not silently substitute a different setting or claim success without platform confirmation.

## Adapter requirements

ENets must:
1. Detect Android version, device ABI, Wi-Fi chipset/driver clues, and available system interfaces.
2. Select a narrowly scoped adapter for a tested device/firmware family.
3. Keep capability detection and operation results separate for each feature.
4. Isolate vendor-specific code so one adapter failure does not crash unrelated ENets components.
5. Report each unsupported or unknown capability explicitly.
6. Avoid blind edits to hostapd, vendor, or system configuration files.
7. Never claim a setting was applied unless the platform confirms it.
8. Fail safely: stop the failed operation, preserve unrelated settings, and keep other working features available.

WPA3, 5 GHz, and 6 GHz depend on hardware, regulatory domain, Android release, and vendor implementation. Driver or radio detection alone does not prove that Android hotspot mode can use a band.
