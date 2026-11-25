# ENets roadmap

## v0.1.0-alpha — read-only foundation
- [x] Public repository and project structure
- [x] English dashboard prototype
- [x] ARM64 Magisk package workflow
- [x] Loopback-only Go server and read-only capability endpoint
- [x] Safe install, service, and uninstall scripts
- [x] Cross-device adapter contract and feature-level fallback policy

## v0.2.0-alpha — VirtualAP integration preview
- [x] Pin and document the VirtualAP upstream source and license
- [x] Build ARM64 VirtualAP, hostapd, and dnsmasq in CI
- [x] Isolate the engine runtime under `/data/local/enets`
- [x] Remove the upstream hardcoded US regulatory-domain setting
- [x] Add same-origin, loopback-only start/stop API endpoints
- [x] Add SSID, password, security, 2.4/5 GHz, channel, width, and hidden-SSID controls
- [x] Add read-only status and connected-client API endpoints
- [x] Add validation tests for configuration and same-origin requests
- [ ] Confirm the latest CI build and corresponding-source artifact
- [ ] Install and test on the first ARM64 device
- [ ] Verify client association and actual internet sharing
- [ ] Verify upgrade, stop, reboot, and uninstall recovery paths

## Beta
- [ ] Introduce device/firmware adapter matching and capability-specific controls
- [ ] Validate supported bands, channels, channel widths, and security modes per device
- [ ] Add safe configuration backup and restore
- [ ] Test multiple OEM families and Android versions
- [ ] Consider additional ABIs only after separate build and device tests

## Stable
- [ ] Broader device/firmware compatibility matrix
- [ ] Complete security and source-compliance review
- [ ] Stable installation and recovery documentation
- [ ] Localization after Alpha/Beta
