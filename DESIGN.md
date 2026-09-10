# Design — Cars24 AI Operations Copilot

## 1. Problem framing

The ops team needs to ask questions that span multiple systems of
record — order, payment, delivery — and sometimes needs the model to
**reason across them**, not just fetch one value. The example query in
the brief makes this explicit:

> "Customer says they've paid for order #1289 but delivery isn't
> scheduled — what's going on?"

Answering this correctly requires checking payment status *and*
delivery status independently, then explaining the discrepancy — the
model cannot answer it from a single lookup, and it must not assume
that "paid" implies "scheduled." That constraint drove most of the
design decisions below.

## 2. Architecture

```
                     ┌──────────────┐
   Browser / curl ─▶ │  POST /query │  (Gin handler)
                     └──────┬───────┘
                            │
                     ┌──────▼───────┐
                     │  ai.Agent    │  tool-calling loop (max 8 rounds)
                     └──────┬───────┘
                            │ tool_use / tool_result
                     ┌──────▼───────┐
                     │ tools.Registry│  routes by tool name → typed struct
                     └──────┬───────┘
                            │
          ┌─────────┬───────┴────────┬────────────┐
          ▼         ▼                ▼            ▼
     OrderTool  PaymentTool   DeliveryTool   SummaryTool
          └─────────┴───────┬────────┴────────────┘
                            ▼
                       GORM / SQLite
```

The LLM never talks to the database directly. It only sees tool
definitions (JSON schemas) and tool results (JSON). This is what makes
the "never invent data" guarantee enforceable rather than just
requested.

## 3. Tool design

Four narrow tools instead of one broad "get order info" tool:

- `get_order_status(order_id)`
- `get_payment_status(order_id)`
- `get_delivery_status(order_id)`
- `get_full_order_summary(order_id)` — a convenience join of all four,
  used for "give me a full status summary" style queries

**Why split them up:** the diagnostic query requires the model to
*choose* to check two independent domains and reconcile them. If
there were only one all-in-one tool, the model would never have to
demonstrate that reasoning — it would just dump every field back at
the user, and a genuinely inconsistent state (e.g. paid, but
delivery record missing) would be easy to paper over. Narrow tools
force explicit, auditable multi-step reasoning, which also shows up in
the server logs (`AI tool call: name=... arguments=...`) for
debugging.

**Why keep `get_full_order_summary` anyway:** without it, "give me a
full status summary" would require the model to make 3 separate calls
every time, which is slower and gives it more chances to only check
some of the tools. A single joined-read tool for the "give me
everything" case is a deliberate shortcut, not a contradiction of the
point above.

## 4. Agent loop

`internal/ai/agent.go` implements a standard tool-calling loop:

1. Send the system prompt + user query + tool definitions to the LLM.
2. If the response contains `tool_calls`, execute each one against the
   registry, append the results as `tool` role messages, and go back
   to step 1.
3. If the response has no tool calls, treat its text as the final
   answer.
4. Hard cap of 8 rounds (`Agent.MaxRounds`) so a confused model can't
   loop forever or run up API costs.

This supports the two-call chain the diagnostic query needs (payment
status → delivery status) without any special-casing — it falls out
of the general loop.

**Tool-call trace.** Every call in the loop is recorded (`name`,
`arguments`, `success`, `error`, `duration_ms`) and returned to the
client alongside the final answer as `tool_calls` in the `/query`
response, not just logged server-side. The rationale: "the model
won't hallucinate" is a claim that's easy to assert and hard to
verify from the outside. Surfacing the actual trace — which tools
ran, with what arguments, in what order — turns that claim into
something a reviewer (or an ops agent) can check on every single
response, and it's what the frontend renders above each answer.

## 5. Grounding / anti-hallucination strategy

This was the highest-risk part of the assignment: an LLM that
confidently invents an order status is worse than no tool at all. Two
layers address it:

1. **System prompt rules** (`internal/ai/agent.go`) explicitly forbid
   inventing payment/delivery/customer/vehicle data, forbid treating
   the user's claims as facts ("customer says they paid" must still be
   verified via the tool, not taken at face value), forbid using one
   order's tool result to answer about a different order, and require
   reporting inconsistencies rather than guessing when tool results
   conflict.
2. **Tool results are the only source of order data** the model ever
   sees — there is no path for the model to answer an order-specific
   question without a tool call returning real rows from SQLite.

