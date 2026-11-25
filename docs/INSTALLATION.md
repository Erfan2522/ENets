# Installing the ENets integration preview

> This is an experimental ARM64 build with VirtualAP-backed hotspot start/stop controls. It has not yet been validated on every device or firmware. Successful installation does not guarantee working Wi-Fi AP mode or client internet access.

## Requirements

- Android device with ARM64 userspace (`arm64-v8a`).
- Magisk installed and working.
- A compatible Wi-Fi driver/firmware with AP mode and the required kernel networking support.
- Enough free storage for the module, backend binaries, logs, and corresponding source archive.
- Download the ZIP from a successful GitHub Actions run of the `integration/virtualap-backend` branch.

## Install

1. Open the [ENets GitHub Actions workflow](https://github.com/Erfan2522/ENets/actions/workflows/build-magisk.yml).
2. Open the latest successful run for `integration/virtualap-backend`.
3. Download the `ENets-ARM64-integration-preview` artifact and extract its ZIP.
4. In Magisk, choose Modules → Install from storage.
5. Select `ENets-ARM64-integration-preview.zip` and wait for installation.
6. Reboot the device.
7. Open a browser on the same device and visit `http://127.0.0.1:8080`.

The server binds to loopback only. It is not available to other hotspot clients or LAN devices. Do not reverse-proxy or otherwise expose it remotely.

## First test

1. Start with a WPA2 password of 8–63 bytes and 2.4 GHz on Auto channel.
2. Tap **Apply & start hotspot** once. Do not repeatedly tap while the engine is starting.
3. Confirm the dashboard reports **Running**.
4. From a second device, check whether the SSID appears and whether it can connect.
5. Test internet access separately; AP startup alone does not prove tethering/NAT works.
6. Test **Stop hotspot** and confirm Wi-Fi returns to its prior state.
7. If 2.4 GHz works, test 5 GHz and other security modes separately.

Do not test WPA3, manual channels, or 80 MHz until the basic 2.4 GHz start/stop and client connectivity tests pass.

## Logs and runtime files

- Module log: `/data/adb/modules/enets/enets.log`
- Engine data/config/logs: `/data/local/enets/`
- Engine binaries: `/data/local/enets/bin/`

The saved hotspot configuration is stored root-only at `/data/local/enets/ap.conf`. The passphrase is not returned by the dashboard API and is not intentionally written to the module log.

## Troubleshooting and recovery

If the hotspot fails to start, stop the operation, check the ENets module and engine logs, and disable/remove the module from Magisk if necessary. Do not manually edit vendor Wi-Fi configuration or force a regulatory country.

The engine currently uses an AP interface named `ap0` and may need a device-specific adapter on other firmware. A successful start command is not sufficient proof that client DHCP, DNS, NAT, or internet access works.

## Uninstall

Disable or remove ENets from Magisk. The uninstall script attempts to stop only the ENets-managed access point and stops the local dashboard server. It preserves `/data/local/enets` configuration and logs for troubleshooting/reinstall; remove that directory manually only if you intentionally want to erase ENets settings.

## Current limitations

- ARM64 only.
- 2.4 GHz and 5 GHz only; no 6 GHz support.
- No maximum-client limit control.
- No client blocking, firewall editor, traffic accounting, or VPN routing.
- Cross-device compatibility and client internet sharing are not yet validated.
