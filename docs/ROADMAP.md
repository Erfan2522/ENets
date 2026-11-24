# ENets roadmap

## v0.1.0-alpha — installable read-only preview
- [x] Public repository and project structure
- [x] English dashboard prototype
- [x] ARM64 Magisk preview package workflow
- [x] Loopback-only Go server and read-only capability endpoint
- [x] Safe install, service, and uninstall scripts
- [x] Cross-device adapter contract and graceful feature-level fallback policy
- [ ] Verify the packaged ZIP installs and starts on the first test device
- [ ] Validate device capability reporting on the first test device

## v0.2.0-alpha — Android integration foundation
- [ ] Authenticated local control API before any write operation
- [ ] Adapter registry and device/firmware matching
- [ ] First tested Android hotspot read-only adapter
- [ ] Read-only connected-client discovery
- [ ] Safe configuration backup and restore design

## Beta
- Apply SSID and password on verified device/firmware combinations
- Detect and expose only supported security modes and Wi-Fi bands
- Add channel and client-limit controls only where validated
- Test install, upgrade, disable, and uninstall paths across multiple OEM families

## Stable
- Broader device/firmware compatibility matrix
- Complete security review and installation documentation
- Localization work after Alpha/Beta
