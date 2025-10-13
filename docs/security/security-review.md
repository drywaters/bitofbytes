# BitOfBytes security review

## Optional hardening backlog

### Logging

- [x] Centralize logger construction so runtime log level and format can be
  controlled with `LOG_LEVEL` and `LOG_FORMAT`.
- [x] Record HTTP request metadata (method, path, status, duration, bytes
  written, and client addressing information) through middleware without logging
  payloads.
- [ ] Ship structured logs to a remote aggregator service with TLS transport.

The initial logging work adds a reusable factory in `models.NewLogger` which
builds either text or JSON handlers with the configured level. Request logging
is implemented as middleware so it executes for every request after security
wrappers (CSRF protection and secure headers) but before application handlers.
It purposely omits bodies and query strings to avoid leaking secrets while still
capturing the operational data required for incident response.
