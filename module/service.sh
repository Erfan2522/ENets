#!/system/bin/sh

MODDIR=${0%/*}
ENETS_BIN="$MODDIR/system/bin/enetsd"
WEB_DIR="$MODDIR/web"
PAYLOAD_DIR="$MODDIR/payload/bin"
ENETS_DATA="/data/local/enets"
LOG_FILE="$MODDIR/enets.log"
PID_FILE="$MODDIR/enets.pid"
WATCHDOG_PID_FILE="$MODDIR/enets-watchdog.pid"

log() { echo "[ENets] $1" >> "$LOG_FILE"; }

is_enets_process() {
  CHECK_PID="$1"
  [ -r "/proc/$CHECK_PID/cmdline" ] || return 1
  tr '\000' ' ' < "/proc/$CHECK_PID/cmdline" 2>/dev/null | grep -F "$ENETS_BIN" >/dev/null 2>&1
}

start_dashboard() {
  chmod 0755 "$ENETS_BIN"
  export ENETS_WEB_DIR="$WEB_DIR"
  nohup "$ENETS_BIN" >> "$LOG_FILE" 2>&1 &
  PID=$!
  echo "$PID" > "$PID_FILE"
  chmod 0600 "$PID_FILE" "$LOG_FILE" 2>/dev/null
  sleep 2
  kill -0 "$PID" 2>/dev/null && is_enets_process "$PID"
}

install_backend_payload() {
  [ -d "$PAYLOAD_DIR" ] || { log "VirtualAP payload missing."; return 0; }
  for name in virtualap hostapd dnsmasq; do
    [ -s "$PAYLOAD_DIR/$name" ] || { log "Backend payload incomplete: missing $name."; return 1; }
  done
  mkdir -p "$ENETS_DATA/bin" "$ENETS_DATA/logs" "$ENETS_DATA/run" || return 1
  chmod 0700 "$ENETS_DATA" "$ENETS_DATA/logs" "$ENETS_DATA/run" || return 1
  if [ -x "$ENETS_DATA/bin/virtualap" ]; then "$ENETS_DATA/bin/virtualap" stop >/dev/null 2>&1; fi
  for name in virtualap hostapd dnsmasq; do
    cp -f "$PAYLOAD_DIR/$name" "$ENETS_DATA/bin/$name" || return 1
    chmod 0755 "$ENETS_DATA/bin/$name" || return 1
  done
  log "Pinned VirtualAP engine payload installed under $ENETS_DATA/bin."
}

ABI="$(getprop ro.product.cpu.abi)"
case "$ABI" in armeabi-v7a) ;; *) log "Unsupported ABI: ${ABI:-unknown}; ARMv7 service not started."; exit 0 ;; esac
[ -x "$ENETS_BIN" ] || { log "Server binary missing."; exit 0; }
[ -d "$WEB_DIR" ] || { log "Dashboard assets missing."; exit 0; }
install_backend_payload || { log "Backend payload setup failed."; exit 0; }

if [ -f "$PID_FILE" ]; then
  OLD_PID="$(cat "$PID_FILE" 2>/dev/null)"
  case "$OLD_PID" in
    ''|*[!0-9]*) rm -f "$PID_FILE" ;;
    *) if kill -0 "$OLD_PID" 2>/dev/null && is_enets_process "$OLD_PID"; then log "Existing ENets server process detected."; exit 0; fi; rm -f "$PID_FILE" ;;
  esac
fi

if start_dashboard; then
  log "Local dashboard started at http://127.0.0.1:8080"
else
  rm -f "$PID_FILE"
  log "Server failed to start. Check enets.log."
  exit 0
fi

# Hotspot auto-start is opt-in. Missing preference means OFF.
if [ "$(cat "$ENETS_DATA/autostart" 2>/dev/null)" = "1" ]; then
  if [ -s "$ENETS_DATA/ap.conf" ] && [ -x "$ENETS_DATA/bin/virtualap" ]; then
    log "Auto Start is enabled; attempting one hotspot start."
    if "$ENETS_DATA/bin/virtualap" start >> "$LOG_FILE" 2>&1; then log "Auto Start command completed."; else log "Auto Start failed; inspect diagnostics."; fi
  else
    log "Auto Start enabled but saved hotspot configuration is unavailable."
  fi
fi

# Bounded watchdog: restart the dashboard at most three times per boot.
(
  RESTARTS=0
  while [ "$RESTARTS" -lt 3 ]; do
    sleep 20
    [ -x "$ENETS_BIN" ] || exit 0
    CURRENT_PID="$(cat "$PID_FILE" 2>/dev/null)"
    case "$CURRENT_PID" in
      ''|*[!0-9]*) CURRENT_PID=0 ;;
    esac
    if [ "$CURRENT_PID" -gt 1 ] && kill -0 "$CURRENT_PID" 2>/dev/null && is_enets_process "$CURRENT_PID"; then
      continue
    fi
    RESTARTS=$((RESTARTS + 1))
    log "Dashboard process exited; recovery attempt $RESTARTS of 3."
    if start_dashboard; then
      log "Dashboard recovery succeeded."
    else
      rm -f "$PID_FILE"
      log "Dashboard recovery attempt failed."
    fi
  done
  log "Dashboard watchdog reached its restart limit."
) >/dev/null 2>&1 &
echo "$!" > "$WATCHDOG_PID_FILE"
chmod 0600 "$WATCHDOG_PID_FILE" 2>/dev/null
