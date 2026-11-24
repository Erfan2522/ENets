# Compatibility policy

The intended minimum is Android 8.0+ with Magisk, and ARM64 is the first target. This is a target, not a verified compatibility claim.

Hotspot configuration is vendor- and release-specific. A rooted shell does not provide a universal supported interface for changing SSID, passphrase, security mode, channel, or band.

ENets must:
1. Detect Android version, device ABI, Wi-Fi chipset/driver clues, and available system interfaces.
2. Use a narrowly scoped adapter for a tested device/firmware family.
3. Report each unsupported capability explicitly.
4. Avoid blind edits to hostapd, vendor, or system configuration files.
5. Never claim a setting was applied unless the platform confirms it.

WPA3, 5 GHz, and 6 GHz depend on hardware, regulatory domain, Android release, and vendor implementation.
