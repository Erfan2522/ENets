# Android hotspot integration findings

## AOSP baseline

Android's Soft AP stack separates Wi-Fi access-point configuration from internet tethering.

The Android Open Source Project documentation describes configurable SSID, security, hidden SSID, band/channel, and maximum client count. It also explains that device capabilities determine which controls are available and that some APIs are system APIs restricted to privileged/system Settings access.

Official references:
- https://source.android.com/docs/core/connect/wifi-softap
- https://android.googlesource.com/platform/packages/modules/Wifi/+/refs/heads/android14-release/service/java/com/android/server/wifi/WifiShellCommand.java
- https://android.googlesource.com/platform/packages/modules/Connectivity/+/refs/heads/main/Tethering/common/TetheringLib/src/android/net/TetheringManager.java

## Important limitation of `cmd wifi start-softap`

AOSP's Wi-Fi shell command calls the Wi-Fi service's `startTetheredHotspot` method, but its own help text explicitly says this command does not activate internet tethering. Therefore it must not be treated as a complete replacement for Android's tethering flow.

The AOSP documentation describes the tethered flow as:
1. Set a `SoftApConfiguration` through `WifiManager#setSoftApConfiguration`.
2. Start internet sharing through `TetheringManager#startTethering`.

These APIs include system/privileged access restrictions. A root shell does not automatically make every framework API callable or stable across OEM builds.

## ENets implementation decision

- Do not use `cmd wifi start-softap` as the universal ENets hotspot backend.
- Do not use hard-coded Binder transaction IDs or undocumented `service call` numbers as a generic fallback.
- Do not edit hostapd/vendor configuration files blindly.
- Evaluate framework/tethering access separately for each tested Android/OEM family.
- A device may expose AP radio support while the ENets tethered-hotspot operation remains unknown or unsupported.
- If ENets can configure the AP but cannot confirm internet tethering, expose those as separate capabilities and never label the complete hotspot as working.

## Required validation

For each adapter/device/firmware tuple, validate:
- Reading the existing Soft AP configuration without exposing the passphrase.
- Applying a temporary SSID/security/band only after explicit user approval.
- Confirming the actual running AP state and configuration.
- Confirming internet tethering separately from AP startup.
- Restoring the original configuration on test failure.
- Reboot, upgrade, disable, and uninstall behavior.

Until these tests pass, write capabilities remain `unknown`.
