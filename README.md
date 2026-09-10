# Talknet

**Good ideas. Better conversations.**

A full-stack discussion community built with **Go, SQLite, server-rendered HTML, and vanilla JavaScript**. Browse topics, share a perspective, and keep a conversation going—with a distinctive editorial interface and a small, self-contained runtime.

[![Quality](https://github.com/hujaafar/Talknet/actions/workflows/ci.yml/badge.svg)](https://github.com/hujaafar/Talknet/actions/workflows/ci.yml)
[![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-d8ef6d.svg)](LICENSE)

![Original chrome and citrus conversation artwork for Talknet](static/images/conversation-art.webp)

## The experience

| Explore                              | Participate                                 | Make it yours                                 |
| ------------------------------------ | ------------------------------------------- | --------------------------------------------- |
| Search titles and discussion content | Create discussions with up to three topics  | Register and sign in                          |
| Filter by topic                      | Reply to posts                              | View your profile and authored posts          |
| Sort by newest or most liked         | Like, dislike, switch, or remove a reaction | Revisit your liked discussions                |
| Browse paginated feeds               | React to individual replies                 | Pause motion or follow your system preference |

The interface pairs graphite surfaces, citrus accents, local Inter typography, and original 3D artwork. It includes scroll-linked artwork, section reveals, a reading-progress indicator, responsive layouts, visible keyboard focus, and reduced-motion support. Reading, search, navigation, authentication, posting, and replies work without JavaScript; reactions and other enhancements use a small client script.

## Run with Docker

Install Docker with Compose v2 and start its Linux container engine.

```bash
git clone https://github.com/hujaafar/Talknet.git
cd Talknet
docker compose up -d --build --wait
```

Open **[localhost:8088](http://localhost:8088)**. Create an account to participate. The database is stored in the `talknet-data` named volume and survives rebuilds and ordinary shutdowns.

For a populated local demo:

```bash
docker compose run --rm talknet -seed-demo
```

Refresh the page. Seeding only runs when both the users and posts tables are empty. The four fictional authors are labeled `Demo`; their accounts have no shared login password. Sample records are never added to an existing community. Register your own account to try posting and reactions.

On Windows, the helper can build, seed, and start everything:

```powershell
.\scripts\start.ps1 -Demo
```

Useful commands:

```bash
docker compose logs -f
docker compose stop
docker compose up -d --wait
```

To choose another port, copy `.env.example` to `.env` and change `TALKNET_PORT`. The default port binding is local to your machine.

## Run with Go

Requires **Go 1.26+** and a C compiler for `go-sqlite3`. On Windows, use a compatible GCC toolchain or the Docker workflow above.

```bash
go mod download
go run . -seed-demo   # optional; seeds an empty database and exits
go run .
```

Open **[localhost:8080](http://localhost:8080)**. Native runs use `data/talknet.db` by default. Templates, CSS, JavaScript, imagery, fonts, and the schema are embedded in the binary.

## Quality and implementation

- Persistent, expiring sessions with hashed bearer tokens and server-side logout.
- Bcrypt passwords, validated forms, bounded requests, and authentication throttling.
- CSRF tokens, cross-origin write protection, secure-cookie configuration, and a restrictive Content Security Policy.
- Parameterized SQL, atomic post/category and reaction writes, foreign keys, WAL, and indexed queries.
- A multi-stage Docker image running as a non-root user, with a health check, persistent data, resource limits, and a read-only root filesystem.
- Integration tests for real account, post, reply, reaction, search, profile, pagination, session, and validation flows—including concurrent reaction toggles under Go’s race detector.

```bash
go test -race -cover ./...
go vet ./...
go build -trimpath ./...
```

Or run tests entirely in Docker:

```bash
docker build --target test -t talknet-tests .
```

[GitHub Actions](https://github.com/hujaafar/Talknet/actions) checks formatting, race/integration tests, static analysis, compilation, and container startup. See [validation notes](docs/VALIDATION.md) for the scope of the recorded checks.

## Inside the project

```text
main.go                 Configuration, health checks, graceful shutdown
internal/forum/         Routing, authentication, queries, writes, and tests
  schema.sql            Idempotent schema initialization and indexes
  seed.go               Explicit, non-destructive demo setup
static/
  assets.go             Embedded interface assets
  pages/                Shared Go templates and seven page views
  styles/app.css        Responsive visual system
  js/app.js             Motion, reactions, form feedback, and enhancements
  images/               Original artwork, local typeface, and favicon
compose.yaml            Persistent local deployment
docs/                   Architecture, operations, validation, and credits
```

Read [the architecture](docs/ARCHITECTURE.md), [configuration and existing-data migration](docs/OPERATIONS.md), or [artwork and font credits](docs/CREDITS.md).

## Project lineage

Originally created by [Mohamed Alasfoor](https://github.com/Mohamed-Alasfoor), [Ali Hasan](https://github.com/alihjmm), [Habib Mansoor](https://github.com/7abib04), and [Hussain Jawad](https://github.com/hujaafar).

The redesign retains Talknet’s Go/SQLite foundation and original contribution history, with a new visual identity and a consolidated application layer. Licensed under [MIT](LICENSE); the bundled typeface has its own [SIL Open Font License](docs/licenses/Inter-OFL.txt).
