# Testing ENets on Xiaomi Poco M5s (rosemary)

This is the first planned reference device. The device codename is expected to be `rosemary`, but the diagnostic output must confirm the currently installed ROM, Android API level, kernel, ABI, and available Wi-Fi tooling.

## Safety

The diagnostic script is read-only. It does not toggle Wi-Fi/hotspot, edit configuration, change firewall rules, or read hotspot credentials. It avoids collecting IMEI, serial number, MAC address, SSID, and passwords.

Do not assume that a band listed by the Wi-Fi driver can be used by Android hotspot mode. AP-mode support, regulatory rules, vendor HAL, tethering framework, and firmware must all agree.

## Run the diagnostic

From a computer with ADB connected to the phone:

```sh
adb push tools/rosemary-diagnostics.sh /data/local/tmp/enets-rosemary-check.sh
adb shell chmod 0755 /data/local/tmp/enets-rosemary-check.sh
adb shell sh /data/local/tmp/enets-rosemary-check.sh
```

If `iw` is unavailable, the script reports that fact rather than installing tools or changing the device. Do not install random Wi-Fi utilities just for this check.

## What to share

Share the text output after reviewing it. Remove any private identifiers if a ROM build unexpectedly includes them. The report should contain Android version/API, kernel, ABI, interface names, tool availability, and filtered radio capability lines.

## Next implementation gate

Only after reviewing the report should ENets select a device-specific adapter strategy. No SSID/password/security/band write operation is enabled until a reversible implementation is verified on this exact device and firmware.
