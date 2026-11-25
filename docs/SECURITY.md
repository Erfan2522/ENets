# Security model

ENets v0.3.0-alpha is an experimental, local-only hotspot manager.

- The Go server binds only to `127.0.0.1:8080` and rejects non-loopback peers and untrusted Host headers.
- State-changing requests require same-origin validation and a root-running service.
- Hotspot settings and profiles are stored under `/data/local/enets` with mode 0600; the parent directory is mode 0700.
- Writes use temporary files, fsync and atomic rename.
- The hotspot passphrase is persisted because the user requested settings restoration. It is returned only to the loopback dashboard so the form can be restored; it is never included in diagnostic responses.
- Diagnostic log export is limited to a recent excerpt and redacts common password/passphrase fields. Treat logs as sensitive nonetheless.
- Auto Start is opt-in and defaults to OFF. Boot performs one start attempt only when the saved preference is enabled.
- WPA3 SoftAP is not enabled based on generic radio detection. A device/firmware-specific SAE capability check must be validated first.
- Updates are informational in this alpha; users review releases and install manually.
- Never expose the server to the LAN or hotspot clients.
