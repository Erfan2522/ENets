#!/system/bin/sh

MODDIR=${0%/*}
PID_FILE="$MODDIR/enets.pid"
LOG_FILE="$MODDIR/enets.log"

if [ -f "$PID_FILE" ]; then
  PID="$(cat "$PID_FILE" 2>/dev/null)"
  case "$PID" in
    ''|*[!0-9]*) ;;
    *) kill "$PID" 2>/dev/null ;;
  esac
  rm -f "$PID_FILE"
fi

# Preserve logs for troubleshooting; do not touch Android system settings.
[ -f "$LOG_FILE" ] && echo "[ENets] Module removed." >> "$LOG_FILE"
