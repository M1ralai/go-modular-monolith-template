# Go Modular Monolith Reference

A reference template for organizing a Go HTTP service as a modular monolith. It demonstrates application composition, module boundaries, PostgreSQL repositories, authentication middleware, migrations, background jobs, WebSockets, logging, and metrics through example planning-domain modules.

This repository is a template and architecture study, not a finished product. The included domain modules are examples of how the structure can be applied; their presence should not be read as a complete application feature set.

## What it demonstrates

- A single deployable application composed from independently structured domain modules
- Explicit dependency construction in one application composition root
- HTTP, service, repository, and domain boundaries inside each module
- PostgreSQL access with `sqlx` and versioned migrations
- Shared JWT, recovery, timeout, logging, validation, and metrics infrastructure
- Background-job and WebSocket infrastructure that modules can opt into

## Architecture and dependency direction

```text
cmd/api
  -> internal/app
       -> modules/<name>/http
            -> modules/<name>/service
                 -> modules/<name>/domain
                 -> repository interface
                      <- modules/<name>/repository/postgres

shared infrastructure: database, middleware, jobs, logging, metrics, WebSockets
```

`cmd/api` owns process startup and shutdown. `internal/app` is the composition root: it creates infrastructure, repositories, services, and handlers, then registers routes. Transport code depends on services; services depend on domain contracts; PostgreSQL implementations satisfy repository interfaces at the outside edge.

Example modules under `internal/modules` cover areas such as authentication, tasks, notes, schedules, habits, events, and users. They exist to exercise the boundaries and registration pattern.

## Module layout

```text
internal/modules/<name>/
  domain/       entities and repository contracts
  repository/   PostgreSQL implementations
  service/      use-case logic
  http/         request handlers and route registration
```

To add a module, define its domain contracts, implement persistence, add service behavior, expose the HTTP adapter, add migrations when needed, and wire the module in `internal/app/server.go`.

## Run locally

Requirements: Go 1.25+ and PostgreSQL.

```bash
git clone https://github.com/M1ralai/go-modular-monolith-template.git
cd go-modular-monolith-template
cp .env.example .env
go mod download
go run cmd/api/main.go
```

Configure the database and JWT values in `.env` before running. Migrations are embedded and applied during startup.

Example infrastructure endpoints include:

```bash
curl http://localhost:8080/health
curl http://localhost:8080/metrics
```

## Limitations

- This is an opinionated reference implementation, not a production-ready starter kit.
- The large set of example modules makes the repository less minimal than a conventional template.
- Automated tests do not yet cover the architecture and example modules comprehensively.
- Deployment hardening, secret management, backups, TLS, and operational policies remain environment-specific.
- Generated and local-development artifacts should be cleaned before using the repository as a new project base.

## License

MIT
