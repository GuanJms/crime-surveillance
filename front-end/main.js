const API_BASE_URL = window.__ENV__?.BROKER_SERVICE_URL;

/** -----------------------------
 * Utility helpers
 * -----------------------------*/
function showLoading() { document.getElementById("loading-overlay").classList.remove("hidden"); }
function hideLoading() { document.getElementById("loading-overlay").classList.add("hidden"); }

// Use Leaflet locate for maps
function locateMap(map, onSuccess) {
  if (!map) return;
  showLoading();
  let handled = false;
  function clear() { if(!handled){handled=true; hideLoading();} }
  map.once('locationfound', (e) => {
    clear();
    map.invalidateSize();
    onSuccess(e.latlng);
  });
  map.once('locationerror', (e) => {
    clear();
    displayMessage(e.message || 'Unable to get location', true);
  });
  map.locate({ setView: true, maxZoom: 16, watch: false, enableHighAccuracy: true, timeout: 8000, maximumAge: 0 });
  // Fallback in case neither event fires (older browsers)
  setTimeout(clear, 9000);
}


function displayMessage(message, isError = false) {
  const el = document.getElementById("message");
  el.textContent = message;
  el.style.color = isError ? "#c62828" : "#388e3c";
  if (!message) return;
  // auto-hide after 4s
  setTimeout(() => {
    el.textContent = "";
  }, 4000);
}

// Helper to get stored token
function getToken() {
  return localStorage.getItem("authToken");
}

// Wrapper around fetch that adds Authorization header and common error handling
async function authFetch(path, options = {}) {
  const token = getToken();
  const headers = options.headers || {};
  if (token) headers["Authorization"] = `Bearer ${token}`;
  const resp = await fetch(`${API_BASE_URL}${path}`, {
    ...options,
    headers,
  });

  // Attempt to parse JSON if possible
  let data = undefined;
  try {
    data = await resp.json();
  } catch (_) {
    /* ignore */
  }

  // Handle 401 (invalid / expired token)
  if (resp.status === 401) {
    displayMessage("Session expired. Please login again.", true);
    logout();
    throw new Error("unauthorized");
  }

  // For forbidden we just bubble up
  return { ok: resp.ok, status: resp.status, data };
}

function logout() {
  document.getElementById("api-panel").classList.remove("hidden");
  document.getElementById("landing").classList.remove("full");

  localStorage.removeItem("userRole");
  localStorage.removeItem("authToken");
  // Reset UI
  document.getElementById("dashboard").classList.add("hidden");
  document.getElementById("role-display").textContent = "";
  document.getElementById("auth-container").classList.remove("hidden");
  displayMessage("Logged out.");
}

/** -----------------------------
 * Authentication (login / register)
 * -----------------------------*/
async function handleLogin() {
  const username = document.getElementById("login-username").value.trim();
  const password = document.getElementById("login-password").value.trim();
  if (!username || !password) {
    displayMessage("Please enter username and password", true);
    return;
  }
  const { ok, data } = await authFetch("/users/login", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username, password }),
  });
  if (!ok) {
    displayMessage(data?.error || "Login failed", true);
    return;
  }
  localStorage.setItem("authToken", data.token);
  localStorage.setItem("userRole", data.role);
  displayMessage(`Login successful. Role: ${data.role}`);
  showDashboard();
}

async function handleRegister() {
  const username = document.getElementById("register-username").value.trim();
  const password = document.getElementById("register-password").value.trim();
  if (!username || !password) {
    displayMessage("Please enter username and password", true);
    return;
  }
  const { ok, status, data } = await authFetch("/users/register", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username, password }),
  });
  if (!ok) {
    if (status === 409) {
      displayMessage("User already exists", true);
    } else {
      displayMessage(data?.error || "Registration failed", true);
    }
    return;
  }
  localStorage.setItem("authToken", data.token);
  localStorage.setItem("userRole", data.role);
  displayMessage(`Registration successful – role: ${data.role}`);
  showDashboard();
}

