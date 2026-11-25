"use strict";

const form = document.getElementById("hotspot-form");
const message = document.getElementById("form-message");
const results = document.getElementById("capability-results");
const clientList = document.getElementById("client-list");
const startButton = document.getElementById("start-button");
const stopButton = document.getElementById("stop-button");
const securityInput = document.getElementById("security");
const autoStartInput = document.getElementById("auto-start");
const profileSelect = document.getElementById("profile-select");
let backendAvailable = false;
let operationBusy = false;
let saveTimer = null;
let savedSettingsReady = false;

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
  if (response.status === 204) return null;
  return response.json();
}

function setMessage(text, isError = false) {
  message.textContent = text;
  message.classList.toggle("error", isError);
}

function formSettings() {
  const security = securityInput.value;
  return {
    ssid: document.getElementById("ssid").value,
    password: security === "open" ? "" : document.getElementById("password").value,
    security,
    band: document.getElementById("band").value,
    channel: document.getElementById("channel").value || "auto",
    width: document.getElementById("width").value,
    hidden: document.getElementById("hidden").checked
  };
}

function populateSettings(settings) {
  if (!settings) return;
  const s = settings.hotspot || {};
  if (s.ssid) document.getElementById("ssid").value = s.ssid;
  if (typeof s.password === "string") document.getElementById("password").value = s.password;
  if (s.security) securityInput.value = s.security;
  if (s.band) document.getElementById("band").value = s.band === "2" ? "2.4" : s.band;
  if (s.channel) document.getElementById("channel").value = s.channel === "auto" ? "" : s.channel;
  if (s.width) document.getElementById("width").value = s.width;
  document.getElementById("hidden").checked = Boolean(s.hidden);
  autoStartInput.checked = Boolean(settings.auto_start);
  updateSecurityFields();
}

function updateSecurityFields() {
  const open = securityInput.value === "open";
  const password = document.getElementById("password");
  password.required = !open;
  password.disabled = open;
  password.placeholder = open ? "Not used for open networks" : "Enter a password";
}

async function saveSettingsNow() {
  if (!backendAvailable || !savedSettingsReady) return;
  const settings = formSettings();
  if (!settings.ssid || (settings.security !== "open" && settings.password.length < 8)) return;
  try {
    await requestJSON("/api/v1/settings", {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ hotspot: settings, auto_start: autoStartInput.checked })
    });
    setMessage("Settings saved securely on this device.");
  } catch (error) {
    setMessage(`Could not save settings: ${error.message}`, true);
  }
}

function scheduleSave() {
  if (!savedSettingsReady) return;
  clearTimeout(saveTimer);
  saveTimer = setTimeout(saveSettingsNow, 450);
}

async function loadSavedSettings() {
  try {
    const saved = await requestJSON("/api/v1/settings");
    if (saved && saved.hotspot) populateSettings(saved);
    savedSettingsReady = true;
    if (!saved || !saved.hotspot || !saved.hotspot.ssid) scheduleSave();
  } catch (error) {
    savedSettingsReady = true;
    setMessage(`Saved settings unavailable: ${error.message}`, true);
  }
}

function updateStatus(status) {
  backendAvailable = Boolean(status.backend_available);
  const connection = document.getElementById("connection-status");
  connection.replaceChildren();
  const dot = document.createElement("i");
  connection.appendChild(dot);
  connection.appendChild(document.createTextNode(backendAvailable ? " Backend ready" : " Backend missing"));
  document.getElementById("hotspot-state").textContent = status.running ? "Running" : "Stopped";
  document.getElementById("hotspot-detail").textContent = status.running ? (status.ssid || "Access point active") : (backendAvailable ? "Ready to configure" : "Install the ENets backend package");
  document.getElementById("client-count").textContent = status.running ? String(status.clients || 0) : "0";
  document.getElementById("active-band").textContent = status.running ? (status.band === "2" ? "2.4 GHz" : status.band === "5" ? "5 GHz" : "Unknown") : "—";
  document.getElementById("hotspot-tag").textContent = status.running ? "Running" : "Not running";
  startButton.disabled = !backendAvailable;
  stopButton.disabled = !backendAvailable || !status.running;
  if (status.running && !savedSettingsReady) {
    if (status.band === "2" || status.band === "5") document.getElementById("band").value = status.band === "2" ? "2.4" : "5";
    if (status.channel) document.getElementById("channel").value = status.channel;
    if (status.security) securityInput.value = status.security;
    updateSecurityFields();
  }
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
    if (!clients.length) { showEmptyClients(); return; }
    clientList.className = "client-table-wrap";
    const table = document.createElement("table");
    const head = document.createElement("tr");
    for (const label of ["MAC address", "IP address", "Device name"]) addText(head, "th", label);
    const thead = document.createElement("thead");
    thead.appendChild(head); table.appendChild(thead);
    const body = document.createElement("tbody");
    for (const client of clients) {
      const row = document.createElement("tr");
      addText(row, "td", client.mac || "—");
      addText(row, "td", client.ip || "—");
      addText(row, "td", client.name && client.name !== "-" ? client.name : "Unknown");
      body.appendChild(row);
    }
    table.appendChild(body); clientList.appendChild(table);
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
    results.replaceChildren();
    const platform = report.platform || {};
    addText(results, "p", `${platform.manufacturer || "Unknown"} ${platform.model || ""} · ${platform.device || "unknown"} · Android ${platform.android_release || "unknown"} (API ${platform.android_api || "unknown"}) · ${platform.primary_abi || "unknown"}`, "muted");
    addText(results, "p", `Wi-Fi interfaces: ${(report.wifi_interfaces || []).join(", ") || "not detected"}`, "muted");
    const list = document.createElement("ul"); list.className = "capability-list";
    for (const [name, item] of Object.entries(report.capabilities || {})) {
      const row = document.createElement("li");
      addText(row, "strong", name.replace(/_/g, " "));
      addText(row, "span", item.status || "unknown", "capability-status");
      addText(row, "small", item.details || "No details available.");
      list.appendChild(row);
    }
    results.appendChild(list);
    for (const warning of report.warnings || []) addText(results, "p", warning, "muted");
    const wpa3 = (report.capabilities || {}).wpa3_softap;
    const option = securityInput.querySelector('option[value="wpa3"]');
    const transition = securityInput.querySelector('option[value="wpa2wpa3"]');
    if (option) option.disabled = !wpa3 || wpa3.status !== "supported";
    if (transition) transition.disabled = !wpa3 || wpa3.status !== "supported";
    if (option && option.disabled && ["wpa3", "wpa2wpa3"].includes(securityInput.value)) securityInput.value = "wpa2";
    updateSecurityFields();
  } catch (error) {
    results.replaceChildren();
    addText(results, "p", `Detection unavailable: ${error.message}`, "muted");
  }
}

