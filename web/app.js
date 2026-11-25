"use strict";

const form = document.getElementById("hotspot-form");
const message = document.getElementById("form-message");
const results = document.getElementById("capability-results");
const clientList = document.getElementById("client-list");
const startButton = document.getElementById("start-button");
const stopButton = document.getElementById("stop-button");
const securityInput = document.getElementById("security");
let backendAvailable = false;
let operationBusy = false;

function addText(parent, tag, value, className) {
  const node = document.createElement(tag);
  node.textContent = value;
  if (className) node.className = className;
  parent.appendChild(node);
  return node;
}

async function requestJSON(path, options = {}) {
  const response = await fetch(path, {
    cache: "no-store",
    ...options,
    headers: { Accept: "application/json", ...(options.headers || {}) }
  });
  if (!response.ok) {
    const detail = await response.text();
    throw new Error(detail || `HTTP ${response.status}`);
  }
  return response.json();
}

function setMessage(text, isError = false) {
  message.textContent = text;
  message.classList.toggle("error", isError);
}

function updateSecurityFields() {
  const open = securityInput.value === "open";
  const password = document.getElementById("password");
  password.required = !open;
  password.disabled = open;
  if (open) password.value = "";
  password.placeholder = open ? "Not used for open networks" : "Enter a password";
}

function updateStatus(status) {
  backendAvailable = Boolean(status.backend_available);
  const connection = document.getElementById("connection-status");
  connection.innerHTML = "";
  const dot = document.createElement("i");
  connection.appendChild(dot);
  connection.appendChild(document.createTextNode(backendAvailable ? " Backend ready" : " Backend missing"));

  document.getElementById("hotspot-state").textContent = status.running ? "Running" : "Stopped";
  document.getElementById("hotspot-detail").textContent = status.running
    ? (status.ssid || "Access point active")
    : (backendAvailable ? "Ready to configure" : "Install the ENets backend package");
  document.getElementById("client-count").textContent = status.running ? String(status.clients || 0) : "0";
  document.getElementById("active-band").textContent = status.running
    ? (status.band === "2" ? "2.4 GHz" : status.band === "5" ? "5 GHz" : "Unknown")
    : "—";
  document.getElementById("hotspot-tag").textContent = status.running ? "Running" : "Not running";
  startButton.disabled = !backendAvailable;
  stopButton.disabled = !backendAvailable || !status.running;

  if (status.running) {
    if (status.band === "2" || status.band === "5") {
      document.getElementById("band").value = status.band === "2" ? "2.4" : "5";
    }
    if (status.channel) document.getElementById("channel").value = status.channel;
    if (status.security) document.getElementById("security").value = status.security;
    updateSecurityFields();
  }
}

function displayCapabilities(report) {
  results.replaceChildren();
  const platform = report.platform || {};
  addText(results, "p", `${platform.manufacturer || "Unknown"} ${platform.model || ""} · Device: ${platform.device || "unknown"} · Android: ${platform.android_release || "unknown"} (API ${platform.android_api || "unknown"}) · ABI: ${platform.primary_abi || "unknown"}`, "muted");
  addText(results, "p", `Wi-Fi interfaces: ${(report.wifi_interfaces || []).join(", ") || "not detected"}`, "muted");

  const list = document.createElement("ul");
  list.className = "capability-list";
  for (const [name, item] of Object.entries(report.capabilities || {})) {
    const row = document.createElement("li");
    addText(row, "strong", name.replace(/_/g, " "));
    addText(row, "span", item.status || "unknown", "capability-status");
    addText(row, "small", item.details || "No details available.");
    list.appendChild(row);
  }
  results.appendChild(list);
  for (const warning of report.warnings || []) addText(results, "p", warning, "muted");
}

async function refreshStatus() {
  if (operationBusy) return;
  try {
    const status = await requestJSON("/api/v1/hotspot/status");
    updateStatus(status);
    if (status.running) await loadClients();
    else showEmptyClients();
  } catch (error) {
    document.getElementById("connection-status").textContent = " Server unavailable";
    setMessage(error.message, true);
  }
}

function showEmptyClients() {
  clientList.replaceChildren();
  addText(clientList, "strong", "No connected devices");
  addText(clientList, "span", "Start the hotspot to discover connected devices.");
  clientList.className = "empty";
}

async function loadClients() {
  try {
    const clients = await requestJSON("/api/v1/hotspot/clients");
    clientList.replaceChildren();
    if (!clients.length) {
      showEmptyClients();
      return;
    }
    clientList.className = "client-table-wrap";
    const table = document.createElement("table");
    const head = document.createElement("tr");
    for (const label of ["Device", "IP address", "Name"]) addText(head, "th", label);
    const thead = document.createElement("thead");
    thead.appendChild(head);
    table.appendChild(thead);
    const body = document.createElement("tbody");
    for (const client of clients) {
      const row = document.createElement("tr");
      addText(row, "td", client.mac || "—");
      addText(row, "td", client.ip || "—");
      addText(row, "td", client.name && client.name !== "-" ? client.name : "Unknown");
      body.appendChild(row);
    }
    table.appendChild(body);
    clientList.appendChild(table);
  } catch (error) {
    clientList.replaceChildren();
    addText(clientList, "strong", "Client list unavailable");
    addText(clientList, "span", error.message);
    clientList.className = "empty";
  }
}

async function loadCapabilities() {
  try {
    const report = await requestJSON("/api/v1/capabilities");
    displayCapabilities(report);
  } catch (error) {
    results.replaceChildren();
    addText(results, "p", `Detection unavailable: ${error.message}`, "muted");
  }
}

form.addEventListener("submit", async (event) => {
  event.preventDefault();
  if (!backendAvailable) return;
  const security = securityInput.value;
  const payload = {
    ssid: document.getElementById("ssid").value,
    password: security === "open" ? "" : document.getElementById("password").value,
    security,
    band: document.getElementById("band").value,
    channel: document.getElementById("channel").value || "auto",
    width: document.getElementById("width").value,
    hidden: document.getElementById("hidden").checked
  };
  operationBusy = true;
  startButton.disabled = true;
  stopButton.disabled = true;
  setMessage("Applying settings and starting hotspot…");
  try {
    const status = await requestJSON("/api/v1/hotspot/start", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload)
    });
    updateStatus(status);
    document.getElementById("password").value = "";
    setMessage(status.running ? "Hotspot started. Verify client connectivity on this device." : "The engine returned, but the hotspot is not running.", !status.running);
    await loadClients();
  } catch (error) {
    setMessage(error.message, true);
  } finally {
    operationBusy = false;
    startButton.disabled = !backendAvailable;
    await refreshStatus();
  }
});

stopButton.addEventListener("click", async () => {
  if (!backendAvailable) return;
  operationBusy = true;
  stopButton.disabled = true;
  startButton.disabled = true;
  setMessage("Stopping hotspot…");
  try {
    const status = await requestJSON("/api/v1/hotspot/stop", { method: "POST" });
    updateStatus(status);
    setMessage("Hotspot stopped.");
    showEmptyClients();
  } catch (error) {
    setMessage(error.message, true);
  } finally {
    operationBusy = false;
    await refreshStatus();
  }
});

securityInput.addEventListener("change", updateSecurityFields);
updateSecurityFields();
refreshStatus();
loadCapabilities();
setInterval(refreshStatus, 5000);