function toggleForms() {
  document.getElementById("login-form").classList.toggle("active");
  document.getElementById("register-form").classList.toggle("active");
}

/** -----------------------------
 * Dashboard navigation
 * -----------------------------*/
function getRole() {
  return localStorage.getItem("userRole") || "CITIZEN";
}

function configureAccessByRole() {
  const role = getRole();
  const allowedMap = {
    CITIZEN: ["report", "crimes", "map"],
    PATROL: ["report", "crimes", "map", "pmap", "patrols", "patrol"],
    DISPATCHER: ["report", "crimes", "map", "pmap", "patrols", "dispatcher"],
    ADMIN: ["report", "crimes", "map", "pmap", "patrols", "patrol", "dispatcher", "admin"],
  };
  const allowed = allowedMap[role] || allowedMap.CITIZEN;

  // Update nav buttons appearance
  document.querySelectorAll("#navbar .nav-btn").forEach((btn) => {
    const sec = btn.dataset.section;
    if (!allowed.includes(sec)) {
      btn.classList.add("locked");
    } else {
      btn.classList.remove("locked");
    }
  });

  // Disable patrol own status/location for ADMIN
  const disableForAdmin = role === "ADMIN";
  document.getElementById("update-status-btn").disabled = disableForAdmin;
  document.getElementById("update-location-btn").disabled = disableForAdmin;

  // role display
  document.getElementById("role-display").textContent = `Role: ${role}`;
}

function isSectionAllowed(section) {
  const role = getRole();
  const acl = {
    CITIZEN: ["report", "crimes", "map"],
    PATROL: ["report", "crimes", "map", "pmap", "patrols", "patrol"],
    DISPATCHER: ["report", "crimes", "map", "pmap", "patrols", "dispatcher"],
    ADMIN: ["report", "crimes", "map", "pmap", "patrols", "patrol", "dispatcher", "admin"],
  };
  return (acl[role] || acl.CITIZEN).includes(section);
}

function showDashboard() {
  document.getElementById("api-panel").classList.add("hidden");
  document.getElementById("landing").classList.add("full");

  document.getElementById("auth-container").classList.add("hidden");
  document.getElementById("dashboard").classList.remove("hidden");
  configureAccessByRole();
  // default view
  showSection("report");
}

function showSection(name) {
  // hide all sections
  document.querySelectorAll(".section").forEach((sec) => {
    sec.classList.remove("active");
  });
  // deactivate nav buttons
  document.querySelectorAll("#navbar .nav-btn").forEach((btn) => {
    btn.classList.remove("active");
  });
  // show selected
  document.getElementById(`section-${name}`).classList.add("active");
  document.querySelector(`#navbar [data-section='${name}']`).classList.add("active");

  if (name === "report") {
    if (reportMap) reportMap.invalidateSize();
    if (!reportMap) initReportMap();
  }

  if (name === "pmap") {
    if (patrolMap) patrolMap.invalidateSize();
    if (!patrolMap) {
      initPatrolMap();
    } else {
      updatePatrolMap();
    }
  }

  if (name === "map") {
    if (crimeMap) crimeMap.invalidateSize();
    if (!crimeMap) {
      initCrimeMap();
    } else {
      updateCrimeMap();
    }
  }
}

/** -----------------------------
 * Crime & Patrol operations
 * -----------------------------*/
