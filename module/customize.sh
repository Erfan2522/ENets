#!/system/bin/sh

SKIPUNZIP=0

ui_print "*******************************"
ui_print "       ENets by Erfan"
ui_print "*******************************"
ui_print "- Installing ENets local dashboard"
ui_print "- Bundles the pinned VirtualAP engine and hostapd/dnsmasq"
ui_print "- Backend files are isolated under /data/local/enets"
ui_print "- Dashboard hotspot controls are not enabled in this preview"
ui_print "- The server listens on 127.0.0.1 only"
ABI="$(getprop ro.product.cpu.abi)"
if [ "$ABI" != "arm64-v8a" ]; then
  abort "This preview currently requires ARM64 (arm64-v8a). Detected: ${ABI:-unknown}"
fi
ui_print "- ARM64 device detected"
