# HTTP interface

Talknet serves complete HTML pages. Two small JSON interfaces support command search and reactions; the bookmark route supports both HTML forms and JSON responses. This is a same-origin, cookie-authenticated application rather than a token-based public API.

## Reading and navigation

| Method and route               | Purpose                                                                     |
| ------------------------------ | --------------------------------------------------------------------------- |
| `GET /`                        | Discussion feed. Optional `q`, `category`, `sort=latest                     | popular`, and `page` parameters. |
| `GET /post-details?post_id=ID` | A discussion and its replies.                                               |
| `GET /profile?user=ID`         | Public authored discussions.                                                |
| `GET /profile?tab=liked`       | The signed-in member’s liked discussions.                                   |
| `GET /profile?tab=saved`       | The signed-in member’s private saved collection, most recently saved first. |
| `GET /api/search?q=TEXT`       | Up to six matching discussions as JSON.                                     |
| `GET /healthz`                 | Database connectivity check; plain-text `ok` on success.                    |

Feeds render 20 discussions per page. Search examines titles and content using SQLite `lower` and `instr`; percent signs and underscores are literal search characters. Queries may contain up to 200 Unicode code points. A missing, negative, invalid, or excessively large page number falls back to page one.

Explicit `user` profile targets must be positive integers and take precedence over the legacy `id` parameter, which identifies the author of a post. Private tabs requested for another member fall back to that member’s public discussions. Counts beside profile tabs describe the current page.

### Command search

```bash
curl --get --data-urlencode "q=technology" http://localhost:8088/api/search
```

Example response shape; IDs and content depend on the database:

```json
{
  "results": [
    {
      "id": 5,
      "title": "A question about technology",
      "author": "DemoAuthor",
      "topic": "Technology"
    }
  ]
}
```

Fewer than two code points returns an empty results array. Overlong queries return HTTP 400. Results contain only public ID, title, author name, and first topic; there are no email addresses, session tokens, post bodies, or private save states in this response. The database query itself is capped at six rows.

## Authentication and CSRF

`GET /register` and `GET /login` render native forms. Registration posts `username`, `email`, and `password` to `/register`; login posts `username` and `password` to `/login`. Successful registration redirects to sign-in. Login replaces the member’s previous session and redirects to the feed. `POST /logout` invalidates the session and clears its cookie.

All writes require a CSRF token from the rendered page, supplied as `csrf_token` in a form or as `X-CSRF-Token`. Keep the corresponding cookies with the request. JavaScript reads the page’s `meta[name="csrf-token"]`; it does not read the HTTP-only session cookie. Cross-origin write protection applies in addition to token verification.

Personalized pages and JSON responses use `Cache-Control: no-store`. Session identity is resolved once per request before write transactions, so error pages can retain the correct navigation without querying an occupied SQLite connection.

## Writing and editing

| Method and route            | Form fields                                                                |
| --------------------------- | -------------------------------------------------------------------------- |
| `GET /post`                 | Render the signed-in composer.                                             |
| `POST /post`                | `title`, `content`, one to three `category[]` IDs, and CSRF token.         |
| `GET /post/edit?post_id=ID` | Render an author’s existing discussion and its current revision.           |
| `POST /post/edit`           | Composer fields plus `post_id` and the hidden `revision` from that editor. |
| `POST /add_comment`         | `post_id`, `content`, and CSRF token.                                      |

Titles allow 5–120 characters, discussion bodies 10–5,000, and replies 1–2,000 after trimming. Topic IDs must exist. Successful writes redirect with HTTP 303. Validation errors retain composer values and return HTTP 400. Only the author may open or submit the editor; ownership failures return HTTP 403.

An edit updates the discussion and topic links atomically while preserving its ID, publication date, replies, reactions, and bookmarks. If the submitted revision is stale, HTTP 409 keeps the submitted text in the editor. Open the current discussion and compare before applying changes through a fresh editor. Revisions are counters, not recoverable edit history.

The browser tracks unsaved editor changes and requests a native leave-page confirmation where supported. Restoring the original text/topics removes that guard. This does not save a draft or guarantee recovery after a closed tab, crash, or mobile app termination.

## Private bookmarks

`POST /bookmarks` accepts `post_id`, `action=save|remove`, and a CSRF token. Native forms redirect to the discussion with HTTP 303. Requests with `Accept: application/json` receive:

```json
{ "saved": true }
```

Removal returns `false`. Explicit actions are idempotent: repeating a save neither duplicates nor reorders it. Removing and saving again puts it at the front of the private collection. An anonymous JSON request with a valid CSRF token gets HTTP 401; the native form redirects to sign-in. Invalid inputs return 400 and missing discussions return 404.

## Reactions

`POST /like_dislike` accepts JSON with `postId`, `type` (`post` or `comment`), and `action` (`like` or `dislike`). For comments, `postId` contains the comment ID. Send the same-origin session cookie and `X-CSRF-Token` header.

```json
{ "postId": 5, "type": "post", "action": "like" }
```

A successful response contains `likeCount`, `dislikeCount`, and `reaction` (`1` like, `0` dislike, `-1` none). Repeating the currently selected reaction removes it; switching replaces it. Unlike bookmark saves, reaction toggles must not be retried automatically after an uncertain response. Refresh the discussion to establish the current state.

Unknown JSON fields, trailing JSON, unsupported types/actions, and invalid IDs are rejected. Anonymous requests return 401; missing targets return 404. Reactions are serialized in a transaction.

## Error handling

Route-level JSON failures use an `error` message. Shared security middleware can return an HTML error page, including CSRF rejection, so clients must check HTTP status before assuming a JSON body. Authentication throttling returns HTTP 429 with `Retry-After: 60`. Use the form’s displayed feedback for recovery; never infer success solely from a completed network request.
