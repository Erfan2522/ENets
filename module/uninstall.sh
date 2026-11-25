#!/system/bin/sh

MODDIR=${0%/*}
ENETS_BIN="$MODDIR/system/bin/enetsd"
PID_FILE="$MODDIR/enets.pid"
WATCHDOG_PID_FILE="$MODDIR/enets-watchdog.pid"
LOG_FILE="$MODDIR/enets.log"
VIRTUALAP="/data/local/enets/bin/virtualap"

is_enets_process() {
  CHECK_PID="$1"
  [ -r "/proc/$CHECK_PID/cmdline" ] || return 1
  tr '\000' ' ' < "/proc/$CHECK_PID/cmdline" 2>/dev/null | grep -F "$ENETS_BIN" >/dev/null 2>&1
}

if [ -f "$WATCHDOG_PID_FILE" ]; then
  WATCHDOG_PID="$(cat "$WATCHDOG_PID_FILE" 2>/dev/null)"
  case "$WATCHDOG_PID" in ''|*[!0-9]*) ;; *) [ "$WATCHDOG_PID" -gt 1 ] && kill "$WATCHDOG_PID" 2>/dev/null ;; esac
  rm -f "$WATCHDOG_PID_FILE"
fi

if [ -x "$VIRTUALAP" ]; then "$VIRTUALAP" stop >/dev/null 2>&1; fi

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

# Preserve settings, profiles, and logs for reinstall/recovery.
[ -f "$LOG_FILE" ] && echo "[ENets] Module removed." >> "$LOG_FILE"
