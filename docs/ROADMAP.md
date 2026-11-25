# ENets roadmap

## v0.1.0-alpha — read-only foundation
- [x] Public repository and project structure
- [x] English dashboard prototype
- [x] ARM64 Magisk package workflow
- [x] Loopback-only Go server and read-only capability endpoint
- [x] Safe install, service, and uninstall scripts
- [x] Cross-device adapter contract and feature-level fallback policy

## v0.2.0-alpha — VirtualAP integration preview
- [x] Pinned VirtualAP engine and ARM64 payload
- [x] ENets-specific runtime path and regulatory-domain adaptations
- [x] Loopback-only hotspot start/stop controls
- [x] SSID, password, security, band, channel, width and hidden SSID
- [x] Read-only hotspot status and connected-client list
- [x] User-reported POCO M5s Android 17 test

## v0.3.0-alpha — persistence and reliability
- [x] Root-only atomic settings persistence with private file permissions
- [x] Save settings as the user changes them
- [x] Opt-in Auto Start after reboot; default OFF
- [x] Named hotspot profiles
- [x] Read-only diagnostics endpoint
- [x] Version and release information endpoint
- [x] Keep dashboard bound to loopback only
- [ ] Run CI tests and inspect ARM64 package
- [ ] On-device test settings persistence across app close and reboot
- [ ] Test Auto Start OFF and ON paths
- [ ] Verify upgrade preserves settings and profiles
- [ ] Validate WPA3 SoftAP capability on supported devices before enabling it
- [ ] Validate connected-client metadata and recovery behavior

## Beta
- [ ] Device/firmware adapter matching and capability-specific controls
- [ ] Validated WPA3/WPA2 transition mode per device
- [ ] Safe configuration backup and restore
- [ ] Multiple OEM families and Android versions

## Stable
- [ ] Broader device/firmware compatibility matrix
- [ ] Complete security and source-compliance review
- [ ] Stable installation and recovery documentation
- [ ] Localization after Alpha/Beta
