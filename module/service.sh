#!/system/bin/sh

MODDIR=${0%/*}
ENETS_BIN="$MODDIR/system/bin/enetsd"
WEB_DIR="$MODDIR/web"
PAYLOAD_DIR="$MODDIR/payload/bin"
ENETS_DATA="/data/local/enets"
LOG_FILE="$MODDIR/enets.log"
PID_FILE="$MODDIR/enets.pid"

log() {
  echo "[ENets] $1" >> "$LOG_FILE"
}

is_enets_process() {
  CHECK_PID="$1"
  [ -r "/proc/$CHECK_PID/cmdline" ] || return 1
  tr '\000' ' ' < "/proc/$CHECK_PID/cmdline" 2>/dev/null | grep -F "$ENETS_BIN" >/dev/null 2>&1
}

# VirtualAP is built from a pinned upstream commit. Keep its runtime files in
# an ENets-owned directory so the original VirtualAP app cannot collide with it.
install_backend_payload() {
  [ -d "$PAYLOAD_DIR" ] || {
    log "VirtualAP payload is not bundled; starting dashboard without hotspot backend."
    return 0
  }

  for name in virtualap hostapd dnsmasq; do
    [ -s "$PAYLOAD_DIR/$name" ] || {
      log "Backend payload incomplete: missing $name."
      return 1
    }
  done

  mkdir -p "$ENETS_DATA/bin" "$ENETS_DATA/logs" "$ENETS_DATA/run" || return 1
  for name in virtualap hostapd dnsmasq; do
    cp -f "$PAYLOAD_DIR/$name" "$ENETS_DATA/bin/$name" || return 1
    chmod 0755 "$ENETS_DATA/bin/$name" || return 1
  done
  log "Pinned VirtualAP engine payload installed under $ENETS_DATA/bin."
  return 0
}

# The dashboard server remains loopback-only. Its same-origin API can control
# the staged engine; remote/LAN access is intentionally not supported.
ABI="$(getprop ro.product.cpu.abi)"
case "$ABI" in
  arm64-v8a) ;;
  *)
    log "Unsupported ABI: ${ABI:-unknown}; service not started."
    exit 0
    ;;
esac

[ -x "$ENETS_BIN" ] || { log "Server binary missing."; exit 0; }
[ -d "$WEB_DIR" ] || { log "Dashboard assets missing."; exit 0; }
install_backend_payload || {
  log "Backend payload setup failed; dashboard will not be started."
  exit 0
}

if [ -f "$PID_FILE" ]; then
  OLD_PID="$(cat "$PID_FILE" 2>/dev/null)"
  case "$OLD_PID" in
    ''|*[!0-9]*) rm -f "$PID_FILE" ;;
    *)
      if kill -0 "$OLD_PID" 2>/dev/null && is_enets_process "$OLD_PID"; then
        log "Existing ENets server process detected."
        exit 0
      fi
      rm -f "$PID_FILE"
      ;;
  esac
fi

chmod 0755 "$ENETS_BIN"
export ENETS_WEB_DIR="$WEB_DIR"
nohup "$ENETS_BIN" >> "$LOG_FILE" 2>&1 &
PID=$!
echo "$PID" > "$PID_FILE"
chmod 0600 "$PID_FILE" "$LOG_FILE" 2>/dev/null

sleep 2
if kill -0 "$PID" 2>/dev/null && is_enets_process "$PID"; then
  log "Read-only dashboard started at http://127.0.0.1:8080"
else
  rm -f "$PID_FILE"
  log "Server failed to start. Check enets.log."
fi
