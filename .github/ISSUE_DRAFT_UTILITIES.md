# Utility Ideas for /utils

This issue tracks proposed additions to the Utilities section. Each item can be split into its own issue/PR.

## Quick Wins (stdlib-only)
- URL Encoder/Decoder (`net/url`): `/utils/url/{encode,decode}`
- Hash Generator (MD5/SHA-1/SHA-256): `/utils/hash`
- HMAC Generator (HMAC-SHA256): `/utils/hmac`
- JWT Decoder (optional verify): `/utils/jwt`
- Timestamp Converter (Unix ↔ RFC3339): `/utils/time`
- JSON Tools (pretty/minify/validate): `/utils/json`
- CSV → JSON (paste/upload): `/utils/csv`
- Slugify: `/utils/slugify`
- HTML Escape/Unescape: `/utils/html`
- Random Password/Token (CSPRNG): `/utils/random`

## Networking/Systems
- DNS Lookup (A/AAAA/CNAME/MX): `/utils/dns`
- CIDR Calculator: `/utils/cidr`

## Optional (adds deps)
- UUID v4 Generator: `/utils/uuid` (dep: github.com/google/uuid)
- QR Code Generator: `/utils/qrcode` (dep: github.com/skip2/go-qrcode)

## Notes
- Pattern: add templates under `templates/utils/<tool>/`, handlers in `controllers/`, routes in `cmd/bob/bob.go`, render via `views.ParseFS(templates.FS, ..., "base.gohtml").`
- Add links on `/utils` and include CSRF field for forms.
