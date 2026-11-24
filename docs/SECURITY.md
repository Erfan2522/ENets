# Security model

This alpha is a UI and repository scaffold, not a deployable hotspot controller.

Before remote dashboard access is enabled:
- Require an administrator-created credential; never ship a shared default password.
- Store password verifiers using a suitable password-hashing algorithm and random salt.
- Use secure session handling, CSRF protection, login throttling, and logout/revocation.
- Bind only to the detected hotspot interface and reject requests from other interfaces.
- Do not expose the service to the general LAN or internet.
- Never log hotspot passphrases, admin passwords, session tokens, or other secrets.
- Validate every setting on both client and server.
- Keep configuration files readable only by the module/service account.
- Do not change firewall rules or block clients.
- Require explicit user approval for updates and configuration restore.

The current Go server binds to loopback only. Keep that default until the controls above are implemented and tested.
