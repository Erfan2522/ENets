#!/system/bin/sh

SKIPUNZIP=0

ui_print "*******************************"
ui_print "       ENets by Erfan"
ui_print "*******************************"
ui_print "- Installing local dashboard preview"
ui_print "- Android hotspot settings are NOT modified"
ui_print "- Wi-Fi controls are not implemented in this build"
ui_print "- The server listens on 127.0.0.1 only"
ABI="$(getprop ro.product.cpu.abi)"
if [ "$ABI" != "arm64-v8a" ]; then
  abort "ENets preview currently requires ARM64 (arm64-v8a). Detected: ${ABI:-unknown}"
fi
ui_print "- ARM64 device detected"
