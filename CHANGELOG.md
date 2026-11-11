# Changelog

All notable changes to this template. Format based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## 2026-11-11

### Added
- `./artisan deploy [--name NAME] [--skip-frontend]` (Linux, systemd): builds the executable, then installs and starts a systemd **user** service on the first run and restarts it afterwards. An existing unit file is never read or rewritten, so local edits (port, environment) survive every deploy. Ported from a downstream project.
- `APP_NAME` in `.env` names the executable (`build/<name>`), the systemd unit, `artisan about`, and the development binaries, so each project is identifiable in htop/ps on a shared server. It is reduced to a lowercase slug and defaults to `inertia-go-solid-template`.
- Default error handling. Go: unhandled errors and panics render a detailed page (message, request, stack trace) or JSON in every environment except `production`, where the response is a generic "500 Internal Server Error" without details. Solid: `AppErrorBoundary`, `ErrorPanel` and `AsyncErrorBanner` (`resources/js/components/error_boundary.tsx`) show a detailed panel only in `./artisan dev` and a generic message in production builds. Both are defaults; pass your own `fallback`, or replace the Fiber `ErrorHandler`, wherever you need a different look.
- `resources/js/components/link.tsx`: Inertia `Link` that preloads the target page on hover/focus by default. New application code should import it instead of `Link` from `inertia-adapter-solid`.

- README section "Live demo websites" with a metrics table (Lighthouse, bundle size, load and SPA navigation times) for sites built with the template.

### Changed
- `scripts/start.mjs`, `scripts/smoke.mjs` and `scripts/dev.mjs` derive the binary name from `APP_NAME` through `scripts/process.mjs`.

## 2026-10-07

### Changed
- Project skill documents local development (ports, LAN access, Windows launchers) and cross-platform checks.

## 2026-10-05

### Added
- Browser-free production smoke test (`npm run test:portable`) that runs the compiled binary from an empty directory on Linux, macOS and Windows; CI runs the native jobs on all three.
- Static asset tests with deliberately wrong system MIME registrations and in-memory asset filesystems.
- `.gitattributes` normalizing launcher and license line endings.

### Changed
- Static asset serving is portable: explicit MIME types for JS/MJS, CSS, fonts and WASM, immutable caching, and rejection of hidden or traversal paths.
- Font licenses use consistent line endings.

### Initial release
- Fiber + inertia-go backend and a common SolidJS frontend in one Go module.
- `./artisan` CLI (also `artisan.cmd` / `artisan.ps1`): `install`, `dev`, `build`, `start`, `serve`, `test`, `route:list`, `key:generate`, `make:controller`, `make:component`, `make:page`, `about`.
- Single-executable production build: the Vite bundle is embedded with `//go:embed` behind `-tags production`.
- Development with Vite HMR and a bundled Node watcher that rebuilds and restarts the Go server; no Air dependency.
- Route registry in `routes/web.go`, one-package-per-controller layout with one snake_case file per action.
- Encrypted session cookies keyed by `APP_KEY`.
- Local helpers for pinned upstream gaps (lazy props, full-URL prefetch cache, Promise-to-signal JSON helper), each covered by known-limit tests.
- Demo pages for the 24 Inertia feature families, Go protocol tests, Playwright scenarios, and a verification report (`reports/template-verification.md`).
- Project skill (`.agents/skills/inertia-go-solid`) and `AGENTS.md` conventions.
