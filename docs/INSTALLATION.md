# Installing the ENets preview

> This is an installable **read-only preview**, not the completed hotspot manager. It does not change SSID, password, security mode, Wi-Fi band, channels, clients, firewall, or routing.

## Requirements

- Android device with ARM64 userspace (`arm64-v8a`).
- Magisk installed and working.
- Enough free storage for the module and its logs.
- The package must be downloaded from a successful GitHub Actions build.

## Install

1. Download `ENets-v0.1.0-alpha-preview.zip` from the latest successful **Build Magisk preview** workflow artifact.
2. Open Magisk.
3. Choose Modules → Install from storage.
4. Select the ZIP and wait for installation.
5. Reboot the device.
6. Open a browser on the same device and visit `http://127.0.0.1:8080`.

The server binds to loopback only. It is not available to other hotspot clients or LAN devices.

## Logs and troubleshooting

The module writes its log to:

`/data/adb/modules/enets/enets.log`

If the dashboard does not open, check that the module is enabled and inspect the log. The service deliberately does not attempt to modify hotspot settings.

## Uninstall

Disable or remove ENets from Magisk. The uninstall script stops the preview server if it is running. It does not alter Android Wi-Fi, hotspot, firewall, routing, or vendor files.

## Current limitations

- Hotspot settings are not implemented.
- Client discovery is not implemented.
- The dashboard is local and read-only.
- ARM64 only; unsupported ABIs are rejected.
- Device compatibility beyond the build/install path has not yet been verified.
