#!/system/bin/sh

# Read-only discovery of Android Wi-Fi and tethering shell interfaces.
# This script prints command help and service names only; it does not dump
# active Wi-Fi state, SSIDs, MAC addresses, passwords, or device identifiers.

echo "===== ENets Android hotspot API check ====="
echo "mode=READ_ONLY"
echo "android_release=$(getprop ro.build.version.release)"
echo "android_api=$(getprop ro.build.version.sdk)"
echo "manufacturer=$(getprop ro.product.manufacturer)"
echo "device=$(getprop ro.product.device)"

for service_name in wifi tethering connectivity; do
  echo
  echo "===== cmd $service_name help (filtered) ====="
  if command -v cmd >/dev/null 2>&1; then
    cmd "$service_name" help 2>&1 | grep -iE 'soft.?ap|hotspot|tether|band|channel|ssid|security|passphrase|start|stop|config|wifi' | head -n 60
    if [ "${PIPESTATUS:-}" = "" ]; then :; fi
  else
    echo "cmd=unavailable"
  fi
done

echo
echo "===== Relevant Binder services (names only) ====="
service list 2>/dev/null | grep -iE 'wifi|tether|connectivity' | sed -E 's/^[[:space:]]*[0-9]+[[:space:]]+//'
echo
echo "===== End ====="
