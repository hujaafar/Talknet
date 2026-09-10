# Validation

The application is tested against temporary, isolated SQLite databases. Tests do not use or mutate the repository’s historical database or an operator’s running community.

## Automated coverage

`go test -race -cover ./...` exercises:

- Registration, login, posting, replies, logout, and profile navigation through HTTP.
- Anonymous write prevention and private liked-discussion visibility.
- Topic/search filtering, sorting requests, empty states, and paginated feeds.
- Like/dislike switching and removal, comment reactions, and concurrent toggles.
- Session token hashing, restart persistence, expiry, and server-side invalidation.
- CSRF rejection, cross-origin rejection, secure-cookie flags, invalid-cookie recovery, and HTML escaping.
- Invalid and oversized topic selections, missing targets, invalid JSON reaction values, and transaction rollback.
- Idempotent schema initialization, non-duplicating sample setup, public assets, and the health endpoint.

`go vet ./...` checks Go source. `node --check static/js/app.js` checks JavaScript syntax. Docker builds and the Compose health check validate the packaged runtime. CI runs the core checks on pushes and pull requests.

The dependency audit uses:

```bash
go run golang.org/x/vuln/cmd/govulncheck@latest ./...
```

The local audit during the redesign reported **zero reachable vulnerabilities**. It also reported one advisory in a required module outside the application’s imported packages; this is not a claim that every dependency is advisory-free. Audit results depend on the toolchain and advisory database at the time of the run.

## Visual and manual checks

The local preview is served from the real Go application. The README artwork is an original asset, not a screenshot or evidence of browser testing.

For manual release checks, review the feed, detail, profile, sign-in, registration, and composer at desktop and narrow mobile widths. Exercise keyboard navigation, 200% zoom, system reduced motion, the pause control, reaction errors, and navigation with JavaScript disabled. Automated HTTP tests do not establish pixel-perfect rendering or browser accessibility conformance.