async function loadProfiles() {
  try {
    const profiles = await requestJSON("/api/v1/profiles");
    profileSelect.replaceChildren();
    addText(profileSelect, "option", "Select a profile");
    profileSelect.options[0].value = "";
    for (const profile of profiles) {
      const option = addText(profileSelect, "option", profile.name);
      option.value = profile.id;
      option.dataset.profile = JSON.stringify(profile);
    }
  } catch (error) {
    setMessage(`Profiles unavailable: ${error.message}`, true);
  }
}

async function saveProfile() {
  const name = window.prompt("Name this hotspot profile (up to 32 characters):");
  if (!name) return;
  try {
    await requestJSON("/api/v1/profiles", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ name, hotspot: formSettings() })
    });
    await loadProfiles();
    setMessage("Profile saved.");
  } catch (error) {
    setMessage(`Could not save profile: ${error.message}`, true);
  }
}

profileSelect.addEventListener("change", () => {
  const option = profileSelect.selectedOptions[0];
  if (!option || !option.dataset.profile) return;
  const profile = JSON.parse(option.dataset.profile);
  populateSettings({ hotspot: profile.hotspot, auto_start: autoStartInput.checked });
  scheduleSave();
});

async function loadDiagnostics() {
  const target = document.getElementById("diagnostic-summary");
  try {
    const report = await requestJSON("/api/v1/diagnostics");
    target.replaceChildren();
    for (const [key, value] of Object.entries(report)) {
      if (key === "warnings") continue;
      addText(target, "p", `${key.replace(/_/g, " ")}: ${String(value)}`, "muted");
    }
    for (const warning of report.warnings || []) addText(target, "p", warning, "muted");
  } catch (error) {
    target.textContent = `Diagnostics unavailable: ${error.message}`;
  }
}

async function exportDiagnostics() {
  try {
    const response = await fetch("/api/v1/diagnostics/logs", { cache: "no-store" });
    if (!response.ok) throw new Error(await response.text());
    const logs = await response.text();
    const blob = new Blob([logs], { type: "text/plain" });
    const url = URL.createObjectURL(blob);
    const link = document.createElement("a");
    link.href = url;
    link.download = "enets-diagnostics.txt";
    link.click();
    URL.revokeObjectURL(url);
    setMessage("Redacted diagnostic log exported.");
  } catch (error) {
    setMessage(`Diagnostic export failed: ${error.message}`, true);
  }
}

async function loadVersion() {
  try {
    const info = await requestJSON("/api/v1/version");
    document.getElementById("version-label").textContent = info.version;
    document.getElementById("release-link").href = info.release_url;
  } catch (_) {
    document.getElementById("version-label").textContent = "Version unavailable";
  }
}

form.addEventListener("submit", async (event) => {
  event.preventDefault();
  if (!backendAvailable) return;
  const payload = formSettings();
  operationBusy = true;
  startButton.disabled = true; stopButton.disabled = true;
  setMessage("Applying settings and starting hotspot…");
  try {
    await saveSettingsNow();
    const status = await requestJSON("/api/v1/hotspot/start", {
      method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify(payload)
    });
    updateStatus(status);
    setMessage(status.running ? "Hotspot started. Verify client connectivity." : "The engine returned, but the hotspot is not running.", !status.running);
    await loadClients();
  } catch (error) {
    setMessage(error.message, true);
  } finally {
    operationBusy = false;
    await refreshStatus();
  }
});

stopButton.addEventListener("click", async () => {
  if (!backendAvailable) return;
  operationBusy = true; stopButton.disabled = true; startButton.disabled = true;
  setMessage("Stopping hotspot…");
  try {
    const status = await requestJSON("/api/v1/hotspot/stop", { method: "POST" });
    updateStatus(status); setMessage("Hotspot stopped."); showEmptyClients();
  } catch (error) { setMessage(error.message, true); }
  finally { operationBusy = false; await refreshStatus(); }
});

document.getElementById("save-profile").addEventListener("click", saveProfile);
document.getElementById("export-diagnostics").addEventListener("click", exportDiagnostics);
for (const control of form.querySelectorAll("input, select")) {
  control.addEventListener("input", scheduleSave);
  control.addEventListener("change", () => {
    updateSecurityFields();
    scheduleSave();
  });
}
autoStartInput.addEventListener("change", scheduleSave);
updateSecurityFields();
refreshStatus();
loadSavedSettings();
loadCapabilities();
loadProfiles();
loadDiagnostics();
loadVersion();
setInterval(refreshStatus, 5000);
