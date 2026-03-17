# Backstock

A household food and non-food inventory tracker. Built to replace a shared spreadsheet with something faster and easier to use on a phone.

## Stack

- **Backend:** Go, Chi router, SQLite (WAL mode)
- **Frontend:** Vue 3, Vite, Bun, Tailwind CSS
- **Deployment:** Docker with a named volume for the database

## Features

- Track items across multiple storage locations (fridge, freezer, pantry, etc.)
- Partial quantity moves between locations
- Decrease/replenish workflow so items persist at zero stock
- Low-quantity alerts per item
- Categories (tag-style, multiple per item)
- Customizable units
- Full-text search and filtering by location, category, or sort order
- Food and non-food mode toggle
- Mobile-first UI

## Running with Docker

```sh
docker compose up -d
```

The app will be available at `http://localhost:8080`. The SQLite database is stored in a Docker volume at `/data/backstock.db`.

## Development

### Prerequisites

- Go 1.25+
- Bun (for the frontend)

### Frontend

```sh
cd frontend
bun install
bun run dev
```

The Vite dev server runs on `http://localhost:5173` and proxies API requests to the Go backend.

### Backend

```sh
go run .
```

Starts the server on `:8080`. Set `BACKSTOCK_ADDR` to change the listen address and `BACKSTOCK_DB_PATH` to change the database file location (defaults to `backstock.db` in the working directory).

### Building

```sh
cd frontend && bun run build && cd ..
go build -o backstock .
```

The frontend is embedded into the Go binary at build time via `//go:embed`.

## Configuration

| Variable | Default | Description |
|---|---|---|
| `BACKSTOCK_ADDR` | `:8080` | Listen address |
| `BACKSTOCK_DB_PATH` | `backstock.db` | Path to the SQLite database file |

## Architecture

The app is a single Go binary serving both the API and the SPA frontend.

- `main.go` / `embed.go` - Entry point and frontend embedding
- `internal/migrate/` - SQL migrations (run automatically on startup)
- `internal/store/` - Database access layer
- `internal/handler/` - HTTP handlers
- `internal/server/` - Router and middleware setup
- `internal/model/` - Shared data types
- `frontend/` - Vue 3 SPA

Items and stock are separate tables. An "item" is a definition (name, unit, category, thresholds). Stock rows represent quantities of that item at specific locations. This allows partial moves and per-location tracking.

No authentication is built in. The app is designed to sit behind a reverse proxy with OIDC or similar.
