# Repository Guidelines

## Project Structure & Module Organization
- `cmd/bob/`: application entrypoint (`main` package).
- `controllers/`: HTTP handlers grouped by feature (e.g., `blog`, `utils`).
- `models/`: services, data types, and config loading.
- `views/`: template execution helpers and shared funcs.
- `templates/`: HTML templates (`.gohtml`) embedded via `templates.FS`.
- `posts/`: markdown posts (embedded via `posts.FS`), named by slug (e.g., `my-post.md`).
- `static/`: public assets (images, compiled CSS, JS).
- `tailwind/`: Tailwind source and config.

## Build, Test, and Development Commands
- `make local`: runs Tailwind watcher and Air for live reload.
- `make run`: starts the app via Air (`.air.toml` builds `./cmd/bob`).
- `make tail-watch` / `make tail-prod`: generate CSS (watch/minified).
- `make docker-build` / `make docker-push`: build and push `drywaters/bob` image.
- Manual build: `go build -o ./tmp/bob ./cmd/bob`.

Prereqs: Go 1.22+, `air`, `tailwindcss` on PATH. Copy `.env.template` to `.env`.

## Coding Style & Naming Conventions
- Go formatting: run `gofmt` (tabs, 1TBS braces). Prefer `go vet` before PRs.
- Packages: lowercase, short; files use underscores (e.g., `base64.go`).
- Exported identifiers: `CamelCase`; unexported: `camelCase`/`lowercase`.
- Templates: `.gohtml`; group by category (e.g., `utils/base64/encode.gohtml`).
- Routes: register in `cmd/bob/bob.go`; keep handlers in `controllers/`.

## Testing Guidelines
- Framework: Go standard testing. Place tests as `*_test.go` beside code.
- Run tests: `go test ./...` (add `-cover` for coverage).
- Table-driven tests for services in `models/`; minimal handler tests for controllers.

## Commit & Pull Request Guidelines
- Commits: short, imperative, scoped (e.g., “Fix embed posts”, “Remove chi”).
- PRs: clear description, linked issues, and screenshots for UI changes.
- Include: what/why, testing notes, relevant commands (e.g., `make tail-prod`).

## Security & Configuration
- `.env`: set `SERVER_ADDRESS`, `CSRF_KEY`, `CSRF_SECURE`.
- Use a 32-byte random `CSRF_KEY`; never commit secrets.
- Static files served from `/static/`; templates and posts are embedded at build time.

## Adding Content
- New post: create `posts/<slug>.md` with TOML front matter (`title`, `slug`, `author`).
- New page: add template under `templates/<section>/`, wire route in `bob.go`, render via `views.ParseFS`.
