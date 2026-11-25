# VirtualAP backend evaluation and integration notes

## Selected upstream

- Project: [VirtualAP](https://github.com/ravindu644/VirtualAP)
- Upstream commit pinned for this evaluation: `19718afbf5df48f12d82a8a30bf205cf7bfa6c6e`
- Upstream repository license: GNU GPL-3.0
- C engine source headers: GPL-3.0-or-later
- ENets license: GPL-3.0-only. GPLv3 code may be incorporated into this GPLv3-only project, subject to preserving notices and complying with GPLv3.

### Bundled upstream components

The upstream repository uses two Git submodules:

| Component | Pinned revision | License reported by upstream |
|---|---|---|
| dnsmasq-vap | `02edbd5d1341112902543671a231d47f374a2778` | GPL-2.0 |
| hostapd-vap | `922f5ed2e91dd95cbd74d56d0678d36e5c79d512` | BSD license, as stated in its COPYING/README |

These are separate executable programs, not libraries linked into the ENets Go server. Their source and notices must remain available in any distributed source bundle. Re-check each pinned revision before release.

## Why this is an engine base, not a drop-in app fork

VirtualAP is not an Android `cmd wifi start-softap` wrapper. Its C engine creates/uses an AP interface (`ap0`), launches hostapd and dnsmasq, manages DHCP/DNS, and configures the routing/NAT needed to share an upstream connection. ENets will keep its existing Go server and web dashboard; the upstream Compose/Kotlin app and Droidspaces/OpenWrt managed mode are not part of the ENets UI.

The backend's minimum routing/NAT behavior is required for internet sharing. ENets will not expose firewall customization, client blocking/allow-lists, TTL spoofing, VPN routing, traffic monitoring, or OpenWrt management controls.

## Required ENets adaptations before enabling the backend

1. The CI build patches `VAP_BASE` from `/data/local/virtualap` to `/data/local/enets`, keeping runtime files isolated from the original VirtualAP app. Uninstall stops only the ENets instance and preserves its config/logs.
2. Make AP interface selection capability-driven. Upstream currently hardcodes `AP_IFACE=ap0`; this is not a cross-device guarantee.
3. The CI build removes the hardcoded `country_code=US` and `ieee80211d=1` lines so hostapd follows the device/kernel regulatory domain instead of forcing a country. Validate actual channel availability on each device.
4. Keep only the required start/stop/status/capability/client-read commands. Do not expose upstream block/allow commands in the ENets API or dashboard.
5. Keep 6 GHz disabled/unsupported until the engine and hostapd configuration are explicitly extended and validated. The inspected upstream CLI currently exposes only 2.4 and 5 GHz.
6. Detect the required interface/driver/hostapd capabilities before attempting a start. Do not infer AP support from `iw list` alone.
7. Make all configuration writes atomic and permission-restricted. Never log the hotspot password.
8. Keep the ENets dashboard bound to loopback. State-changing requests require same-origin checks; do not enable remote exposure without authentication and full CSRF protection.
9. Validate start, stop, status, reboot recovery, failure cleanup, and internet sharing on each device/firmware adapter. A successful process exit alone is not proof that clients can connect or reach the internet.

## Initial scope

- Magisk-only installation.
- ARM64 first; other ABIs remain unsupported until built and tested.
- SSID, passphrase, security, 2.4/5 GHz band, channel/width only when reported and validated.
- Start/stop and truthful runtime status.
- Read-only connected-client information only if it can be obtained reliably.
- No 6 GHz claim, no client blocking, no traffic accounting, no firewall editor.

## Status

The integration branch now has a CI workflow that checks out a pinned upstream commit, applies the ENets data-path and regulatory-domain patches, builds the ARM64 engine/daemons, and packages them into a Magisk preview. The module service stages the binaries under `/data/local/enets/bin`. The ENets Go server calls the engine directly, validates configuration, and requires a same-origin POST plus root execution for state changes. The dashboard now exposes same-origin, loopback-only start/stop controls and read-only status/client endpoints. The API writes a root-only configuration file and never returns the saved passphrase. A successful CI build is not device validation. The default `main` branch remains unchanged; integration work is isolated on `integration/virtualap-backend`.