async function reportCrime() {
  const desc = document.getElementById("crime-description").value.trim();
  const street = document.getElementById("crime-street").value.trim();
  const city = document.getElementById("crime-city").value.trim();
  const state = document.getElementById("crime-state").value.trim();
  const latitude = parseFloat(document.getElementById("crime-latitude").value);
  const longitude = parseFloat(document.getElementById("crime-longitude").value);
  if (!desc || !street || !city || !state || isNaN(latitude) || isNaN(longitude)) {
    displayMessage("Please fill all fields", true);
    return;
  }
  const body = {
    location: { street, city, state, latitude, longitude },
    description: desc,
  };
  const { ok, data } = await authFetch("/crimes", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  if (ok) {
    displayMessage("Crime reported successfully");
    // clear inputs
    ["crime-description", "crime-street", "crime-city", "crime-state", "crime-latitude", "crime-longitude"].forEach((id) => {
      document.getElementById(id).value = "";
    });
  } else {
    displayMessage(data?.error || "Failed to report crime", true);
  }
}

async function fetchCrimes() {
  const { ok, data } = await authFetch("/crimes");
  const tbody = document.querySelector("#crimes-table tbody");
  tbody.innerHTML = "";
  if (ok && Array.isArray(data)) {
    data.forEach((c) => {
      const tr = document.createElement("tr");
      tr.innerHTML = `
        <td>${c.id ?? ""}</td>
        <td>${c.description ?? ""}</td>
        <td>${c.status ?? ""}</td>
        <td>${c.location?.street ?? ""}</td>
        <td>${c.location?.city ?? ""}</td>
        <td>${c.reported_at ? new Date(c.reported_at).toLocaleString() : ""}</td>`;
      tbody.appendChild(tr);
    });
  } else if (!ok && data?.error) {
    displayMessage(data.error, true);
  }
}

async function fetchPatrols() {
  const { ok, data } = await authFetch("/patrols/register");
    const tbody = document.querySelector("#patrols-table tbody");
  tbody.innerHTML = "";
  if (ok && Array.isArray(data)) {
    data.forEach((p) => {
      const tr = document.createElement("tr");
      tr.innerHTML = `
        <td>${p.userID ?? ""}</td>
        <td>${p.officerID ?? ""}</td>
        <td>${p.officerName ?? ""}</td>
        <td>${p.status ?? ""}</td>
        <td>${p.location?.street ?? ""}</td>
        <td>${p.location?.city ?? ""}</td>`;
      tbody.appendChild(tr);
    });
  } else if (!ok && data?.error) {
    displayMessage(data.error, true);
  }
}

async function updateCrime() {
  const id = document.getElementById("patrol-crime-id").value.trim();
  const status = document.getElementById("patrol-crime-status").value;
  if (!id) {
    displayMessage("Crime ID required", true);
    return;
  }
  const { ok, data } = await authFetch(`/crimes/${id}`, {
    method: "PATCH",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ status }),
  });
  if (ok) {
    displayMessage("Crime updated");
  } else {
    displayMessage(data?.error || "Failed to update", true);
  }
}

async function deleteCrime() {
  const id = document.getElementById("patrol-crime-id").value.trim();
  if (!id) {
    displayMessage("Crime ID required", true);
    return;
  }
  const { ok, data } = await authFetch(`/crimes/${id}`, { method: "DELETE" });
  if (ok) {
    displayMessage("Crime deleted");
  } else {
    displayMessage(data?.error || "Failed to delete", true);
  }
}

async function updatePatrolStatus() {
  const status = document.getElementById("patrol-status-select").value;
  const { ok, data } = await authFetch("/patrols/status", {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ status }),
  });
  if (ok) {
    displayMessage("Status updated");
  } else {
    displayMessage(data?.error || "Failed to update status", true);
  }
}

function initPatrolUpdMap() {
  if (patrolUpdMap) return;
    patrolUpdMap = L.map("patrol-update-map").setView([40.7128, -74.0060], 12);
  L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", { attribution: "© OpenStreetMap contributors" }).addTo(patrolUpdMap);
  if (navigator.geolocation) {
    navigator.geolocation.getCurrentPosition(
      (pos) => {
        const { latitude, longitude } = pos.coords;
        patrolUpdMap.setView([latitude, longitude], 14);
        setPatrolUpdMarker(latitude, longitude);
      },
      () => {
        patrolUpdMap.setView([40.7128, -74.0060], 12);
      }
    );
  } else {
    patrolUpdMap.setView([40.7128, -74.0060], 12);
  }
  patrolUpdMap.on("click", (e) => {
    const { lat, lng } = e.latlng;
    setPatrolUpdMarker(lat, lng);
  });
}

