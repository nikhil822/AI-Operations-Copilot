const QUERY_ENDPOINT = "/query";

const SUGGESTIONS = [
  "What's the payment status for order #4521?",
  "Customer says they've paid for order #1289 but delivery isn't scheduled — what's going on?",
  "Give me a full status summary for order #2231.",
  "Why is order #3345 stuck?",
  "What's the status of order #99999?",
];

// Mirrors seed/seed.go — kept in the frontend purely as a demo/testing
// reference for whoever is trying the tool. The backend is the only
// source of truth; this list is never sent to the LLM.
const SEEDED_ORDERS = [
  { id: 1289, tag: "paid / not scheduled", cls: "amber" },
  { id: 2231, tag: "delivered", cls: "green" },
  { id: 4521, tag: "pending / scheduled", cls: "amber" },
  { id: 3345, tag: "payment failed", cls: "red" },
  { id: 5678, tag: "out for delivery", cls: "green" },
  { id: 6789, tag: "paid / scheduled", cls: "green" },
  { id: 7890, tag: "pending / not scheduled", cls: "amber" },
  { id: 9345, tag: "cancelled", cls: "red" },
];

const transcript = document.getElementById("transcript");
const composer = document.getElementById("composer");
const queryInput = document.getElementById("queryInput");
const submitBtn = document.getElementById("submitBtn");
const statusDot = document.getElementById("statusDot");
const chipsContainer = document.getElementById("suggestionChips");
const orderListContainer = document.getElementById("orderList");

function renderSuggestions() {
  SUGGESTIONS.forEach((text) => {
    const chip = document.createElement("button");
    chip.type = "button";
    chip.className = "chip";
    chip.textContent = text;
    chip.addEventListener("click", () => {
      queryInput.value = text;
      queryInput.focus();
    });
    chipsContainer.appendChild(chip);
  });
}

function renderOrderList() {
  SEEDED_ORDERS.forEach((order) => {
    const row = document.createElement("div");
    row.className = "order-row";
    row.innerHTML = `
      <span class="order-row__id">#${order.id}</span>
      <span class="order-row__tag tag--${order.cls}">${order.tag}</span>
    `;
    row.addEventListener("click", () => {
      queryInput.value = `Give me a full status summary for order #${order.id}.`;
      queryInput.focus();
    });
    orderListContainer.appendChild(row);
  });
}

function stripMarkdown(text) {
  return text
    .replace(/\*\*(.*?)\*\*/g, "$1")   // **bold**
    .replace(/\*(.*?)\*/g, "$1")        // *italic*
    .replace(/`([^`]*)`/g, "$1")        // `code`
    .replace(/^#{1,6}\s+/gm, "")        // # headers
    .replace(/^[-*]\s+/gm, "");         // - bullet / * bullet
}

function addEntry(label, body, variant) {
  const entry = document.createElement("div");
  entry.className = `entry entry--${variant}`;

  const labelEl = document.createElement("div");
  labelEl.className = "entry__label";
  labelEl.textContent = label;

  const bodyEl = document.createElement("div");
  bodyEl.className = "entry__body";
  bodyEl.textContent = variant === "assistant" ? stripMarkdown(body) : body;

  entry.appendChild(labelEl);
  entry.appendChild(bodyEl);
  transcript.appendChild(entry);
  transcript.scrollTop = transcript.scrollHeight;

  return entry;
}

// Renders the sequence of tool calls the agent made while answering,
// e.g. "get_payment_status(order_id: 1289) — 42ms". This is what lets
// someone verify the answer was actually grounded in tool results
// rather than just trusting the text.
function addToolTrace(toolCalls) {
  if (!toolCalls || toolCalls.length === 0) return;

  const wrapper = document.createElement("div");
  wrapper.className = "entry entry--trace";

  const labelEl = document.createElement("div");
  labelEl.className = "entry__label";
  labelEl.textContent = `tools used (${toolCalls.length})`;
  wrapper.appendChild(labelEl);

  const list = document.createElement("div");
  list.className = "trace-list";

  toolCalls.forEach((call) => {
    const row = document.createElement("div");
    row.className = `trace-row ${call.success ? "trace-row--ok" : "trace-row--fail"}`;

    let args = call.arguments;
    try {
      const parsed = JSON.parse(call.arguments);
      args = Object.entries(parsed)
        .map(([k, v]) => `${k}: ${v}`)
        .join(", ");
    } catch (_) {
      // fall back to raw string
    }

    row.innerHTML = `
      <span class="trace-row__status">${call.success ? "✓" : "✗"}</span>
      <span class="trace-row__name">${call.name}</span>
      <span class="trace-row__args">(${args})</span>
      <span class="trace-row__duration">${call.duration_ms}ms</span>
    `;

    if (!call.success && call.error) {
      const errEl = document.createElement("div");
      errEl.className = "trace-row__error";
      errEl.textContent = call.error;
      row.appendChild(errEl);
    }

    list.appendChild(row);
  });

  wrapper.appendChild(list);
  transcript.appendChild(wrapper);
  transcript.scrollTop = transcript.scrollHeight;
}

function setBusy(isBusy) {
  submitBtn.disabled = isBusy;
  queryInput.disabled = isBusy;
  statusDot.classList.toggle("is-busy", isBusy);
  if (!isBusy) statusDot.classList.remove("is-error");
}

async function submitQuery(query) {
  addEntry("you", query, "user");
  const thinkingEntry = addEntry("copilot", "Checking order records…", "thinking");
  setBusy(true);

  try {
    const response = await fetch(QUERY_ENDPOINT, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ query }),
    });

    const data = await response.json().catch(() => ({}));

    thinkingEntry.remove();

    if (!response.ok) {
      const message = data.error || `Request failed (HTTP ${response.status})`;
      addEntry("copilot", message, "error");
      statusDot.classList.add("is-error");
      return;
    }

    addToolTrace(data.tool_calls);
    addEntry("copilot", data.answer || "No answer returned.", "assistant");
  } catch (err) {
    thinkingEntry.remove();
    addEntry("copilot", `Couldn't reach the server: ${err.message}`, "error");
    statusDot.classList.add("is-error");
  } finally {
    setBusy(false);
  }
}

composer.addEventListener("submit", (e) => {
  e.preventDefault();
  const query = queryInput.value.trim();
  if (!query) return;
  queryInput.value = "";
  submitQuery(query);
});

renderSuggestions();
renderOrderList();
queryInput.focus();
