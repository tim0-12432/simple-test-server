Build / Lint / Test
- Build Go backend: `go build ./...` (works from repo root)
- Run backend tests: `go test ./...`
- Run a single Go test (package): `go test ./pkg/progress -run TestName`
- Run a single test by file: `go test ./controllers -run TestName` (package-level run)
- Frontend dev: `cd frontend && npm run dev`
- Frontend build: `cd frontend && npm run build`
- Frontend lint: `cd frontend && npm run lint`

Code Style Guidelines
- Formatting: use `gofmt`/`go fmt` for Go and `gofumpt` when available; run `go vet` before commits.
- Imports: keep standard library imports separate from third-party; use `goimports` to auto-group and remove unused imports.
- Types: prefer concrete types in package APIs; expose only necessary fields and methods (use unexported fields where appropriate).
- Naming: use mixedCase for local vars, PascalCase for exported identifiers, and short receiver names (e.g., `s *Server`).
- Errors: return wrapped errors using `fmt.Errorf("...: %w", err)` or `errors.Join`/`errors.Is` where appropriate; check and handle errors explicitly.
- Context: accept `context.Context` as first arg for long-running or I/O functions.
- Logging: keep logs in controllers or entrypoints; avoid global mutable loggers. Prefer structured logs.
- Tests: keep tests small and deterministic; use table-driven tests and `t.Parallel()` when safe.
- Concurrency: prefer channels and context cancellation; avoid data races — use `go test -race` when adding concurrency.

Cursor/Copilot Rules
- No `.cursor` or Copilot instruction files were found in this repository.

Notes
- Do not modify generated PocketBase or frontend build artifacts. Commit only source changes.
- When introducing new packages, add unit tests and update `go.mod` only via `go get` or `go mod tidy`.
