# Architecture

Talknet serves complete HTML pages from Go. SQLite is the source of truth; JavaScript enhances the rendered interface without duplicating the application state in a client framework.

```mermaid
flowchart LR
    Browser[Browser] --> HTTP[Go HTTP server]
    HTTP --> Security[Headers, body limits, CSRF and origin checks]
    Security --> Routes[Forum handlers]
    Routes --> Views[Embedded Go templates]
    Routes --> DB[(SQLite + WAL)]
    HTTP --> Assets[Embedded public assets]
    Views --> Browser
```

## Request flow

`forum.New` builds an isolated router and template set. Security middleware applies response headers and request limits, issues a random CSRF cookie, and checks state-changing requests against a hidden form token or `X-CSRF-Token`. Go’s cross-origin protection also checks browser request metadata.

The session cookie is a cryptographically random bearer token. Only its SHA-256 hash is stored in `Sessions`. Queries check expiry on every authenticated request; logging in replaces that user’s old session. Logout removes the server record. Sessions survive process restarts and expire after 24 hours.

Authentication attempts are bounded by remote IP, with a one-minute window and ten attempts. The limiter has bounded memory and is local to each process. It intentionally does not trust arbitrary forwarded IP headers.

## Data and consistency

The original users, posts, comments, topics, post/topic links, reactions, and session tables remain compatible. Schema initialization is idempotent. Two additive tables support the new experience: `Bookmarks` has a unique `(user_id, post_id)` key, and `Post_Revisions` records a revision number while older posts default to revision one. Foreign keys are enabled on the connection; WAL and a five-second busy timeout support normal local concurrency.

A single database connection serializes writes. Transactions commit a post and its topic links together, and serialize reaction toggles so a member has at most one reaction per target through the application. Deferred rollbacks cover every failure path. Query rows are closed before additional queries run to avoid blocking the only connection.

Editing rechecks the author and current revision inside the write transaction. A stale submission returns HTTP 409 and retains the submitted text and topics in the form; it never silently overwrites the newer post. Successful edits preserve the post ID, publication date, replies, votes, and saves. The revision table records a number, not a recoverable edit history.

Bookmark writes use explicit `save` or `remove` actions, making retries idempotent. Both native forms and JSON requests use the same authenticated, CSRF-protected route. Saved and liked profile tabs are restricted to their owner.

`GET /api/search?q=` returns at most six matching discussions with only ID, title, author name, and first topic. Search matches titles and content. The palette debounces requests, aborts superseded requests, and ignores late responses; rendered results use `textContent`.

The discussion feed uses a fixed page size of 20, fetching one additional record to discover a next page. Search and topic filters carry through pagination. Detail requests load the requested post directly. Discussion replies are currently rendered as one chronological list.

## Rendering and assets

Templates use Go’s contextual escaping. Templates render to a buffer before writing a response, avoiding a partial HTML success page when rendering fails. The server exposes only the public `styles`, `js`, and `images` paths; template source is not a public asset route.

The binary embeds all assets. Inter and Instrument Serif are served locally, and the original artwork is delivered as a 1536 × 1024 WebP. The browser makes no font, script, or image requests to third-party services. Documentation screenshots are outside the embedded asset tree.

## Motion and progressive enhancement

`app.js` uses a passive scroll listener and one animation frame at a time. Its `motion.mjs` import calculates bounded transforms for the artwork, title, and conversation note, plus reading progress. Intersection Observer reveals individual discussion rows. Essential content is never hidden while waiting for JavaScript. The pure motion functions have dependency-free Node tests covering overscroll, short documents, different layer depths, and user preferences.

The fixed motion control remembers a device-local preference. The system reduced-motion setting applies by default; an explicit visitor choice can enable or pause motion. Small screens keep the artwork static, and the topic navigation becomes horizontally scrollable. JavaScript adds reaction requests, character counters, password visibility, copy-link feedback, and duplicate-submit protection.

`theme.js` reads a validated dark/light preference before styles render. `experience.mjs` adds the command palette, private save shortcuts, feed density, plain-text editor preview, word counts, focus mode, and bounded pointer depth. Preferences use local storage with a fallback when storage is unavailable. Native dialog semantics provide modal focus containment; keyboard handlers support Ctrl/Cmd K, arrow navigation, and Escape.

The native submit guard prevents repeat submissions without disabling the submitter before its name/value is serialized. Returning through browser history resets the guard. Editor preview reveals the writable fields if native form validation finds an invalid field.

## Scope

This is a discussion application for a single Go process with SQLite storage. Moderation tooling, post deletion, email verification, password recovery, uploads, and multi-instance operation are not implemented. The `/profile?user=` parameter identifies a member; legacy `/profile?id=` links still resolve the author of that post. Liked and saved discussions are visible only to their owner. Draft autosave, Markdown rendering, and edit-history recovery are not implemented.
