# Cars24 AI Operations Copilot

An AI-powered backend that lets an operations team ask natural-language
questions about orders — payment status, delivery status, or a full
order summary — and get answers grounded in real database records, not
LLM guesses.

```
"Customer says they've paid for order #1289 but delivery isn't
scheduled — what's going on?"

→ Payment for order #1289 is confirmed (paid), but delivery has not
  been scheduled yet. This looks like a delivery-scheduling gap rather
  than a payment issue — the order should be escalated to logistics.
```

A minimal web UI is included so the service can be demoed without curl,
but the core deliverable is the API: `POST /query`.

---

## Stack

| Layer      | Choice                                |
|------------|----------------------------------------|
| Language   | Go 1.25                                |
| HTTP       | Gin                                    |
| ORM        | GORM (SQLite driver)                   |
| LLM access | OpenRouter (`/v1/chat/completions`, OpenAI-compatible tool-calling) |
| Frontend   | Static HTML/CSS/JS, served by the same Go binary |

No external services required beyond an OpenRouter API key — the
database is a local SQLite file.

---

## Project structure

```
cmd/server/main.go       — entrypoint: wiring, routes, HTTP server
internal/config/         — env var loading
internal/database/       — SQLite connection + GORM AutoMigrate
internal/models/         — Customer, Vehicle, Order, Payment, Delivery
internal/tools/          — one Go type per LLM tool (order/payment/delivery/summary)
internal/ai/             — LLM types, OpenRouter client, tool-calling agent loop, system prompt
internal/handlers/       — POST /query HTTP handler
seed/seed.go              — standalone seed script (realistic demo data)
web/                      — static frontend (index.html, style.css, app.js)
```

---

## Setup

### Prerequisites

- Go 1.25 or later
- An [OpenRouter](https://openrouter.ai) API key (free tier works — the
  default model is `openrouter/free`)

### 1. Clone and install dependencies

```bash
git clone <your-repo-url>
cd AI-Operations-Copilot
go mod download
```

### 2. Configure environment variables

Create a `.env` file in the project root:

```bash
PORT=8080
DATABASE_PATH=./data/cars24.db
MODEL=openrouter/free
API_KEY=your-openrouter-api-key-here
```

| Variable        | Required | Default                | Description                                  |
|------------------|:--------:|-------------------------|-----------------------------------------------|
| `PORT`           | No       | `8080`                  | HTTP port the server listens on               |
| `DATABASE_PATH`  | No       | `./data/cars24.db`      | Path to the SQLite database file              |
| `MODEL`          | No       | `openrouter/free`       | OpenRouter model ID used for tool-calling      |
| `API_KEY`        | **Yes**  | —                       | OpenRouter API key. Server refuses to start without it |

> **Never commit `.env`.** It's already listed in `.gitignore` — keep it
> that way, and rotate any key that's ever been pasted outside your own
> machine (e.g. into a chat, ticket, or screenshot).

### 3. Seed the database

```bash
mkdir -p data
go run ./seed
```

This clears and repopulates the database with 15 orders covering every
status combination (paid+undelivered, payment failed, out for
delivery, cancelled, etc.) — see [DESIGN.md](./DESIGN.md#seed-data)
for the full scenario table.

### 4. Run the server

```bash
go run ./cmd/server
```

You should see:

```
Cars24 AI Operations Copilot listening on port 8080
Using OpenRouter model: openrouter/free
```

### 5. Open the app

Visit **http://localhost:8080** — the frontend is served directly from
the Go binary. Try one of the suggested queries, or click a seeded
order ID to auto-fill a question about it.

No separate frontend server or build step is needed.

---

## API Documentation

### `GET /health`

Liveness check.

**Response `200 OK`**
```json
{ "status": "ok" }
```

---

### `POST /query`

Ask a natural-language operations question. The agent decides which
internal tools to call (order / payment / delivery / full summary),
executes them against the database, and returns a grounded answer.

**Request**

```http
POST /query
Content-Type: application/json
```

```json
{ "query": "Customer says they've paid for order #1289 but delivery isn't scheduled — what's going on?" }
```

| Field   | Type   | Required | Constraints                     |
|---------|--------|:--------:|----------------------------------|
| `query` | string | Yes      | Non-empty, max 2000 characters   |

**Response `200 OK`**

```json
{
  "answer": "Payment for order #1289 is confirmed (paid), but delivery has not been scheduled yet. This is a scheduling gap, not a payment issue.",
  "tool_calls": [
    {
      "name": "get_payment_status",
      "arguments": "{\"order_id\":1289}",
      "success": true,
      "duration_ms": 4
    },
    {
      "name": "get_delivery_status",
      "arguments": "{\"order_id\":1289}",
      "success": true,
      "duration_ms": 3
    }
  ]
}
```

`tool_calls` is the full, ordered trace of every tool the agent invoked
while producing the answer — the mechanism that makes "grounded, not
guessed" verifiable rather than just claimed. Each entry:

| Field | Type | Description |
|-------|------|--------------|
| `name` | string | Tool name, e.g. `get_payment_status` |
| `arguments` | string | Raw JSON arguments sent to the tool |
| `success` | boolean | Whether the tool call succeeded |
| `error` | string | Present only when `success` is `false` (e.g. order not found) |
| `duration_ms` | number | Tool execution time in milliseconds |

The frontend renders this trace above each answer so it's visible
which tools were checked (and whether any failed) before the model
responded.

**Error responses**

| Status | Body | Cause |
|--------|------|-------|
| `400 Bad Request` | `{"error": "query is required"}` | Missing/empty `query`, or `query` exceeds 2000 characters, or malformed JSON |
| `502 Bad Gateway`  | `{"error": "AI service temporarily unavailable"}` | Upstream LLM call failed after retries, or the agent exceeded its tool-calling round limit |

**Example**

```bash
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{"query": "Give me a full status summary for order #2231."}'
```

#### Tools available to the agent

| Tool | Purpose |
|------|---------|
| `get_order_status` | Order status, order date, linked customer/vehicle IDs |
| `get_payment_status` | Payment status, amount, paid-at timestamp |
| `get_delivery_status` | Delivery status, scheduled date, delivered-at timestamp |
| `get_full_order_summary` | Order + customer + vehicle + payment + delivery, joined |

Each tool takes a single `order_id` (integer) argument. See
[DESIGN.md](./DESIGN.md#tool-design) for why they're split this way
instead of one catch-all tool.

---

## Seeded reference data

For manual testing, the seed script guarantees these order IDs exist
with specific, useful scenarios:

| Order ID | Payment | Delivery | Useful for testing |
|----------|---------|----------|----------------------|
| `#1289`  | Paid    | Not scheduled | The core diagnostic query (paid but undelivered) |
| `#2231`  | Paid    | Delivered | Full summary on a clean, completed order |
| `#4521`  | Pending | Scheduled | Delivery scheduled ahead of payment — inverse edge case |
| `#3345`  | Failed  | Not scheduled | Payment-failure path |
| `#5678`  | Paid    | Out for delivery | Mid-flight delivery state |
| `#9345`  | Failed  | Not scheduled | Cancelled order |
| `#99999` | —       | — | Does not exist — tests "order not found" handling |

---

## Notes

- The service is **read-only** — it answers questions but never
  performs actions (no cancellations, refunds, or rescheduling). This
  is a deliberate scope boundary; see DESIGN.md for what a write path
  would need.
- CORS is enabled (`Access-Control-Allow-Origin: *`) so the frontend
  can also be run from a separate dev server if you'd rather not use
  the bundled one.
