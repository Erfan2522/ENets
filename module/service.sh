#!/system/bin/sh
# ENets boot service placeholder.
# A production service must verify the binary, configuration permissions,
# authentication state, and supported network interface before starting.
# Do not expose a dashboard on a network interface in this scaffold.

MODDIR=${0%/*}
ENETS_BIN="$MODDIR/system/bin/enetsd"

[ -x "$ENETS_BIN" ] || exit 0
# Intentionally disabled until the authenticated server and Android adapter
# are implemented and tested.
exit 0