function setPatrolUpdMarker(lat, lng) {
  if (!patrolUpdMap) return;
  if (patrolUpdMarker) {
    patrolUpdMarker.setLatLng([lat, lng]);
  } else {
    patrolUpdMarker = L.marker([lat, lng], { draggable: true }).addTo(patrolUpdMap);
    patrolUpdMarker.on("dragend", () => {
      const pos = patrolUpdMarker.getLatLng();
      syncPatrolLatLngInputs(pos.lat, pos.lng);
    });
  }
  syncPatrolLatLngInputs(lat, lng);
  reverseGeocodePatrol(lat, lng);
}

function syncPatrolLatLngInputs(lat, lng) {
  document.getElementById("patrol-lat").value = lat.toFixed(6);
  document.getElementById("patrol-lng").value = lng.toFixed(6);
}

async function reverseGeocodePatrol(lat, lng) {
  try {
    const url = `https://nominatim.openstreetmap.org/reverse?format=jsonv2&lat=${lat}&lon=${lng}`;
    const res = await fetch(url, { headers: { "User-Agent": "crime-app" } });
    const data = await res.json();
    if (data.address) {
      if (data.address.road) document.getElementById("patrol-street").value = data.address.road;
      if (data.address.city || data.address.town)
        document.getElementById("patrol-city").value = data.address.city || data.address.town;
      if (data.address.state) document.getElementById("patrol-state").value = data.address.state;
    }
  } catch {}
}

async function updatePatrolLocation() {
  const latitude = parseFloat(document.getElementById("patrol-lat").value);
  const longitude = parseFloat(document.getElementById("patrol-lng").value);
  if (isNaN(latitude) || isNaN(longitude)) {
    displayMessage("Latitude & Longitude required", true);
    return;
  }
  const street = document.getElementById("patrol-street").value.trim();
  const city = document.getElementById("patrol-city").value.trim();
  const state = document.getElementById("patrol-state").value.trim();
  const locObj = { latitude, longitude };
  if (street) locObj.street = street;
  if (city) locObj.city = city;
  if (state) locObj.state = state;

  const { ok, data } = await authFetch("/patrols/location", {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ location: locObj }),
  });
  if (ok) {
    displayMessage("Location updated");
  } else {
    displayMessage(data?.error || "Failed to update location", true);
  }
}

async function dispatcherRefresh() {
  // Fetch both crimes and patrols in parallel
    const [crimesRes, patrolsRes] = await Promise.all([
    authFetch("/crimes"),
    authFetch("/patrols/register")
  ]);
  // populate crimes table
  const ctbody = document.querySelector("#dispatcher-crimes-table tbody");
  ctbody.innerHTML = "";
  if (crimesRes.ok && Array.isArray(crimesRes.data)) {
    crimesRes.data.forEach(c=>{
      const tr=document.createElement('tr');
      tr.innerHTML=`<td>${c.id}</td><td>${c.description}</td><td>${c.status}</td>`;
      ctbody.appendChild(tr);
    });
  }
  // populate patrols table
  const ptbody = document.querySelector("#dispatcher-patrols-table tbody");
  ptbody.innerHTML = "";
  if (patrolsRes.ok && Array.isArray(patrolsRes.data)) {
    patrolsRes.data.forEach(p=>{
      const tr=document.createElement('tr');
      tr.innerHTML=`<td>${p.userID}</td><td>${p.officerID}</td><td>${p.officerName}</td><td>${p.status}</td>`;
      ptbody.appendChild(tr);
    });
  }
}

async function dispatchPatrol() {
  const crimeId = document.getElementById("dispatch-crime-id").value.trim();
  const patrolId = document.getElementById("dispatch-patrol-id").value.trim();
  if (!crimeId || !patrolId) {
    displayMessage("Crime ID and Patrol ID required", true);
    return;
  }
  const { ok, data } = await authFetch("/patrols/dispatch", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ crimeId, patrolId }),
  });
  if (ok) {
    displayMessage("Dispatch successful");
  } else {
    displayMessage(data?.error || "Dispatch failed", true);
  }
}

