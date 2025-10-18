# BitOfBytes Project Notes


## Overview
- This repository contains a Go 1.22 web application for the BitOfBytes site.
- The HTTP entrypoint is `cmd/bob/bob.go`; it wires up controllers, parses templates from `templates.FS`, embeds blog posts from `posts.FS`, and serves static assets from `/static`.
- Controllers live in `controllers/`, services and data types live in `models/`, reusable view helpers live in `views/`, and Go HTML templates are under `templates/`.

## Configuration
- Runtime configuration is loaded from environment variables via `models.LoadEnvConfig()`.
- Copy `.env.template` to `.env` (the app uses `github.com/joho/godotenv` to load it automatically) and fill in:
  - `SERVER_ADDRESS` (default `:3000` for local development).
  - `CSRF_KEY` (32-byte base64 string recommended) and `CSRF_SECURE` (set `false` for HTTP during local work).

## Running the server
- Run `go run ./cmd/bob` after preparing the `.env` file to start the web server locally.
- For live-reload development you can install `air` and use `make run`, which runs `air` alongside the Tailwind watcher (`make tail-watch`).
- Tailwind CSS assets are generated from `tailwind/styles.css` into `static/styles.css`. Use `make tail-prod` to build a minified stylesheet for production.

## Content & rendering
- Blog posts live in `posts/*.md` with TOML frontmatter parsed by `github.com/adrg/frontmatter`. Markdown content is rendered to HTML with `github.com/yuin/goldmark` plus the syntax-highlighting extension.
- Utility pages for Base64 encoding/decoding are backed by `models.Base64Service` and templates in `templates/utils/`.
- Gorilla CSRF middleware wraps the router; if you add new POST routes make sure to include CSRF tokens in the forms (`csrf.TemplateField`).

## Tests & tooling
- There are currently no automated tests. Add Go unit tests under the relevant package directories and run them with `go test ./...`.
- Docker builds use the recipes in the `Makefile` (see `make docker-build` / `make docker-push`) and the multi-stage image defined in `Docker/Dockerfile`.

## Miscellaneous tips
- Static assets (images, PDFs, JS) are stored under `static/` and served at `/static/` via `http.FileServer`.
- Templates use the helper types in `views/` (e.g., `views.Page`, `views.ParseFS`); new templates should follow the existing pattern and be parsed via the views package.
- If you add new directories, remember that `templates/fs.go` and `posts/fs.go` use `go:embed` directives, so keep file names glob-friendly.