`temperature: 0.1` is also set on every request to reduce variance in
borderline cases (it doesn't replace the rules above, it just makes
outputs more consistent run to run).

## 6. Data model

```
Customer 1───* Order *───1 Vehicle
                 │
                 ├──1 Payment
                 └──1 Delivery
```

Payment and Delivery are separate tables (not columns on Order)
specifically so they can go out of sync — the diagnostic query only
means something if payment and delivery are independent facts that
*can* disagree. `Payment.Status` and `Delivery.Status` are plain
strings with named constants rather than an enum type, matching
GORM's straightforward SQLite mapping without extra migration
complexity.

## 7. Seed data

`seed/seed.go` seeds 15 orders across 10 customers and 10 vehicles,
covering every state combination needed to make the tool's reasoning
demonstrable rather than trivial:

| Scenario | Payment | Delivery | Order IDs |
|---|---|---|---|
| Clean, delivered | paid | delivered | `#2231`, `#9012` |
| **Paid, not scheduled** (core diagnostic case) | paid | not_scheduled | `#1289` |
| Pending payment, delivery already scheduled (inverse edge case) | pending | scheduled | `#4521` |
| Payment failed | failed | not_scheduled | `#3345` |
| Out for delivery | paid | out_for_delivery | `#5678` |
| Normal in-flight | paid | scheduled | `#6789`, `#8901`, `#9123`, `#9567`, `#9678`, `#9789` |
| Pending, nothing scheduled | pending | not_scheduled | `#7890`, `#9456` |
| Cancelled | failed | not_scheduled | `#9345` |
| Non-existent order (not seeded) | — | — | `#99999` — tests the "order not found" path |

The seed script is idempotent — it clears child tables before parent
tables and can be re-run safely.

## 8. Error handling

- **Invalid input:** empty query, query >2000 chars, or malformed JSON
  → `400` before the LLM is ever called.
- **Order not found:** the tool layer returns a descriptive error
  (`"order %d not found"`) which is serialized and handed back to the
  model as a `tool` message, so the model can tell the user the order
  doesn't exist instead of silently failing or hallucinating one.
- **Upstream LLM failures:** `OpenRouterProvider.Chat` retries
  transient failures (HTTP 429/5xx, and provider errors whose message
  mentions "temporarily"/"overloaded") up to 3 times with exponential
  backoff, then surfaces a `502` to the client.
- **Runaway loops:** capped at 8 tool-calling rounds; exceeding it
  returns a `502` rather than hanging the request. The handler also
  applies a 45s context timeout independent of the round cap.

## 9. Frontend

A static HTML/CSS/JS page served by the same Go binary (`web/`,
mounted via `router.Static`) — no build step, no separate deployment.
It's a thin client over `POST /query`: a scrolling transcript, an
input box, and a reference sidebar listing the seeded order IDs and
their scenarios so a reviewer can test the diagnostic case immediately
without reading the seed script first. Each answer is preceded by its
tool-call trace (tool name, arguments, success/fail, latency) rendered
directly from the `tool_calls` array in the response — so the
"grounded, not guessed" claim is visible per-query, not just asserted
in this document. All grounding and reasoning logic lives in the
backend; the frontend has no business logic of its own.

## 10. Known limitations & what a v2 would add

- **Read-only.** The agent can't take action (reschedule delivery,
  retry payment, cancel an order). Adding this safely would need a
  second tool category with explicit confirmation semantics — an LLM
  should not autonomously mutate financial records — plus an audit
  log of who approved what.
- **No conversation memory.** Each `/query` call is stateless; there's
  no follow-up ("...and what about the customer's other orders?")
  without repeating full context. A session/thread ID with stored
  message history would fix this.
- **No auth.** Anyone who can reach the port can query any order. A
  real deployment needs at minimum an API key or SSO-gated access,
  since this surfaces customer PII (name, phone, email).
- **Single-tenant SQLite.** Fine for a take-home; a production version
  would move to Postgres and add connection pooling.
- **No automated tests.** Given the time constraint, verification was
  manual against the seed scenarios above. Unit tests for the tool
  layer (especially the not-found and conflicting-state paths) and an
  integration test that stubs the LLM provider would be the first
  addition.
- **No observability beyond logs.** Tool-call latency is logged per
  call, but there's no tracing/metrics export — useful once multiple
  tool calls per query start adding up in production.
