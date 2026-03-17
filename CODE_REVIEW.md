# Code Review Report — Tundra

**Date:** 2026-03-16
**Reviewer:** Claude Opus 4.6

---

## Critical (4)

### 1. No transaction in migration
**File:** `internal/database/migrations.go`

The `CREATE TABLE` statements run as individual `Exec` calls. If the second table creation fails, the first table is left behind with no rollback, leaving the database in an inconsistent state.

### 2. No transaction in `UpdateLocation`
**File:** `internal/database/database.go`

`UpdateLocation` does a SELECT then an UPDATE without a transaction. Under concurrent access, another request could modify or delete the row between the two queries (TOCTOU race).

### 3. No transaction in `CreateItem`/`UpdateItem`
**File:** `internal/database/database.go`

These insert/update an item and then insert stock rows in separate statements. If stock insertion fails partway through, you get an item with partial stock data and no rollback.

### 4. N+1 query in `ListItems`
**File:** `internal/database/database.go`

`ListItems` fetches all items, then loops and runs a separate query per item to fetch stock. This will degrade linearly as inventory grows. Use a JOIN or a single batch query for stock.

---

## Warnings (9)

### 5. No foreign key on `stock.location`
The `stock` table stores `location TEXT` instead of referencing `locations.id`. This allows orphaned/invalid location references and makes renames inconsistent.

### 6. Orphaned stock on item delete
`DeleteItem` only deletes the item row. Stock rows referencing that item are left behind (no `ON DELETE CASCADE` and no manual cleanup).

### 7. Error messages leak internals
API handlers return raw `err.Error()` to the client. This can expose SQL errors, file paths, or internal state to end users.

### 8. No request body size limit
`http.ListenAndServe` with no `MaxBytesReader` means a client can send an arbitrarily large JSON body and exhaust server memory.

### 9. Missing `Content-Type` validation
Handlers don't check that incoming requests have `Content-Type: application/json` before attempting `json.Decode`.

### 10. `go.mod` declares Go 1.24.1
This pins to a patch version that may not exist on all build machines. Prefer `go 1.24` (minor only).

### 11. Unused `radix-vue` dependency
`frontend/package.json` lists `radix-vue` but no imports of it were found in the component source.

### 12. Falsy-value bug in `ItemDrawer.vue`
Uses `||` for defaults (e.g., `props.item.quantity || ''`). If `quantity` is `0`, it's falsy and gets replaced with `''`, silently losing valid zero values. Use `??` (nullish coalescing) instead.

---

## Suggestions (11)

### 13. Missing indexes
No index on `stock.item_id` or `stock.location`. As data grows, JOIN and filter queries will slow down.

### 14. Timer/interval leak in `StatusBar.vue`
If the component uses `setInterval` for a clock or status refresh, ensure it's cleared in `onUnmounted` to prevent leaks on navigation.

### 15. No loading/error states in views
Views fetch data but don't show loading spinners or error messages if the API is slow or fails. Users see a blank screen.

### 16. Accessibility gaps
Drawer and modal components lack `aria-*` attributes, focus trapping, and keyboard dismiss (`Escape`). Screen readers and keyboard users can't interact properly.

### 17. Animation loss from `v-if` mount/unmount
Drawer uses `v-if` which destroys the DOM node. Transition animations on close won't play. Consider `v-show` or `<Transition>` wrapper.

### 18. No input validation on the frontend
Forms allow submitting empty names, negative quantities, etc. Validate before sending to the API.

### 19. Docker: no `.dockerignore` for Go cache
The `.dockerignore` exists but may not exclude Go build cache or test artifacts, increasing image size.

### 20. No graceful shutdown
`main.go` calls `http.ListenAndServe` without signal handling. On container stop, in-flight requests are dropped instead of draining.

### 21. Embed directive and dev mode
`embed.go` embeds the frontend dist. Ensure the dev workflow doesn't accidentally serve stale embedded assets when the frontend hasn't been rebuilt.

### 22. No CORS configuration
If the frontend dev server runs on a different port, API calls will be blocked by the browser. No CORS headers are set on the Go server.

### 23. Route handler organization
All routes and handlers live in `main.go`. As the API grows, extract handlers into an `internal/handler` or `internal/api` package.

---

## Recommended Priority

1. **Items 1–4** (transactions + N+1 query) — risk data corruption and performance degradation
2. **Item 12** (falsy-value bug) — silently corrupts user input
3. **Item 6** (orphaned stock) — data integrity
4. **Items 7–9** (security hardening) — defense in depth
5. Remaining suggestions as time permits