#!/system/bin/sh

MODDIR=${0%/*}
ENETS_BIN="$MODDIR/system/bin/enetsd"
WEB_DIR="$MODDIR/web"
LOG_FILE="$MODDIR/enets.log"
PID_FILE="$MODDIR/enets.pid"

log() {
  echo "[ENets] $1" >> "$LOG_FILE"
}

# This package currently contains a read-only dashboard preview.
# Never change hotspot, firewall, routing, or vendor configuration.
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

if [ -f "$PID_FILE" ]; then
  OLD_PID="$(cat "$PID_FILE" 2>/dev/null)"
  case "$OLD_PID" in
    ''|*[!0-9]*) ;;
    *) kill -0 "$OLD_PID" 2>/dev/null && { log "Existing server process detected."; exit 0; } ;;
  esac
fi

chmod 0755 "$ENETS_BIN"
export ENETS_WEB_DIR="$WEB_DIR"
nohup "$ENETS_BIN" >> "$LOG_FILE" 2>&1 &
PID=$!
echo "$PID" > "$PID_FILE"
chmod 0600 "$PID_FILE" "$LOG_FILE" 2>/dev/null

sleep 2
if kill -0 "$PID" 2>/dev/null; then
  log "Read-only dashboard started at http://127.0.0.1:8080"
else
  rm -f "$PID_FILE"
  log "Server failed to start. Check enets.log."
fi