/** Map operations */
let crimeMap;
let mapMarkers;
let crimeData = [];
let crimeMapInterval;

// Patrol map variables
let patrolMap;
let patrolMarkers;
let patrolData = [];
let patrolMapInterval;

// Report map interaction
let reportMap;
let reportMarker;

// Patrol update map
let patrolUpdMap;
let patrolUpdMarker;
function populateCityDropdown() {
  const select = document.getElementById("city-filter");
  const current = select.value;
  // build unique city list
  const cities = Array.from(new Set(crimeData.map((c) => c.location?.city).filter(Boolean)));
  // reset options
  select.innerHTML = '<option value="">All Cities</option>' + cities.map((c)=>`<option value="${c}">${c}</option>`).join("");
  // restore selection if still present
  if (cities.includes(current)) select.value = current;
}

function renderCrimeMarkers() {
  if (!crimeMap) return;
  mapMarkers.clearLayers();
  const filterCity = document.getElementById("city-filter").value;
  let bounds = [];
  crimeData.forEach((c) => {
    if (!c.location) return;
    if (filterCity && c.location.city !== filterCity) return;
    const { latitude, longitude } = c.location;
    if (typeof latitude === "number" && typeof longitude === "number") {
      const marker = L.marker([latitude, longitude]).addTo(mapMarkers);
      marker.bindPopup(`<strong>${c.status}</strong><br/>${c.description}`);
      bounds.push([latitude, longitude]);
    }
  });
  if (bounds.length) {
    crimeMap.fitBounds(bounds, { padding: [50, 50] });
  } else if (!crimeMap._loaded) {
    // set default center if none yet
    crimeMap.setView([40.7128, -74.0060], 12);
  }
}

function initCrimeMap() {
  if (crimeMap) return;
  crimeMap = L.map("crime-map").setView([39.784, -89.652], 13);
  L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
    attribution: "© OpenStreetMap contributors",
  }).addTo(crimeMap);
  mapMarkers = L.layerGroup().addTo(crimeMap);
  updateCrimeMap();
  crimeMapInterval = setInterval(updateCrimeMap, 30000);
  document.getElementById("city-filter").addEventListener("change", renderCrimeMarkers);
}

async function updateCrimeMap() {
  if (!crimeMap) return;
    const { ok, data } = await authFetch("/crimes");
  if (ok && Array.isArray(data)) {
    crimeData = data;
    populateCityDropdown();
    renderCrimeMarkers();
  }
}

/* ------------ Patrols Map ------------*/
function populatePatrolCityDropdown() {
  const select = document.getElementById("pmap-city-filter");
  const current = select.value;
  const cities = Array.from(new Set(patrolData.map((p)=>p.location?.city).filter(Boolean)));
  select.innerHTML = '<option value="">All Cities</option>' + cities.map(c=>`<option value="${c}">${c}</option>`).join("");
  if (cities.includes(current)) select.value = current;
}

function renderPatrolMarkers() {
  if (!patrolMap) return;
  patrolMarkers.clearLayers();
  const filterCity = document.getElementById("pmap-city-filter").value;
  let bounds = [];
  patrolData.forEach((p)=>{
    if (!p.location) return;
    const { latitude, longitude, city, street } = p.location;
    if (filterCity && city !== filterCity) return;
    if (typeof latitude === "number" && typeof longitude === "number") {
      const m = L.marker([latitude, longitude]).addTo(patrolMarkers);
      m.bindPopup(`<strong>${p.officerName}</strong><br/>Status: ${p.status}<br/>${street || ''}`);
      bounds.push([latitude, longitude]);
    }
  });
  if (bounds.length) {
    patrolMap.fitBounds(bounds,{padding:[50,50]});
  } else if (!patrolMap._loaded) {
    patrolMap.setView([40.7128, -74.0060], 12);
  }
}

function initPatrolMap(){
  if(patrolMap) return;
  patrolMap = L.map("patrol-map").setView([40.7128,-74.0060],12);
  L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png",{attribution:"© OpenStreetMap contributors"}).addTo(patrolMap);
  patrolMarkers = L.layerGroup().addTo(patrolMap);
  updatePatrolMap();
  patrolMapInterval = setInterval(updatePatrolMap,30000);
  document.getElementById("pmap-city-filter").addEventListener("change", renderPatrolMarkers);
}

