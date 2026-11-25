#!/system/bin/sh

SKIPUNZIP=0

ui_print "*******************************"
ui_print "       ENets by Erfan"
ui_print "*******************************"
ui_print "- Installing ENets local dashboard"
ui_print "- Bundles the pinned VirtualAP engine and hostapd/dnsmasq"
ui_print "- Backend files are isolated under /data/local/enets"
ui_print "- Experimental hotspot start/stop controls are enabled"
ui_print "- The server listens on 127.0.0.1 only"
ABI="$(getprop ro.product.cpu.abi)"
case "$ABI" in
  armeabi-v7a) ui_print "- ARMv7 32-bit device detected" ;;
  *) abort "This preview requires 32-bit ARMv7 (armeabi-v7a). Detected: ${ABI:-unknown}" ;;
esac
