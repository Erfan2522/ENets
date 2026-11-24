"use strict";

// The alpha UI is intentionally a non-operational preview.
// Do not connect these fields to privileged Android commands until the
// capability adapter, validation, authentication, and audit-safe handling exist.
const form = document.getElementById("hotspot-form");
const message = document.getElementById("form-message");

form.addEventListener("submit", (event) => {
  event.preventDefault();
  message.textContent = "Not applied: Android hotspot integration is not implemented in this alpha.";
});
