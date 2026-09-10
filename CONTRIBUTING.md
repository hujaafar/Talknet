# Contributing to Talknet

Start with the [README](README.md) for local setup and the [architecture](docs/ARCHITECTURE.md) for request flow, data consistency, and current scope.

## Work locally

Create a branch for one coherent change. Use an empty local database or the optional demo seed; do not develop against a live community’s data. The Go application embeds its interface, so rebuild the binary or run `docker compose up -d --build --wait` after template, style, or script edits.

Keep database writes parameterized and transactional where multiple records must change together. Preserve server-rendered navigation and native form submissions. For interface changes, follow the [design system](docs/DESIGN.md), check a narrow phone layout, and exercise the motion toggle.

## Validate a change

With Go 1.26+, a C compiler, and Node.js 22+:

```bash
gofmt -w main.go internal/forum/*.go static/*.go
go test -race -cover ./...
go vet ./...
go build -trimpath ./...
node --check static/js/app.js
node --test tests/*.test.mjs
```

The Go test suite uses temporary databases. The JavaScript tests use Node’s built-in runner and require no package installation. `docker build --target test -t talknet-tests .` provides the Go checks in a container; CI also starts the packaged application and checks its health endpoint.

Add a regression test when fixing behavior that can recur. For a visual change, include a capture from the running application and describe the viewport and interactions checked. Avoid committing databases, credentials, generated binaries, or local environment files.

## Open a pull request

Explain the concrete problem, what the change does, and how it was checked. Include any migration or operating impact. Keep unrelated formatting or refactoring separate so reviewers can assess the behavior. Retain attribution and font licenses when changing assets.

For bug reports, include the route, steps to reproduce, expected behavior, and relevant browser or runtime version. Use fictional data and remove session cookies, passwords, database contents, and personal details from logs or screenshots.
