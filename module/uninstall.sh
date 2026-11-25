#!/system/bin/sh

MODDIR=${0%/*}
ENETS_BIN="$MODDIR/system/bin/enetsd"
PID_FILE="$MODDIR/enets.pid"
LOG_FILE="$MODDIR/enets.log"
VIRTUALAP="/data/local/enets/bin/virtualap"

is_enets_process() {
  CHECK_PID="$1"
  [ -r "/proc/$CHECK_PID/cmdline" ] || return 1
  tr '\000' ' ' < "/proc/$CHECK_PID/cmdline" 2>/dev/null | grep -F "$ENETS_BIN" >/dev/null 2>&1
}

# Stop only an access point started by ENets. Keep its saved configuration and
# logs in place so reinstalling ENets does not silently discard user settings.
if [ -x "$VIRTUALAP" ]; then
  "$VIRTUALAP" stop >/dev/null 2>&1
fi

if [ -f "$PID_FILE" ]; then
  PID="$(cat "$PID_FILE" 2>/dev/null)"
  case "$PID" in
    ''|*[!0-9]*) ;;
    *)
      if [ "$PID" -gt 1 ] && is_enets_process "$PID"; then
        kill -TERM "$PID" 2>/dev/null
        sleep 1
        kill -0 "$PID" 2>/dev/null && is_enets_process "$PID" && kill -KILL "$PID" 2>/dev/null
      fi
      ;;
  esac
  rm -f "$PID_FILE"
fi

# Preserve logs and /data/local/enets configuration for troubleshooting/reinstall.
[ -f "$LOG_FILE" ] && echo "[ENets] Module removed." >> "$LOG_FILE"
