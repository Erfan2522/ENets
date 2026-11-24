"use strict";

// Capability detection is read-only. It never applies hotspot settings.
const form = document.getElementById("hotspot-form");
const message = document.getElementById("form-message");
const results = document.getElementById("capability-results");

form.addEventListener("submit", (event) => {
  event.preventDefault();
  message.textContent = "Not applied: Android hotspot configuration is not implemented in this alpha.";
});

function clearChildren(parent) {
  while (parent.firstChild) parent.removeChild(parent.firstChild);
}

function addText(parent, tag, value, className) {
  const node = document.createElement(tag);
  node.textContent = value;
  if (className) node.className = className;
  parent.appendChild(node);
  return node;
}

function displayCapabilities(report) {
  clearChildren(results);
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

  for (const warning of report.warnings || []) {
    addText(results, "p", warning, "muted");
  }
}

async function loadCapabilities() {
  clearChildren(results);
  addText(results, "p", "Checking device capabilities…", "muted");
  try {
    const response = await fetch("/api/v1/capabilities", { headers: { Accept: "application/json" }, cache: "no-store" });
    if (!response.ok) throw new Error(`HTTP ${response.status}`);
    displayCapabilities(await response.json());
  } catch (error) {
    clearChildren(results);
    addText(results, "p", `Detection unavailable: ${error.message}`, "muted");
    addText(results, "p", "Open this dashboard through the ENets local server to query the Android device.", "muted");
  }
}

loadCapabilities();
