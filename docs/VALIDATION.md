# Validation

The application is tested against temporary, isolated SQLite databases. Tests do not use or mutate the repository’s historical database or an operator’s running community.

## Automated coverage

`go test -race -cover ./...` exercises:

- Registration, login, posting, replies, logout, and profile navigation through HTTP.
- Anonymous write prevention, private liked/saved collections, idempotent bookmark writes, and removal.
- Author-only editing, topic replacement, preserved replies and reactions, revision increments, and rejected stale edits with retained form values.
- Command search result limits, Unicode query boundaries, and public-metadata-only responses.
- Saved-item ordering and repeat-save stability; positive profile targets and explicit-user precedence.
- Authentication quotas at limiter capacity and signed-in navigation during transaction errors.
- Topic/search filtering, sorting requests, empty states, and paginated feeds.
- Like/dislike switching and removal, comment reactions, and concurrent toggles.
- Session token hashing, restart persistence, expiry, and server-side invalidation.
- CSRF rejection, cross-origin rejection, secure-cookie flags, invalid-cookie recovery, and HTML escaping.
- Invalid and oversized topic selections, missing targets, invalid JSON reaction values, and transaction rollback.
- Idempotent schema initialization, non-duplicating sample setup, public assets, and the health endpoint.
- Exact ARIA reaction states on replies before and after a saved vote, including a page reload.

`go vet ./...` checks Go source. Node syntax checks cover `app.js`, `theme.js`, and `experience.mjs`. `node --test tests/*.test.mjs` covers system motion preferences, explicit visitor choices, bounded overscroll, short documents, independent layer movement, and editor word counts/reading estimates. Docker builds and the Compose health check validate the packaged runtime. CI runs the core checks on pushes and pull requests.

The dependency audit uses:

```bash
go run golang.org/x/vuln/cmd/govulncheck@latest ./...
```

The local audit during the redesign reported **zero reachable vulnerabilities**. It also reported one advisory in a required module outside the application’s imported packages; this is not a claim that every dependency is advisory-free. Audit results depend on the toolchain and advisory database at the time of the run.

## Visual and manual checks

The September 10, 2026 immersive design pass used the real Go application. Public views used `localhost:8088` with the opt-in demo dataset. Writing and saving checks used a disposable container, an isolated SQLite database, and a synthetic local account. The images in `docs/screenshots/` are unedited browser captures; the operator’s preview database was not used for test posts.

Completed browser checks:

- Dark opening, geometric feature cards, discussion feed, and search palette at a reported desktop viewport of 1294 × 912.
- Phone opening and feed at 390 × 844; narrow feed and focused editor at 320 × 760. Document width stayed within the viewport at these checked widths.
- Dark/light switching and Comfortable/Compact controls, including retained preferences after navigation.
- Ctrl K, live search for “technology,” arrow-key result selection, Enter to open a matching discussion, and Escape to dismiss a populated search and restore focus.
- Synthetic account sign-in, editor word counts, preview, focus mode, native publication, author editing, and the visible edited label. Submitting a preview with a missing title returned to Write and focused that field.
- Native bookmark submission, private Saved for later navigation, and card removal updating the count and empty state.
- Reopening a cancelled command search returned fresh results; modal expansion state followed opening and closing.
- Editor change tracking switched to unsaved when a title changed and returned to clean when the original text was restored. Native leave-page dialogs and print output remain browser-dependent release checks.
- Local fonts and artwork loading without third-party asset requests. The browser reported one skipped view transition during rapid navigation; the application remained usable. No application-script error was observed.
- System reduced motion enabled by default, explicit motion opt-in, distinct scroll transforms, and a persistent pause control.

The final local Docker test stage passed Go race tests and static analysis. All seven JavaScript unit tests passed, including editor change tracking. The packaged application built and reported healthy through Compose. HTTP integration tests cover security, ownership, privacy, concurrency, and validation independently of the browser walkthrough.

Broader release checks still include 200% zoom, browser/screen-reader combinations, reaction network failures, and navigation with JavaScript disabled. Recorded checks do not establish accessibility conformance or full production readiness. See [application scope](ARCHITECTURE.md#scope) for features that remain outside this project.
