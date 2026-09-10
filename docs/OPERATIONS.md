# Operations

## Configuration

| Setting                  | Default           | Purpose                                                                   |
| ------------------------ | ----------------- | ------------------------------------------------------------------------- |
| `TALKNET_ADDR`           | `:8080`           | Native server listen address; keep `:8080` inside the supplied container. |
| `TALKNET_DB`             | `data/talknet.db` | Native SQLite path. The container sets `/app/data/talknet.db`.            |
| `TALKNET_SECURE_COOKIES` | `false`           | Set to `true` when users connect over HTTPS.                              |
| `TALKNET_PORT`           | `8088`            | Host port used by Docker Compose.                                         |

`-seed-demo` initializes an empty community with fictional records and exits. It does nothing if users or posts already exist. `-seed-showcase` explicitly adds the larger fictional pack to either an empty or populated database. Its `showcase-v1` marker in `Demo_Seeds` prevents duplicate imports. All records and the marker commit in one transaction; account-name/email collisions fail without adopting or modifying the existing account. `-healthcheck` calls the local health endpoint and exits with a nonzero status when it is unavailable.

## Existing Talknet data

The repository no longer distributes a working database. Existing local `talknet.db` files are ignored by Git and are not removed by the application. The original database schema is supported; a backup is recommended before any software upgrade.

For a native run, stop the old application and point Talknet at the existing file:

```powershell
$env:TALKNET_DB = '.\talknet.db'
go run .
```

```bash
TALKNET_DB=./talknet.db go run .
```

For Docker, mount a directory containing your existing database at `/app/data`. The runtime user is UID/GID `10001`; the directory and database must be writable by that user on Linux. Keep the database name `talknet.db` or update `TALKNET_DB`. A bind mount replaces the default named volume for that deployment; stop the old app before switching. Do not run two application versions against the same file during migration.

The original table names and relationships are retained. Startup creates missing tables and indexes and inserts only missing standard topics. It does not replace existing posts, comments, users, or passwords. The additive `Bookmarks` and `Post_Revisions` tables are created automatically; existing posts begin at revision one without rewriting their content. Existing bcrypt password hashes remain valid. Historical sessions must sign in again because session tokens are now hashed and expiry is enforced by the server.

## Persistence and backups

`docker compose stop` and ordinary `docker compose down` preserve the database volume. `docker compose down -v` deletes named volumes and therefore the forum’s data; it is not part of the normal run instructions.

For a native database, use SQLite’s online backup API or stop the application before copying the database and any accompanying WAL files. For a named volume, stop Talknet, back up the volume through your normal Docker backup tooling, and then restart it. Test restoration before relying on a backup.

## Public deployment

The supplied Compose configuration binds only to `127.0.0.1`. To host the forum for other people, put it behind an HTTPS reverse proxy, configure the proxy route, and set `TALKNET_SECURE_COOKIES=true`. Proxy traffic to the configured local port. Keep SQLite on persistent storage and retain the single-instance deployment model.

The built-in authentication throttle uses the socket peer address. A reverse proxy can make many users share that address; add an appropriate per-client rate limit at the proxy and account for the application’s ten-attempt-per-minute limit. Forwarded headers are not accepted as trusted identity information.

Before operating a public community, provide the moderation, account-recovery, and privacy processes your audience needs. The repository does not include an email service or moderation console.

## Troubleshooting

- **Port already in use:** change `TALKNET_PORT` in `.env`, then recreate the container.
- **Native SQLite build fails:** `go-sqlite3` requires CGO and a C compiler. Use Docker if you do not want a local compiler setup.
- **Sign-in does not persist over local HTTP:** keep `TALKNET_SECURE_COOKIES=false` locally. A secure cookie requires HTTPS.
- **Form expired:** refresh the page to obtain a current CSRF token, then retry.
- **Database is read-only:** check ownership and write access for the database directory, including WAL files.
- **Seeded authors cannot sign in:** these are fictional demonstration records. Register your own account.
- **Auth attempts are throttled:** wait one minute. `429` responses include `Retry-After: 60`.
- **A changed asset is not visible:** assets are embedded at build time. Rebuild and restart the binary or container.

- **An edit changed in another tab:** the server keeps your submitted text and returns a conflict. Open the current discussion in a separate tab, compare the versions, and apply your changes through a freshly opened editor. There is no draft autosave or recoverable revision history.