async function updatePatrolMap(){
  if(!patrolMap) return;
  const {ok,data}= await authFetch("/patrols/register");
  if(ok && Array.isArray(data)){
    patrolData= data;
    populatePatrolCityDropdown();
    renderPatrolMarkers();
  }
}

/* ------------ Report Map ------------*/
function initReportMap() {
  if (reportMap) return;
    reportMap = L.map("report-map").setView([40.7128, -74.0060], 12);
  L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
    attribution: "© OpenStreetMap contributors",
  }).addTo(reportMap);
  // Try browser geolocation
  if (navigator.geolocation) {
    navigator.geolocation.getCurrentPosition(
      (pos) => {
        const { latitude, longitude } = pos.coords;
        reportMap.setView([latitude, longitude], 14);
        setReportMarker(latitude, longitude);
      },
      () => {
        reportMap.setView([40.7128, -74.0060], 12); // fallback NYC
      }
    );
  } else {
    reportMap.setView([40.7128, -74.0060], 12);
  }

  reportMap.on("click", (e) => {
    const { lat, lng } = e.latlng;
    setReportMarker(lat, lng);
  });
}

function setReportMarker(lat, lng) {
  if (!reportMap) return;
  if (reportMarker) {
    reportMarker.setLatLng([lat, lng]);
  } else {
    reportMarker = L.marker([lat, lng], { draggable: true }).addTo(reportMap);
    reportMarker.on("dragend", () => {
      const pos = reportMarker.getLatLng();
      syncLatLngInputs(pos.lat, pos.lng);
    });
  }
  syncLatLngInputs(lat, lng);
  reverseGeocodeAndFill(lat, lng);
}

async function reverseGeocodeAndFill(lat, lng) {
  try {
    const url = `https://nominatim.openstreetmap.org/reverse?format=jsonv2&lat=${lat}&lon=${lng}`;
    const res = await fetch(url, { headers: { "User-Agent": "crime-app" } });
    const data = await res.json();
    if (data.address) {
      if (data.address.road) document.getElementById("crime-street").value = data.address.road;
      if (data.address.city || data.address.town) document.getElementById("crime-city").value = data.address.city || data.address.town;
      if (data.address.state) document.getElementById("crime-state").value = data.address.state;
    }
  } catch (e) {
    // ignore network errors
  }
}

function syncLatLngInputs(lat, lng) {
  document.getElementById("crime-latitude").value = lat.toFixed(6);
  document.getElementById("crime-longitude").value = lng.toFixed(6);
}

/** Admin functions */
async function adminChangeRole() {
  const userId = document.getElementById("admin-user-id").value.trim();
  const role = document.getElementById("admin-role-select").value;
  if (!userId) {
    displayMessage("User ID required", true);
    return;
  }
  const { ok, data } = await authFetch(`/admin/users/${userId}/role`, {
    method: "PATCH",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ role }),
  });
  if (ok) {
    displayMessage("Role updated");
  } else {
    displayMessage(data?.error || "Failed to update role", true);
  }
}

async function adminRefreshUsers(){
  const {ok,data}= await authFetch("/admin/users");
  const tbody=document.querySelector("#admin-users-table tbody");
  tbody.innerHTML="";
  if(ok && Array.isArray(data)){
    data.forEach(u=>{
      const tr=document.createElement('tr');
      tr.innerHTML=`<td>${u.id}</td><td>${u.username}</td><td>${u.role}</td><td>${new Date(u.createdAt).toLocaleString()}</td><td>${new Date(u.updatedAt).toLocaleString()}</td><td>${u.lastLogin?new Date(u.lastLogin).toLocaleString():''}</td><td>${u.lastActivity?new Date(u.lastActivity).toLocaleString():''}</td>`;
      tbody.appendChild(tr);
    });
  } else if(!ok){
    displayMessage(data?.error||"Failed to load users",true);
  }
}

