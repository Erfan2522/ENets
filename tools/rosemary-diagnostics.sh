#!/system/bin/sh
# ENets read-only diagnostic collector for Xiaomi/MediaTek rosemary.
# This script does not change hotspot, Wi-Fi, firewall, or system settings.
# It intentionally avoids serial numbers, IMEI, MAC addresses, SSIDs and passwords.

section() { printf '\n===== %s =====\n' "$1"; }

section "ENets rosemary diagnostic"
echo "mode=READ_ONLY"
echo "timestamp=$(date 2>/dev/null)"
echo "device=$(getprop ro.product.device)"
echo "manufacturer=$(getprop ro.product.manufacturer)"
echo "model=$(getprop ro.product.model)"
echo "android_release=$(getprop ro.build.version.release)"
echo "android_api=$(getprop ro.build.version.sdk)"
echo "build_fingerprint=$(getprop ro.build.fingerprint | sed 's/\/[^/]*user\/.*/\/user\/[redacted]/')"
echo "kernel=$(uname -r 2>/dev/null)"
echo "abi=$(getprop ro.product.cpu.abi)"
echo "magisk_version=$(magisk -V 2>/dev/null || echo unavailable)"

section "Wi-Fi interfaces (names only)"
found=0
for path in /sys/class/net/*; do
    [ -e "$path" ] || continue
    name=${path##*/}
    if [ -e "$path/wireless" ] || [ -e "$path/phy80211" ]; then
        echo "$name"
        found=1
    fi
done
[ "$found" -eq 1 ] || echo "No standard wireless sysfs marker found"

section "Available diagnostic tools"
for tool in iw ip cmd dumpsys; do
    if command -v "$tool" >/dev/null 2>&1; then
        echo "$tool=available"
    else
        echo "$tool=unavailable"
    fi
done

section "Wi-Fi radio capabilities (if iw is available)"
if command -v iw >/dev/null 2>&1; then
    iw list 2>&1 | awk '
        /Wiphy/ || /Band [0-9]+:/ || /Frequencies:/ || /MHz/ ||
        /HT Capabilities/ || /VHT Capabilities/ || /HE Iftypes/ ||
        /EHT Iftypes/ || /Supported interface modes:/ || /\* AP/ {
            print
        }
    ' | head -n 180
else
    echo "iw command unavailable; radio capability data not collected"
fi

section "Relevant Wi-Fi properties"
getprop | grep -iE '\[(wifi|wlan|vendor\.wifi|ro\.vendor\.wifi)[^]]*\]' |
    sed -E 's/(password|passwd|psk|ssid)([^]]*)\]:.*/\1\2]:[redacted]/I' |
    head -n 80

section "End"
echo "No settings were changed."
echo "Review the report before sharing it. Do not post private identifiers."
