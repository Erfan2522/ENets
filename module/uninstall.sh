#!/system/bin/sh

MODDIR=${0%/*}
PID_FILE="$MODDIR/enets.pid"
LOG_FILE="$MODDIR/enets.log"
VIRTUALAP="/data/local/enets/bin/virtualap"

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
      kill -TERM "$PID" 2>/dev/null
      sleep 1
      kill -0 "$PID" 2>/dev/null && kill -KILL "$PID" 2>/dev/null
      ;;
  esac
  rm -f "$PID_FILE"
fi

# Preserve logs and /data/local/enets configuration for troubleshooting/reinstall.
[ -f "$LOG_FILE" ] && echo "[ENets] Module removed." >> "$LOG_FILE"