async function adminRegisterPatrol() {
  const userId = document.getElementById("patrol-user-id").value.trim();
  const officerId = document.getElementById("patrol-officer-id").value.trim();
  const officerName = document.getElementById("patrol-officer-name").value.trim();
  const status = document.getElementById("patrol-register-status").value;
  if (!userId || !officerId || !officerName) {
    displayMessage("All fields required", true);
    return;
  }
  const body = { userId, officerId, officerName, status };
  const { ok, data } = await authFetch("/patrols/register", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  if (ok) {
    displayMessage("Patrol registered");
  } else {
    displayMessage(data?.error || "Failed to register patrol", true);
  }
}

/** -----------------------------
 * Initialise
 * -----------------------------*/

document.addEventListener("DOMContentLoaded", () => {
  // Auth events
  document.getElementById("login-btn").addEventListener("click", handleLogin);
  document.getElementById("register-btn").addEventListener("click", handleRegister);
  document.getElementById("show-register").addEventListener("click", (e) => {
    e.preventDefault();
    toggleForms();
  });
  document.getElementById("show-login").addEventListener("click", (e) => {
    e.preventDefault();
    toggleForms();
  });

  // If token already stored, skip to dashboard
  if (getToken()) {
    showDashboard();
  }

  // Navbar navigation
  document.querySelectorAll("#navbar .nav-btn").forEach((btn) => {
    btn.addEventListener("click", () => {
      const sec = btn.dataset.section;
      if (!isSectionAllowed(sec)) {
        displayMessage("No access to this section", true);
        return;
      }
      showSection(sec);
    });
  });
  document.getElementById("logout-btn").addEventListener("click", logout);

  // Buttons
  document
    .getElementById("report-crime-btn")
    .addEventListener("click", reportCrime);
  document
    .getElementById("refresh-crimes-btn")
    .addEventListener("click", fetchCrimes);
  document
    .getElementById("refresh-patrols-btn")
    .addEventListener("click", fetchPatrols);
  document
    .getElementById("update-crime-btn")
    .addEventListener("click", updateCrime);
  document
    .getElementById("delete-crime-btn")
    .addEventListener("click", deleteCrime);
  document
    .getElementById("update-status-btn")
    .addEventListener("click", updatePatrolStatus);
  document
    .getElementById("update-location-btn")
    .addEventListener("click", updatePatrolLocation);
  initPatrolUpdMap();

  // location buttons
  document.getElementById("report-loc-btn").addEventListener("click", () => {
    if (navigator.geolocation) {
      locateMap(reportMap, ({lat, lng}) => {
        setReportMarker(lat, lng);
      });
    }
  });
  document.getElementById("crime-map-loc-btn").addEventListener("click", () => {
    if (crimeMap) {
      locateMap(crimeMap, ({lat, lng}) => {
        crimeMap.setView([lat, lng], 14);
      });
    }
  });
  document.getElementById("pmap-loc-btn").addEventListener("click", () => {
    if (patrolMap) {
      locateMap(patrolMap, ({lat, lng}) => {
        patrolMap.setView([lat, lng], 14);
      });
    }
  });
  document.getElementById("patrol-upd-loc-btn").addEventListener("click", () => {
    if (patrolUpdMap) {
      locateMap(patrolUpdMap, ({lat, lng}) => {
        setPatrolUpdMarker(lat, lng);
        patrolUpdMap.setView([lat, lng], 14);
      });
    }
  });
  document
    .getElementById("dispatcher-refresh-btn")
    .addEventListener("click", dispatcherRefresh);
  document
    .getElementById("dispatch-btn")
    .addEventListener("click", dispatchPatrol);
  document
    .getElementById("change-role-btn")
    .addEventListener("click", adminChangeRole);
  document
    .getElementById("register-patrol-btn")
    .addEventListener("click", adminRegisterPatrol);
document.getElementById("admin-refresh-users").addEventListener("click",adminRefreshUsers);

//
 

});
