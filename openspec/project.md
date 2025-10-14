# Project Context

## Purpose
Simple Test Server is a small, developer-focused application for spinning up short-lived, disposable protocol servers (MQTT, SMTP, FTP, SMB, HTTP/web, etc.) inside Docker containers. It exists to help QA engineers and developers quickly run and inspect lightweight servers and services for integration and manual testing without needing production infrastructure. The app provides an API + web UI to configure, start/stop, observe, and interact with those test servers and their logs/messages.

## Tech Stack
- Backend: Go 1.25, Gin (HTTP framework), Viper (configuration)
- Database: PocketBase (embedded NoSQL used by the Go backend, exposed at `/pb/`)
- Container/runtime: Docker (builders and runners, custom images under `custom_images/`)
- Frontend: React (v19) + TypeScript (v5), Vite (build), Bun (optional dev/runtime), TailwindCSS (styling), Shadcn UI components
- Tooling: go toolchain (`go test`, `go build`), Vitest + React Testing Library for frontend tests, Docker Compose for local orchestration

## Project Conventions

### Code Style
- Go:
  - Follow Go idioms: `gofmt`/`go fmt` formatting, short package names, clear exported identifiers and package comments.
  - Context-first function signatures (`context.Context` as first parameter).
  - Return `(T, error)` pairs; wrap errors with `%w` and use `errors.Is`/`errors.As` for checks.
  - Keep functions small and focused; avoid panics in libraries.
  - Place `go test ./...` as the standard test command.
- Frontend (TypeScript + React):
  - Single quotes, semicolons enabled, named exports preferred.
  - Files and components use PascalCase; hooks/utilities use camelCase.
  - Prefer controlled components for forms; type all props and avoid `any`.
  - Use TailwindCSS classes (or the chosen styling approach) consistently per component.
  - Run `bun install` (or `npm install`) then `bun run dev` (or `npm run dev`) for local frontend dev.
- Linting & Formatting:
  - Keep lints in CI (ESLint for frontend, `go vet`/`golangci-lint` optionally for Go).

### Architecture Patterns
- Separation of concerns: `controllers/` for HTTP handlers, `docker/` for container build/run logic, `db/` for PocketBase DTOs and services, `protocols/` for per-protocol controllers/services.
- Backend serves the frontend build from `frontend/dist/` and also exposes a REST API for container management and logs.
- PocketBase stores collections (containers, logs, servers) and is bootstrapped from `db/collection_bootstrap.go` when needed.
- Docker image build is handled by `docker/builder.go` with custom image sources in `custom_images/`.
- Progress/events: `progress/hub.go` provides an event hub for streaming progress and logs to clients.

### Testing Strategy
- Backend:
  - Unit tests with `go test ./...` covering packages; tests live next to code as *_test.go files.
  - Use table-driven tests and `t.Parallel()` where safe.
- Frontend:
  - Component and hook tests via Vitest + React Testing Library. Run with `bun run test` or `npm run test`.
  - Prefer testing behavior over implementation; use role/name queries.
- Integration / Manual:
  - Docker-based manual/integration testing by starting services via the UI or `docker-compose.yml`.
  - Keep integration tests lightweight; prefer manual verification for heavy e2e scenarios.

### Git Workflow
- Branching: feature branches off `main` (e.g., `feature/xxx`), open Pull Requests for review, merge with meaningful commit messages.
- Commits: Prefer conventional, descriptive commits (e.g., `feat(mqtt): add topic inspector` or `fix(docker): handle image build error`).
- PRs should include a short summary of the why and any migration or data effects.
- NOTE: automation/agents should not perform git pushes or commits without explicit human approval.

## Domain Context
- The system models "test servers" as Docker containers with a protocol type (mqtt, web, mail, ftp, smb, etc.).
- Key domain entities:
  - Container / Server: configuration and runtime metadata for a running test server.
  - Logs / Messages: protocol-specific data captured from running containers (MQTT messages, SMTP messages, HTTP access logs).
  - Images: prebuilt or custom images used to start test servers; custom variants live in `custom_images/`.
- UI organizes servers as tabs (one per running server instance) and provides tools for inspecting live messages and logs.
- PocketBase collections used by backend: containers, logs, servers — see `db/dtos/` and `db/services/` for concrete fields.

## Important Constraints
- Keep the application self-contained for local development: PocketBase is embedded and Docker must be available on the host.
- Avoid adding heavy external dependencies without justification; prefer stdlib and small, well-maintained libs.
- Containers may require elevated privileges on some hosts (Docker Desktop/WSL/Podman differences) — document host platform caveats.
- Data persistence: `pb_data/` (or configured path) stores PocketBase data; treat it as the canonical local DB for dev.
- Security: the app is intended for local/dev environments. Do not expose to untrusted networks without proper hardening.

## External Dependencies
- Docker Engine / Docker Desktop (required to build and run test server containers)
- Optional: Docker Compose for multi-service scenarios (provided `docker-compose.yml`)
- PocketBase (embedded; managed by the Go backend)
- Bun (optional dev runtime for frontend) or Node/npm as alternative
- CI runners: typical Linux-based runners for tests and builds (ensure Docker availability where integration tests run)


## Useful Commands (Developer Quick Reference)
- Backend run: `go run .` or `go run main.go`
- Backend tests: `go test ./...`
- Frontend dev: `cd frontend && bun install && bun run dev` (or `npm install && npm run dev`)
- Frontend tests: `cd frontend && bun run test` (or `npm run test`)
- Build frontend for shipping: `cd frontend && bun run build` (or `npm run build`)
- Start locally with compose: `docker compose up --build`


---

This file should be kept up-to-date with any architectural or operational changes. For larger proposals or breaking changes, follow the project's `openspec/AGENTS.md` guidance for formal change proposals and reviews.
