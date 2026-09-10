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
- Exact ARIA reaction states on replies before and after a saved vote, including a page reload.

`go vet ./...` checks Go source. `node --check static/js/app.js` checks JavaScript syntax. `node --test tests/*.test.mjs` covers system motion preferences, explicit visitor choices, bounded overscroll, short documents, and independent layer movement. Docker builds and the Compose health check validate the packaged runtime. CI runs the core checks on pushes and pull requests.

The dependency audit uses:

```bash
go run golang.org/x/vuln/cmd/govulncheck@latest ./...
```

The local audit during the redesign reported **zero reachable vulnerabilities**. It also reported one advisory in a required module outside the application’s imported packages; this is not a claim that every dependency is advisory-free. Audit results depend on the toolchain and advisory database at the time of the run.

## Visual and manual checks

The September 10, 2026 design pass used the real Go application at `localhost:8088` with the opt-in demo dataset. The images in `docs/screenshots/` are unedited browser captures of that application.

Completed browser checks:

- Opening and discussion feed at a reported desktop viewport of 1294 × 912 and a phone viewport of 390 × 844.
- Feed at 320 × 760, including correction and recheck of a heading/action overflow. Document width stayed within the viewport at the checked widths.
- Search for “technology,” opening its matching discussion, navigating to the author profile, and inspecting registration and sign-in layouts.
- Loaded local fonts and artwork; no browser warning/error log entries during the check.
- System reduced motion enabled by default, explicit motion opt-in, separate scroll transforms for artwork/title/note, a working pause control, and preference retention through navigation.
- Visible account-page primary headings on mobile, singular reply labels, and exact reply reaction state tokens.

The final local Docker test stage passed the Go race tests and static analysis. All four JavaScript motion tests passed. The packaged application started successfully and reported healthy through Compose.

For broader release checks, also exercise the signed-in composer visually, keyboard navigation, 200% zoom, browser/screen-reader combinations, reaction network failures, and navigation with JavaScript disabled. Account creation, posting, replies, and reactions are covered by HTTP integration tests; the recorded browser pass did not submit an account or publish content. These checks do not establish accessibility conformance.
