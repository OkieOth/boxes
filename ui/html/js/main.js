// Return all box ids present in the current SVG (top-level boxes only, not blacklisted)
// Base path detection: derive deployment base (e.g., "/xxx") at runtime
function __detectBasePath() {
    try {
        // Prefer script src containing "/html/" to locate the app root sibling folders
        const scripts = document.getElementsByTagName("script");
        for (let i = 0; i < scripts.length; i++) {
            const src = scripts[i].src || "";
            if (!src) continue;
            try {
                const u = new URL(src, window.location.origin);
                const p = u.pathname || "";
                const idx = p.indexOf("/html/");
                if (idx >= 0) {
                    const bp = p.substring(0, idx) || "/";
                    return bp.endsWith("/") ? bp.slice(0, -1) : bp;
                }
            } catch {}
        }
    } catch {}
    try {
        const p = window.location.pathname || "/";
        const idx = p.indexOf("/html/");
        if (idx >= 0) {
            const bp = p.substring(0, idx) || "/";
            return bp.endsWith("/") ? bp.slice(0, -1) : bp;
        }
        // Served directly from base (e.g., "/xxx/")
        return p.replace(/\/$/, "");
    } catch {}
    return "";
}

window.getBasePath = function () {
    if (typeof window.basePath === "string") return window.basePath;
    const bp = __detectBasePath();
    window.basePath = bp || "";
    return window.basePath;
};

window.getAllBoxIds = function () {
    const svg = getSvg();
    if (!svg) return [];
    // Assume box elements have id matching the box id pattern (e.g., box_1, box_2, ...)
    // We'll collect all elements with an id that matches the box prefix logic
    const elements = svg.querySelectorAll("[id]");
    const ids = new Set();
    elements.forEach((el) => {
        const boxId = getBoxPrefix(el.id);
        if (boxId) ids.add(boxId);
    });
    return Array.from(ids);
};

async function promptForUploadedMixin() {
    const upload = await new Promise((resolve, reject) => {
        const input = document.createElement("input");
        input.type = "file";
        input.accept = ".yaml,.yml";
        input.style.display = "none";
        document.body.appendChild(input);
        input.addEventListener("change", () => {
            const file = input.files && input.files[0];
            document.body.removeChild(input);
            if (!file) {
                reject(new Error("No file selected"));
                return;
            }
            const reader = new FileReader();
            reader.onload = () =>
                resolve({
                    content: String(reader.result || ""),
                    fileName: file.name || `mixin_${Date.now()}`,
                });
            reader.onerror = (err) => reject(err);
            reader.readAsText(file);
        });
        input.click();
    });
    const uploadedName = upload.fileName || `mixin_${Date.now()}`;
    const title = parseTitleFromYaml(upload.content, uploadedName);
    const entry = upsertUploadedMixin(uploadedName, title, upload.content);
    return {
        id: entry.id,
        label: entry.title || entry.id,
        content: upload.content,
    };
}

async function triggerToolbarComboUpload() {
    try {
        closeToolbarComboDropdown(false);
        const { select } = getToolbarComboElements();
        const upload = await promptForUploadedMixin();
        if (!upload || !select) return;
        ensureUploadOptionsInCombo(select);
        const uploadedVal = `uploaded::${upload.id}`;
        toolbarComboState.contentCache.set(uploadedVal, upload.content);
        toolbarComboState.selectionMeta.set(uploadedVal, {
            ...(toolbarComboState.selectionMeta.get(uploadedVal) || {}),
            label: upload.label,
        });
        if (toolbarComboState.selectedValues.includes(uploadedVal)) {
            refreshToolbarComboUI();
            applySelectedMixins();
            showUploadRefreshToast(
                upload.label
                    ? `Updated: ${upload.label}`
                    : "Uploaded mixin updated"
            );
        } else {
            addToolbarComboSelection(uploadedVal);
        }
    } catch (err) {
        console.error("Upload cancelled or failed:", err);
    }
}

let uploadToastTimer = null;
function showUploadRefreshToast(message) {
    const el = document.getElementById("upload-refresh-toast");
    if (!el) return;
    el.textContent = message || "Uploaded mixin updated";
    el.classList.remove("hidden");
    el.setAttribute("aria-hidden", "false");
    requestAnimationFrame(() => {
        el.classList.add("show");
    });
    if (uploadToastTimer) clearTimeout(uploadToastTimer);
    uploadToastTimer = setTimeout(() => {
        el.classList.remove("show");
        el.setAttribute("aria-hidden", "true");
        setTimeout(() => {
            el.classList.add("hidden");
        }, 180);
    }, 1600);
}

function collectBadgeIds(listId) {
    const list = document.getElementById(listId);
    if (!list) return [];
    return Array.from(list.querySelectorAll(".badge"))
        .map((b) => b.dataset.hid)
        .filter(Boolean);
}

function restoreBadgeCollectorFromIds(savedBadgeIds) {
    const list = document.getElementById("badge-list");
    if (!list || !Array.isArray(savedBadgeIds)) return;
    list.innerHTML = "";
    savedBadgeIds.forEach((hid) => {
        if (!hid) return;
        const svg = document.querySelector("#canvas svg");
        let el = svg ? svg.querySelector(`[id='${hid}']`) : null;
        if (!el && svg) {
            el = svg.querySelector(`[id^='${hid}']`);
        }
        if (el && window.createBadgeForShape) {
            const badge = window.createBadgeForShape(el);
            badge.dataset.hid = hid;
            list.appendChild(badge);
        } else if (window.getCaptionForId) {
            const span = document.createElement("span");
            span.className = "badge";
            span.dataset.hid = hid;
            const label = document.createElement("span");
            label.textContent = window.getCaptionForId(hid);
            span.appendChild(label);
            list.appendChild(span);
        }
    });
    requestAnimationFrame(window.refitAllBadges || (() => {}));
}

async function loadYamlForComboValue(value) {
    if (!value || value === "__upload__") return "";
    if (toolbarComboState.contentCache.has(value)) {
        return toolbarComboState.contentCache.get(value);
    }
    if (value.startsWith("uploaded::")) {
        const id = value.slice("uploaded::".length);
        const list = getUploadedMixins();
        const found = list.find((x) => x && x.id === id);
        if (found) {
            const content = String(found.content || "");
            toolbarComboState.contentCache.set(value, content);
            return content;
        }
        console.warn("Uploaded mixin not found:", id);
        return "";
    }
    try {
        const resp = await fetch(window.getBasePath() + "/data/" + value, {
            cache: "no-cache",
        });
        if (!resp.ok) throw new Error("HTTP " + resp.status);
        const text = await resp.text();
        toolbarComboState.contentCache.set(value, text);
        return text;
    } catch (err) {
        console.error("Failed to load YAML for", value, err);
        return "";
    }
}

let formatMixinContent = "";
let formatMixinLoaded = false;
let formatMixinLoadPromise = null;

async function ensureFormatMixinLoaded() {
    const name =
        typeof window.formatMixin === "string"
            ? window.formatMixin.trim()
            : "";
    if (!name) {
        formatMixinLoaded = true;
        formatMixinContent = "";
        return;
    }
    if (formatMixinLoaded) return;
    if (formatMixinLoadPromise) return formatMixinLoadPromise;
    formatMixinLoadPromise = (async () => {
        try {
            const resp = await fetch(window.getBasePath() + "/data/" + name, {
                cache: "no-cache",
            });
            if (!resp.ok) throw new Error("HTTP " + resp.status);
            formatMixinContent = await resp.text();
        } catch (err) {
            console.error("Failed to load format mixin", name, err);
            formatMixinContent = "";
        } finally {
            formatMixinLoaded = true;
            formatMixinLoadPromise = null;
        }
    })();
    return formatMixinLoadPromise;
}

function getCombinedMixins() {
    const combined = [];
    if (formatMixinContent) combined.push(formatMixinContent);
    if (Array.isArray(mixins) && mixins.length) {
        combined.push(...mixins);
    }
    return combined;
}

function normalizeCreateSvgResult(result, fallbackExpanded, fallbackBlacklisted) {
    let svgStr = "";
    let expanded = fallbackExpanded;
    let blacklisted = fallbackBlacklisted;

    if (result && typeof result === "object") {
        if (typeof result.SVG === "string") {
            svgStr = result.SVG;
        } else if (typeof result.svg === "string") {
            svgStr = result.svg;
        }

        if (Array.isArray(result.Expanded)) {
            expanded = result.Expanded;
        } else if (Array.isArray(result.expanded)) {
            expanded = result.expanded;
        }

        if (Array.isArray(result.Blacklisted)) {
            blacklisted = result.Blacklisted;
        } else if (Array.isArray(result.blacklisted)) {
            blacklisted = result.blacklisted;
        }
    } else if (typeof result === "string") {
        svgStr = result;
    }

    return { svgStr, expanded, blacklisted };
}

function applyExpandedAndBlacklistState(expandedIds, blacklistIds) {
    restoreBadgeCollectorFromIds(expandedIds || []);
    if (Array.isArray(blacklistIds)) {
        blacklist = blacklistIds.slice();
        window.blacklist = blacklist;
    } else {
        blacklist = [];
        window.blacklist = blacklist;
    }
    updateBlacklistUI();
}

async function regenerateSvgWithState(expandedIds, blacklistIds) {
    if (!getActiveCreateSvgFunction()) {
        return;
    }
    const canvas = document.getElementById("canvas");
    if (!canvas) return;
    const arg =
        typeof window.input === "string" && window.input.length > 0
            ? window.input
            : "";
    console.log(
        "Refreshing SVG: ",
        expandedIds,
        "blacklist ids: ",
        blacklistIds,
        "comments hidden: ",
        window.hideCommentsEnabled
    );
    try {
        let result = await callCreateSvgWithMode(
            arg,
            expandedIds,
            blacklistIds
        );
        result = result && typeof result.then === "function" ? await result : result;
        const normalized = normalizeCreateSvgResult(
            result,
            expandedIds,
            blacklistIds
        );
        if (
            typeof normalized.svgStr !== "string" ||
            !normalized.svgStr.trim().startsWith("<svg")
        ) {
            console.error("createSvg did not return a valid SVG string.");
            console.error(result);
            return;
        }
        canvas.innerHTML = normalized.svgStr;
        const evtSwap = new Event("htmx:afterSwap", { bubbles: true });
        canvas.dispatchEvent(evtSwap);
        applyExpandedAndBlacklistState(
            normalized.expanded,
            normalized.blacklisted
        );
    } catch (err) {
        console.error("Error updating SVG via createSvg:", err);
    }
}

async function regenerateSvgWithConnectedOnce(expandedIds, blacklistIds) {
    if (typeof window.createSvgForConnected !== "function") {
        return;
    }
    const canvas = document.getElementById("canvas");
    if (!canvas) return;
    const arg =
        typeof window.input === "string" && window.input.length > 0
            ? window.input
            : "";
    console.log(
        "Refreshing SVG (connected): ",
        expandedIds,
        "blacklist ids: ",
        blacklistIds,
        "comments hidden: ",
        window.hideCommentsEnabled
    );
    try {
        await ensureFormatMixinLoaded();
        let result = window.createSvgForConnected(
            arg,
            getCombinedMixins(),
            window.debug
        );
        result = result && typeof result.then === "function" ? await result : result;
        const normalized = normalizeCreateSvgResult(
            result,
            expandedIds,
            blacklistIds
        );
        if (
            typeof normalized.svgStr !== "string" ||
            !normalized.svgStr.trim().startsWith("<svg")
        ) {
            console.error("createSvg did not return a valid SVG string.");
            console.error(result);
            return;
        }
        canvas.innerHTML = normalized.svgStr;
        const evtSwap = new Event("htmx:afterSwap", { bubbles: true });
        canvas.dispatchEvent(evtSwap);
        applyExpandedAndBlacklistState(
            normalized.expanded,
            normalized.blacklisted
        );
    } catch (err) {
        console.error("Error updating SVG via createSvg:", err);
    }
}

async function applySelectedMixins() {
    const seq = ++toolbarComboState.applyToken;
    const selectedValues = toolbarComboState.selectedValues.slice();
    const contents = [];
    for (const value of selectedValues) {
        if (toolbarComboState.hiddenValues.has(value)) continue;
        const yamlContent = await loadYamlForComboValue(value);
        if (toolbarComboState.applyToken !== seq) {
            return;
        }
        if (yamlContent) {
            contents.push(yamlContent);
        }
    }
    mixins = contents;
    window.currentYamlFile = selectedValues.length ? selectedValues.join(",") : undefined;
    if (!getActiveCreateSvgFunction()) {
        return;
    }
    const savedBadgeState = collectBadgeIds("badge-list");
    const savedBlacklistState = collectBadgeIds("blacklist-list");
    if (toolbarComboState.applyToken !== seq) {
        return;
    }
    if (connectedModeActive) {
        await regenerateSvgWithConnectedOnce(savedBadgeState, savedBlacklistState);
    } else {
        await regenerateSvgWithState(savedBadgeState, savedBlacklistState);
    }
}
let state = { scale: 1, tx: 0, ty: 0 };
let baseSize = { width: 0, height: 0 };
let minimapVisible = false;
// Pan tool state
let panToolActive = false;
let spacePressed = false;
let isDragging = false;
let dragStart = { x: 0, y: 0, tx: 0, ty: 0 };
// NEW: undo stack of previous badge states
let undoStack = [];

// Update toolbar active states
let blacklistMode = false;
let blacklist = [];

// NEW: global additional mixins for persistence
let mixins = [];
let collectorPanelVisible = false;
let collectorVisibilityGuardAttached = false;
let commentLegendPanelVisible = false;
let connectionAnimationEnabled = false;
let connectedModeActive = false;
const commentLegendState = {
    byNodeId: new Map(),
    byGroupClass: new Map(),
    itemActiveCount: new Map(),
};
const SVG_NS = "http://www.w3.org/2000/svg";
window.connectionAnimationEnabled = connectionAnimationEnabled;

// --- Editable toolbar combobox state ---
const toolbarComboState = {
    isOpen: false,
    filterText: "",
    highlightedIndex: -1,
    filteredOptions: [],
    selectedValues: [], // ordered list of selected mixin option values
    hiddenValues: new Set(), // values that are selected but temporarily hidden
    selectionMeta: new Map(), // value -> { label, content? }
    contentCache: new Map(), // value -> yaml content
    applyToken: 0,
};

function getToolbarComboElements() {
    return {
        root: document.getElementById("toolbar-combobox"),
        input: document.getElementById("toolbar-combo-input"),
        toggle: document.getElementById("toolbar-combo-toggle"),
        dropdown: document.getElementById("toolbar-combo-dropdown"),
        select: document.getElementById("toolbar-combo"),
    };
}

function getSelectedCollectorElements() {
    return {
        panel: document.getElementById("selected-collector"),
        list: document.getElementById("selected-mixins-list"),
    };
}

function getToolbarComboOptions() {
    const { select } = getToolbarComboElements();
    if (!select) return [];
    return Array.from(select.options || []).map((opt) => ({
        value: opt.value,
        label: opt.textContent || opt.value || "",
    }));
}

function resetToolbarComboInput() {
    const { input } = getToolbarComboElements();
    if (!input) return;
    input.value = toolbarComboState.filterText || "";
    if (!toolbarComboState.filterText) {
        input.setAttribute("aria-activedescendant", "");
    }
}

function getToolbarComboOptionLabel(value) {
    if (!value) return "";
    const options = getToolbarComboOptions();
    const found = options.find((opt) => opt.value === value);
    return (found && found.label) || value;
}

function isToolbarComboValueSelected(value) {
    return toolbarComboState.selectedValues.includes(value);
}

function createSelectedCollectorBadge(value, label) {
    const badge = document.createElement("div");
    const isHidden = toolbarComboState.hiddenValues.has(value);
    badge.className = "selected-mixin-badge" + (isHidden ? " selected-mixin-badge--hidden" : "");
    badge.dataset.value = value;

    const text = document.createElement("span");
    text.className = "selected-mixin-label";
    text.textContent = label;

    const toggleBtn = document.createElement("button");
    toggleBtn.type = "button";
    toggleBtn.className = "selected-mixin-toggle tool-btn";
    toggleBtn.title = isHidden ? `Show ${label}` : `Hide ${label}`;
    toggleBtn.innerHTML = isHidden
        ? '<i class="fa-solid fa-eye-slash"></i>'
        : '<i class="fa-solid fa-eye"></i>';

    const removeBtn = document.createElement("button");
    removeBtn.type = "button";
    removeBtn.className = "selected-mixin-remove tool-btn";
    removeBtn.title = `Remove ${label}`;
    removeBtn.innerHTML = '<i class="fa-solid fa-xmark"></i>';

    const btnGroup = document.createElement("span");
    btnGroup.className = "selected-mixin-btn-group";
    btnGroup.appendChild(toggleBtn);
    btnGroup.appendChild(removeBtn);

    badge.appendChild(text);
    badge.appendChild(btnGroup);
    return badge;
}

function updateToolbarComboSelectedList() {
    const { panel, list } = getSelectedCollectorElements();
    if (!list) return;
    list.innerHTML = "";
    const vals = toolbarComboState.selectedValues;
    if (!vals.length) {
        const empty = document.createElement("div");
        empty.className = "collector-empty";
        empty.textContent = "No mixins selected.";
        list.appendChild(empty);
    } else {
        vals.forEach((value) => {
            const label =
                toolbarComboState.selectionMeta.get(value)?.label ||
                getToolbarComboOptionLabel(value) ||
                value;
            list.appendChild(createSelectedCollectorBadge(value, label));
        });
    }
    if (panel) {
        panel.setAttribute(
            "aria-hidden",
            panel.classList.contains("hidden") ? "true" : "false"
        );
    }
    positionBlacklistCollector();
}

function initSelectedCollectorInteractions() {
    const { list } = getSelectedCollectorElements();
    if (!list || list.__selectedHandlerAttached) return;
    list.__selectedHandlerAttached = true;
    list.addEventListener("click", (evt) => {
        const removeBtn = evt.target.closest(".selected-mixin-remove");
        if (removeBtn) {
            const badge = removeBtn.closest(".selected-mixin-badge");
            const value = badge ? badge.dataset.value : null;
            if (value) {
                evt.preventDefault();
                evt.stopPropagation();
                removeToolbarComboSelection(value);
            }
            return;
        }
        const toggleBtn = evt.target.closest(".selected-mixin-toggle");
        if (toggleBtn) {
            const badge = toggleBtn.closest(".selected-mixin-badge");
            const value = badge ? badge.dataset.value : null;
            if (value) {
                evt.preventDefault();
                evt.stopPropagation();
                if (toolbarComboState.hiddenValues.has(value)) {
                    toolbarComboState.hiddenValues.delete(value);
                } else {
                    toolbarComboState.hiddenValues.add(value);
                }
                updateToolbarComboSelectedList();
                applySelectedMixins();
            }
        }
    });
}

function renderToolbarComboDropdown() {
    const { dropdown } = getToolbarComboElements();
    if (!dropdown) return;
    const filter = toolbarComboState.filterText.trim().toLowerCase();
    const options = getToolbarComboOptions();
    const filtered = filter
        ? options.filter((opt) => {
              const label = opt.label.toLowerCase();
              const value = String(opt.value || "").toLowerCase();
              return label.includes(filter) || value.includes(filter);
          })
        : options;
    toolbarComboState.filteredOptions = filtered;
    dropdown.innerHTML = "";

    if (!filtered.length) {
        const empty = document.createElement("div");
        empty.className = "toolbar-combo-empty";
        empty.textContent = "No matches";
        dropdown.appendChild(empty);
        toolbarComboState.highlightedIndex = -1;
        return;
    }

    if (
        toolbarComboState.highlightedIndex < 0 ||
        toolbarComboState.highlightedIndex >= filtered.length
    ) {
        toolbarComboState.highlightedIndex = 0;
    }

    filtered.forEach((opt, idx) => {
        const row = document.createElement("div");
        row.className = "toolbar-combo-option";
        row.dataset.value = opt.value;
        row.id = `toolbar-combo-option-${idx}`;
        row.setAttribute("role", "option");
        const displayLabel = opt.label || opt.value || "(unnamed)";
        row.title = displayLabel;
        const isUpload = opt.value === "__upload__";
        if (isUpload) {
            row.classList.add("toolbar-combo-option-upload");
            const icon = document.createElement("i");
            icon.className = "fa-solid fa-upload";
            const label = document.createElement("span");
            label.textContent = displayLabel;
            row.appendChild(icon);
            row.appendChild(label);
            row.addEventListener("mousedown", (evt) => {
                evt.preventDefault();
                triggerToolbarComboUpload();
            });
        } else {
            const checkbox = document.createElement("input");
            checkbox.type = "checkbox";
            checkbox.tabIndex = -1;
            checkbox.checked = isToolbarComboValueSelected(opt.value);
            checkbox.setAttribute("aria-hidden", "true");
            const label = document.createElement("span");
            label.className = "toolbar-combo-option-label";
            label.textContent = displayLabel;
            row.appendChild(checkbox);
            row.appendChild(label);
            row.addEventListener("mousedown", (evt) => {
                evt.preventDefault();
                toggleToolbarComboValue(opt.value);
            });
            row.setAttribute("aria-selected", checkbox.checked ? "true" : "false");
        }
        row.classList.toggle("active", idx === toolbarComboState.highlightedIndex);
        row.addEventListener("mouseenter", () => {
            toolbarComboState.highlightedIndex = idx;
            updateToolbarComboActiveOption();
        });
        dropdown.appendChild(row);
    });

    updateToolbarComboActiveOption();
    if (toolbarComboState.isOpen) {
        positionToolbarComboDropdown();
    }
}

function updateToolbarComboActiveOption() {
    const { dropdown, input } = getToolbarComboElements();
    if (!dropdown || !input) return;
    const rows = dropdown.querySelectorAll(".toolbar-combo-option");
    let activeId = "";
    rows.forEach((row, idx) => {
        const isActive = idx === toolbarComboState.highlightedIndex;
        row.classList.toggle("active", isActive);
        if (isActive) {
            activeId = row.id;
            ensureToolbarComboOptionVisible(row);
        }
    });
    if (activeId) {
        input.setAttribute("aria-activedescendant", activeId);
    } else {
        input.removeAttribute("aria-activedescendant");
    }
}

function ensureToolbarComboOptionVisible(row) {
    if (!row) return;
    try {
        row.scrollIntoView({ block: "nearest" });
    } catch {
        /* ignore */
    }
}

function getToolbarComboAnchorElement() {
    const { root } = getToolbarComboElements();
    if (!root) return null;
    const input = root.querySelector(".toolbar-combo-input");
    if (input) return input;
    return root.querySelector(".toolbar-combo-field") || root;
}

function positionToolbarComboDropdown() {
    if (!toolbarComboState.isOpen) return;
    const { dropdown } = getToolbarComboElements();
    const anchor = getToolbarComboAnchorElement();
    if (!dropdown || !anchor) return;
    const rect = anchor.getBoundingClientRect();
    const viewportWidth = document.documentElement.clientWidth || window.innerWidth || 0;
    const viewportHeight = document.documentElement.clientHeight || window.innerHeight || 0;
    const margin = 8;
    let width = rect.width;
    if (!Number.isFinite(width) || width <= 0) {
        width = 240;
    }
    let left = rect.left;
    if (left < margin) {
        left = margin;
    }
    const availableWidth = viewportWidth - margin - left;
    if (availableWidth > 0) {
        width = Math.min(width, availableWidth);
    }
    let top = rect.bottom + 6;
    dropdown.style.width = `${Math.round(width)}px`;
    dropdown.style.left = `${Math.round(left)}px`;
    dropdown.style.top = `${Math.round(top)}px`;
    dropdown.dataset.placement = "below";
    requestAnimationFrame(() => {
        if (!toolbarComboState.isOpen) return;
        const dropdownHeight = dropdown.offsetHeight;
        if (!dropdownHeight) return;
        if (top + dropdownHeight > viewportHeight - margin && rect.top - 6 - dropdownHeight >= margin) {
            const aboveTop = Math.max(rect.top - 6 - dropdownHeight, margin);
            dropdown.style.top = `${Math.round(aboveTop)}px`;
            dropdown.dataset.placement = "above";
        }
    });
}

function handleToolbarComboReposition() {
    if (toolbarComboState.isOpen) {
        positionToolbarComboDropdown();
    }
}

function setupToolbarComboRepositionListeners() {
    if (window.__toolbarComboRepositionBound) return;
    window.addEventListener("resize", handleToolbarComboReposition);
    window.addEventListener("scroll", handleToolbarComboReposition, true);
    const menuEl = document.getElementById("menu");
    if (menuEl) {
        menuEl.addEventListener("scroll", handleToolbarComboReposition);
    }
    window.__toolbarComboRepositionBound = true;
}

function openToolbarComboDropdown(options = {}) {
    const { dropdown, root, input, select } = getToolbarComboElements();
    if (!dropdown || !root || !input || !select) return;
    if (!select.options.length) return;
    if (!options.preserveFilter) {
        toolbarComboState.filterText = options.filterText || "";
    }
    toolbarComboState.isOpen = true;
    dropdown.classList.remove("hidden");
    root.classList.add("open");
    input.setAttribute("aria-expanded", "true");
    renderToolbarComboDropdown();
    positionToolbarComboDropdown();
}

function closeToolbarComboDropdown(resetFilter = true) {
    const { dropdown, root, input } = getToolbarComboElements();
    if (!dropdown || !root || !input) return;
    toolbarComboState.isOpen = false;
    toolbarComboState.highlightedIndex = -1;
    dropdown.classList.add("hidden");
    dropdown.dataset.placement = "";
    root.classList.remove("open");
    input.setAttribute("aria-expanded", "false");
    if (resetFilter) {
        toolbarComboState.filterText = "";
        resetToolbarComboInput();
    }
}

function moveToolbarComboHighlight(delta) {
    const options = toolbarComboState.filteredOptions;
    if (!options.length) return;
    let idx = toolbarComboState.highlightedIndex;
    if (idx < 0) idx = 0;
    idx = (idx + delta + options.length) % options.length;
    toolbarComboState.highlightedIndex = idx;
    updateToolbarComboActiveOption();
}

function toggleToolbarComboValue(value) {
    if (!value || value === "__upload__") return;
    if (isToolbarComboValueSelected(value)) {
        removeToolbarComboSelection(value);
    } else {
        addToolbarComboSelection(value);
    }
}

function handleToolbarComboKeydown(evt) {
    switch (evt.key) {
        case "ArrowDown":
            evt.preventDefault();
            if (!toolbarComboState.isOpen) {
                openToolbarComboDropdown({ filterText: toolbarComboState.filterText });
            }
            moveToolbarComboHighlight(1);
            break;
        case "ArrowUp":
            evt.preventDefault();
            if (!toolbarComboState.isOpen) {
                openToolbarComboDropdown({ filterText: toolbarComboState.filterText });
            }
            moveToolbarComboHighlight(-1);
            break;
        case "Enter": {
            if (!toolbarComboState.isOpen) {
                closeToolbarComboDropdown();
                return;
            }
            evt.preventDefault();
            const choice = toolbarComboState.filteredOptions[toolbarComboState.highlightedIndex];
            if (choice) {
                if (choice.value === "__upload__") {
                    triggerToolbarComboUpload();
                } else {
                    toggleToolbarComboValue(choice.value);
                }
            }
            break;
        }
        case "Escape":
            evt.preventDefault();
            closeToolbarComboDropdown();
            break;
        case "Tab":
            closeToolbarComboDropdown();
            break;
        default:
            break;
    }
}

function handleToolbarComboDocumentClick(evt) {
    const { root } = getToolbarComboElements();
    if (!root) return;
    if (!toolbarComboState.isOpen) return;
    if (!root.contains(evt.target)) {
        closeToolbarComboDropdown();
    }
}

function initToolbarComboUI() {
    const els = getToolbarComboElements();
    const { select, input, toggle } = els;
    if (!select || select.__comboInitialized) return;
    select.__comboInitialized = true;

    resetToolbarComboInput();
    renderToolbarComboDropdown();
    updateToolbarComboSelectedList();
    setupToolbarComboRepositionListeners();

    if (input) {
        input.addEventListener("focus", () => {
            openToolbarComboDropdown({ filterText: "" });
        });
        input.addEventListener("input", (evt) => {
            toolbarComboState.filterText = evt.target.value;
            openToolbarComboDropdown({ preserveFilter: true });
            renderToolbarComboDropdown();
        });
        input.addEventListener("keydown", handleToolbarComboKeydown);
    }

    if (toggle) {
        toggle.addEventListener("click", () => {
            if (toolbarComboState.isOpen) {
                closeToolbarComboDropdown();
            } else {
                if (input) input.focus();
                openToolbarComboDropdown({ filterText: "" });
            }
        });
    }

    if (!window.__toolbarComboDocListener) {
        document.addEventListener("mousedown", handleToolbarComboDocumentClick);
        window.__toolbarComboDocListener = true;
    }
}

function refreshToolbarComboUI() {
    resetToolbarComboInput();
    updateToolbarComboSelectedList();
    if (toolbarComboState.isOpen) {
        renderToolbarComboDropdown();
    }
}

function addToolbarComboSelection(value, options = {}) {
    if (!value || value === "__upload__") return;
    if (toolbarComboState.selectedValues.includes(value)) return;
    toolbarComboState.selectedValues.push(value);
    const label = getToolbarComboOptionLabel(value) || value;
    const prevMeta = toolbarComboState.selectionMeta.get(value) || {};
    toolbarComboState.selectionMeta.set(value, { ...prevMeta, label });
    updateToolbarComboSelectedList();
    if (toolbarComboState.isOpen) renderToolbarComboDropdown();
    if (!options.silent) {
        applySelectedMixins();
    }
}

function removeToolbarComboSelection(value, options = {}) {
    const idx = toolbarComboState.selectedValues.indexOf(value);
    if (idx === -1) return;
    toolbarComboState.selectedValues.splice(idx, 1);
    if (!options.keepMeta) {
        toolbarComboState.selectionMeta.delete(value);
    }
    updateToolbarComboSelectedList();
    if (toolbarComboState.isOpen) renderToolbarComboDropdown();
    if (!options.silent) {
        applySelectedMixins();
    }
}

function setToolbarComboSelections(values, options = {}) {
    const unique = [];
    (values || []).forEach((value) => {
        if (!value || value === "__upload__") return;
        if (!unique.includes(value)) unique.push(value);
    });
    toolbarComboState.selectedValues = unique;
    const nextMeta = new Map();
    unique.forEach((value) => {
        const existing = toolbarComboState.selectionMeta.get(value) || {};
        const label = getToolbarComboOptionLabel(value) || value;
        nextMeta.set(value, { ...existing, label });
    });
    toolbarComboState.selectionMeta = nextMeta;
    updateToolbarComboSelectedList();
    if (toolbarComboState.isOpen) renderToolbarComboDropdown();
    if (!options.silent) {
        applySelectedMixins();
    }
}

window.getToolbarComboSelectionValues = function () {
    return toolbarComboState.selectedValues.slice();
};

function getFirstSelectableComboValue() {
    const { select } = getToolbarComboElements();
    if (!select) return null;
    const opt = Array.from(select.options).find((o) => o.value && o.value !== "__upload__");
    return opt ? opt.value : null;
}

function getDefaultComboSelectionValues() {
    if (Array.isArray(window.queryComboValues) && window.queryComboValues.length) {
        return window.queryComboValues.slice();
    }
    const first = getFirstSelectableComboValue();
    return first ? [first] : [];
}

// --- Uploaded mixins persistence helpers ---
// Stored in localStorage under 'uploadedMixins' as [{ id, title, content }]
function getUploadedMixins() {
    try {
        const raw = localStorage.getItem("uploadedMixins");
        if (!raw) return [];
        const arr = JSON.parse(raw);
        return Array.isArray(arr) ? arr : [];
    } catch {
        return [];
    }
}

function setUploadedMixins(list) {
    try {
        localStorage.setItem("uploadedMixins", JSON.stringify(list || []));
    } catch {
        // ignore storage errors
    }
}

function upsertUploadedMixin(id, title, content) {
    const list = getUploadedMixins();
    const idx = list.findIndex((x) => x && x.id === id);
    const entry = { id, title: title || id, content: content || "" };
    if (idx >= 0) list[idx] = entry; else list.push(entry);
    setUploadedMixins(list);
    return entry;
}

function parseTitleFromYaml(text, fallback) {
    try {
        const lines = String(text).split(/\r?\n/);
        for (let line of lines) {
            const clean = line.replace(/#.*/, "").trim();
            if (!clean) continue;
            const m = clean.match(/^title\s*:\s*(.+)$/i);
            if (m) {
                let val = m[1].trim();
                if ((val.startsWith('"') && val.endsWith('"')) || (val.startsWith("'") && val.endsWith("'"))) {
                    val = val.slice(1, -1);
                }
                return val || fallback || "Uploaded";
            }
        }
    } catch {
        /* ignore */
    }
    return fallback || "Uploaded";
}

function ensureUploadOptionsInCombo(sel) {
    if (!sel) return;
    // Remove stale uploaded options that are no longer in storage
    const current = new Set(getUploadedMixins().map((u) => `uploaded::${u.id}`));
    Array.from(sel.querySelectorAll("option")).forEach((opt) => {
        if (opt && typeof opt.value === "string" && opt.value.startsWith("uploaded::")) {
            if (!current.has(opt.value)) opt.remove();
        }
    });
    // Append existing uploaded mixins as options
    const uploaded = getUploadedMixins();
    for (const u of uploaded) {
        const val = `uploaded::${u.id}`;
        const existing = sel.querySelector(`option[value='${CSS.escape(val)}']`);
        if (!existing) {
            const opt = document.createElement("option");
            opt.value = val;
            opt.textContent = u.title || u.id;
            sel.appendChild(opt);
        } else {
            const label = u.title || u.id;
            if (existing.textContent !== label) existing.textContent = label;
        }
    }
    // Always ensure the upload sentinel exists as last option
    const sentinelVal = "__upload__";
    let sentinel = sel.querySelector(`option[value='${CSS.escape(sentinelVal)}']`);
    if (!sentinel) {
        sentinel = document.createElement("option");
        sentinel.value = sentinelVal;
        sentinel.textContent = "Upload file…";
        sel.appendChild(sentinel);
    }
    refreshToolbarComboUI();
}

function removeUploadedMixin(id) {
    const list = getUploadedMixins();
    const next = list.filter((x) => x && x.id !== id);
    setUploadedMixins(next);
}

function clearUploadedMixins() {
    setUploadedMixins([]);
}

// Spinner helpers
window.showSpinner = function () {
    const el = document.getElementById("spinner");
    console.log("Spinner start");
    if (!el) return;
    el.classList.remove("hidden");
    el.setAttribute("aria-hidden", "false");
};
window.hideSpinner = function () {
    const el = document.getElementById("spinner");
    console.log("Spinner stop");
    if (!el) return;
    el.classList.add("hidden");
    el.setAttribute("aria-hidden", "true");
};

function installCreateSvgSpinnerWrapper(fnName) {
    const fn = window[fnName];
    if (typeof fn !== "function") return;
    if (fn.__wrappedWithSpinner) return;
    const raw = fn;
    const wrapped = async function (...args) {
        try {
            if (typeof window.showSpinner === "function") window.showSpinner();
            // Yield to the browser so the spinner can paint before heavy work
            await new Promise((r) => requestAnimationFrame(() => r()));
            // A second rAF improves reliability on some browsers
            await new Promise((r) => requestAnimationFrame(() => r()));
            let res = raw.apply(this, args);
            if (res && typeof res.then === "function") {
                res = await res;
            }
            return res;
        } finally {
            if (typeof window.hideSpinner === "function") window.hideSpinner();
        }
    };
    wrapped.__wrappedWithSpinner = true;
    wrapped.__original = raw;
    window[fnName] = wrapped;
}

function getActiveCreateSvgFunction() {
    return typeof window.createSvgExt === "function" ? window.createSvgExt : null;
}

async function callCreateSvgWithMode(arg, filterTexts, blacklistIds, depthOverride) {
    const fn = getActiveCreateSvgFunction();
    if (!fn) return null;
    await ensureFormatMixinLoaded();
    return fn(
        arg,
        getCombinedMixins(),
        Number.isFinite(depthOverride) ? depthOverride : window.defaultDepth,
        filterTexts || [],
        blacklistIds || [],
        window.hideCommentsEnabled,
        window.debug
    );
}

function initPage() {
    window.addEventListener("resize", positionBlacklistCollector);
    window.addEventListener("DOMContentLoaded", positionBlacklistCollector);
    setCollectorPanelVisibility(false);
    setCommentLegendPanelVisibility(false);
    initCollectorVisibilityGuard();
    // Attach dummy handler for toolbar combo box
    document.addEventListener("DOMContentLoaded", function () {
        const selectedPanel = document.getElementById("selected-collector");
        if (selectedPanel) {
            selectedPanel.classList.remove("hidden");
            selectedPanel.setAttribute("aria-hidden", "false");
        }
        initToolbarComboUI();
        initSelectedCollectorInteractions();
        const { select: combo, root: comboRoot } = getToolbarComboElements();
        if (combo) {
            if (comboRoot) {
                comboRoot.style.display = window.queryOptions ? "" : "none";
                if (!window.queryOptions) {
                    closeToolbarComboDropdown();
                }
            } else if (!window.queryOptions) {
                combo.style.display = "none";
            }
            if (typeof window.loadComboOptionsFromYaml === "function") {
                window.loadComboOptionsFromYaml();
            } else {
                ensureUploadOptionsInCombo(combo);
                const defaults = getDefaultComboSelectionValues();
                if (defaults.length) {
                    setToolbarComboSelections(defaults, { silent: true });
                }
                updateToolbarComboSelectedList();
                applySelectedMixins();
            }
        }

        // Manage Uploads UI wiring
        const manageBtn = document.getElementById("btn-manage-uploads");
        const popup = document.getElementById("uploads-popup");
        const listEl = document.getElementById("uploads-list");
        const closeBtn = document.getElementById("uploads-close-btn");
        const doneBtn = document.getElementById("uploads-done");
        const clearAllBtn = document.getElementById("uploads-clear-all");
        const connectedBtn = document.getElementById("btn-connected");
        const clearExpandedLink = document.getElementById("link-clear-expanded");
        const clearBlacklistLink = document.getElementById("link-clear-blacklist");
        const toolbarHandle = document.getElementById("toolbar-collapse-handle");
        const menuWrapper = document.getElementById("menu-wrapper");

        function refreshComboAfterStorageChange(removedId) {
            const sel = document.getElementById("toolbar-combo");
            if (!sel) return;
            ensureUploadOptionsInCombo(sel);

            const storedUploaded = new Set(
                getUploadedMixins().map((u) => `uploaded::${u.id}`)
            );

            Array.from(toolbarComboState.contentCache.keys()).forEach((key) => {
                if (key.startsWith("uploaded::") && !storedUploaded.has(key)) {
                    toolbarComboState.contentCache.delete(key);
                }
            });

            if (removedId) {
                const removedVal = `uploaded::${removedId}`;
                toolbarComboState.contentCache.delete(removedVal);
                removeToolbarComboSelection(removedVal, { silent: true });
            }

            const available = new Set(
                Array.from(sel.options).map((o) => o.value)
            );
            const filtered = toolbarComboState.selectedValues.filter((value) => {
                if (value.startsWith("uploaded::")) {
                    return storedUploaded.has(value);
                }
                return available.has(value);
            });
            if (filtered.length !== toolbarComboState.selectedValues.length) {
                setToolbarComboSelections(filtered, { silent: true });
            }

            if (!toolbarComboState.selectedValues.length) {
                const fallback = getFirstSelectableComboValue();
                if (fallback) {
                    addToolbarComboSelection(fallback, { silent: true });
                }
            }
            updateToolbarComboSelectedList();
            applySelectedMixins();
        }

        function populateUploadsPopup() {
            if (!listEl) return;
            const items = getUploadedMixins();
            listEl.innerHTML = "";
            if (!items.length) {
                const empty = document.createElement("div");
                empty.textContent = "No uploaded mixins stored.";
                empty.style.color = "#666";
                listEl.appendChild(empty);
                return;
            }
            items.forEach((it) => {
                const row = document.createElement("div");
                row.style.display = "flex";
                row.style.alignItems = "center";
                row.style.justifyContent = "space-between";
                row.style.gap = "8px";
                row.style.padding = "4px 0";

                const info = document.createElement("div");
                info.style.display = "flex";
                info.style.flexDirection = "column";
                const title = document.createElement("span");
                title.textContent = it.title || it.id;
                title.style.fontWeight = "600";
                const sub = document.createElement("span");
                sub.textContent = it.id;
                sub.style.fontSize = "0.85em";
                sub.style.color = "#666";
                info.appendChild(title);
                info.appendChild(sub);

                const actions = document.createElement("div");
                const del = document.createElement("button");
                del.className = "tool-btn";
                del.title = "Remove";
                del.innerHTML = '<i class="fa-solid fa-trash"></i>';
                del.addEventListener("click", () => {
                    removeUploadedMixin(it.id);
                    populateUploadsPopup();
                    refreshComboAfterStorageChange(it.id);
                });
                actions.appendChild(del);

                row.appendChild(info);
                row.appendChild(actions);
                listEl.appendChild(row);
            });
        }

        function showUploadsPopup() {
            if (!popup) return;
            populateUploadsPopup();
            popup.classList.remove("hidden");
            popup.style.display = "block";
            popup.focus();
        }
        function hideUploadsPopup() {
            if (!popup) return;
            popup.classList.add("hidden");
            popup.style.display = "none";
        }

        if (manageBtn) manageBtn.addEventListener("click", showUploadsPopup);
        if (closeBtn) closeBtn.addEventListener("click", hideUploadsPopup);
        if (doneBtn) doneBtn.addEventListener("click", hideUploadsPopup);
        if (clearAllBtn) clearAllBtn.addEventListener("click", () => {
            clearUploadedMixins();
            populateUploadsPopup();
            refreshComboAfterStorageChange(null);
        });
        if (connectedBtn) {
            connectedBtn.addEventListener("click", () => {
                if (typeof window.createSvgForConnected !== "function") {
                    alert("Connected render is not available yet.");
                    return;
                }
                connectedModeActive = !connectedModeActive;
                connectedBtn.classList.toggle("active", connectedModeActive);
                const savedBadgeState = collectBadgeIds("badge-list");
                const savedBlacklistState = collectBadgeIds("blacklist-list");
                if (connectedModeActive) {
                    regenerateSvgWithConnectedOnce(savedBadgeState, savedBlacklistState);
                } else {
                    regenerateSvgWithState(savedBadgeState, savedBlacklistState);
                }
            });
        }
        if (toolbarHandle && menuWrapper) {
            const toggleToolbar = () => {
                const collapsed = menuWrapper.classList.toggle("toolbar-collapsed");
                toolbarHandle.title = collapsed ? "Expand Toolbar" : "Collapse Toolbar";
                toolbarHandle.setAttribute(
                    "aria-label",
                    collapsed ? "Expand Toolbar" : "Collapse Toolbar"
                );
            };
            toolbarHandle.addEventListener("click", toggleToolbar);
            toolbarHandle.addEventListener("keydown", (evt) => {
                if (evt.key === "Enter" || evt.key === " ") {
                    evt.preventDefault();
                    toggleToolbar();
                }
            });
        }
        if (clearExpandedLink) {
            clearExpandedLink.addEventListener("click", (evt) => {
                evt.preventDefault();
                const list = document.getElementById("badge-list");
                if (list) list.innerHTML = "";
                if (typeof reloadSvgFromBadges === "function") {
                    reloadSvgFromBadges();
                }
                positionBlacklistCollector();
            });
        }
        if (clearBlacklistLink) {
            clearBlacklistLink.addEventListener("click", (evt) => {
                evt.preventDefault();
                blacklist = [];
                window.blacklist = [];
                updateBlacklistUI();
                if (typeof reloadSvgFromBadges === "function") {
                    reloadSvgFromBadges();
                }
                positionBlacklistCollector();
            });
        }
        // Dismiss on ESC and outside click
        if (popup) {
            popup.addEventListener("keydown", function (e) {
                if (e.key === "Escape") hideUploadsPopup();
            });
            document.addEventListener("mousedown", function (e) {
                if (!popup.classList.contains("hidden") && !popup.contains(e.target) && e.target !== manageBtn) {
                    hideUploadsPopup();
                }
            });
        }
    });
    // Also reposition after collector content changes
    const observer = new MutationObserver(positionBlacklistCollector);
    observer.observe(document.getElementById("collector"), {
        childList: true,
        subtree: true,
    });
    const selectedPanel = document.getElementById("selected-collector");
    if (selectedPanel) {
        const selectedObserver = new MutationObserver(positionBlacklistCollector);
        selectedObserver.observe(selectedPanel, {
            childList: true,
            subtree: true,
        });
    }
    (function () {
        // Polyfill: CSS.escape (minimal) for older browsers
        if (typeof CSS === "undefined" || typeof CSS.escape !== "function") {
            window.CSS = window.CSS || {};
            CSS.escape = function (sel) {
                return String(sel).replace(/[^a-zA-Z0-9_\-]/g, "\\$&");
            };
        }

        // Collect texts from a node (including nested text/tspan)
        function collectTextsFromNode(node, acc) {
            if (!node) return;
            // caption companion
            if (node.id) {
                const capt = document.getElementById(`${node.id}_capt`);
                if (capt && String(capt.tagName).toLowerCase() === "text") {
                    const t = (capt.textContent || "").trim();
                    if (t)
                        t.split(/\r?\n/)
                            .map((s) => s.trim())
                            .filter(Boolean)
                            .forEach((x) => acc.add(x));
                }
            }
            // direct text nodes
            if (node.nodeType === Node.ELEMENT_NODE) {
                const tag = String(node.tagName).toLowerCase();
                if (tag === "text" || tag === "tspan") {
                    const t = (node.textContent || "").trim();
                    if (t)
                        t.split(/\r?\n/)
                            .map((s) => s.trim())
                            .filter(Boolean)
                            .forEach((x) => acc.add(t));
                }
                // aria-label/title or title element
                const titleAttr =
                    node.getAttribute &&
                    (node.getAttribute("title") ||
                        node.getAttribute("aria-label"));
                if (titleAttr) acc.add(titleAttr.trim());
                // traverse children
                for (let i = 0; i < node.childNodes.length; i++) {
                    collectTextsFromNode(node.childNodes[i], acc);
                }
            } else if (node.nodeType === Node.TEXT_NODE) {
                const t = (node.textContent || "").trim();
                if (t) acc.add(t);
            }
        }

        // Robust: aggregate filter texts for all elements in the clicked box
        function getFilterTextsForBox(el) {
            if (!el || !el.id) return [];

            const svg = getSvg();
            if (!svg) return getTextContentArray(el);

            const boxId = getBoxPrefix(el.id);

            // Prefer exact container match
            let container = svg.querySelector(`#${CSS.escape(boxId)}`);
            // Fallback: any descendant whose id starts with boxId
            const acc = new Set();
            if (container) {
                collectTextsFromNode(container, acc);
            } else {
                const related = svg.querySelectorAll(
                    `[id^="${CSS.escape(boxId)}"]`
                );
                related.forEach((node) => collectTextsFromNode(node, acc));
            }

            // If nothing collected, include element-level text arrays across related nodes
            if (acc.size === 0) {
                const related = svg.querySelectorAll(
                    `[id^="${CSS.escape(boxId)}"]`
                );
                related.forEach((node) => {
                    getTextContentArray(node).forEach((t) => acc.add(t));
                });
            }

            // Always include a readable caption for the box itself
            acc.add(getCaptionForId(boxId));

            // As last resort, include ids to make filters stable
            const relatedIds = svg.querySelectorAll(
                `[id^="${CSS.escape(boxId)}"]`
            );
            relatedIds.forEach((node) => acc.add(node.id));

            return Array.from(acc)
                .map((s) => s.trim())
                .filter(Boolean);
        }

        // Global click handler for SVG shapes
        window.shapeClick = function (evt) {
            const el = evt.target;
            if (!el || !el.id) return;

            // If blacklist collector is visible, collect to blacklist; otherwise, collect to expanded collector
            const blacklistBox = document.getElementById("blacklist-collector");
            const blacklistVisible =
                blacklistBox && !blacklistBox.classList.contains("hidden");
            if (blacklistVisible) {
                addToBlacklist(el);
                // Do not forcibly show the blacklist collector here; only toggleBlacklist controls it
                return;
            }

            // --- Existing expanded collector logic ---
            // NEW: ignore click on a parent if a child is already in the collector
            const clickedHid = getBoxPrefix(el.id);
            if (anyBadgeIsChildOf(clickedHid)) {
                console.log(
                    "Ignoring parent click because a child is already selected:",
                    el.id
                );
                return;
            }

            // Snapshot current badges BEFORE applying click changes
            const prevState = getCurrentBadgeState();

            // NEW: if a parent badge exists and user clicked a child, remove parent badge(s) and deselect parent shape(s)
            const ancestorBadges = findAncestorBadgesOf(clickedHid);
            if (ancestorBadges.length > 0) {
                ancestorBadges.forEach((b) => {
                    const hid =
                        b.dataset.hid ||
                        (b.dataset.id ? getBoxPrefix(b.dataset.id) : "");
                    if (hid) deselectElementByHid(hid);
                    b.remove();
                });
            }

            // Simple highlight: toggle a data-selected flag and adjust stroke
            const selected = el.getAttribute("data-selected") === "true";
            el.setAttribute("data-selected", selected ? "false" : "true");
            if (selected) {
                el.setAttribute(
                    "stroke-width",
                    el.getAttribute("data-original-stroke-width") || "3"
                );
                el.removeAttribute("filter");
            } else {
                // Store original stroke width once
                if (!el.hasAttribute("data-original-stroke-width")) {
                    el.setAttribute(
                        "data-original-stroke-width",
                        el.getAttribute("stroke-width") || "3"
                    );
                }
                el.setAttribute(
                    "stroke-width",
                    String(
                        Number(el.getAttribute("data-original-stroke-width")) +
                            2
                    )
                );
            }
            console.log("Clicked item:", el.id);

            // Toggle badge in collector: remove if exists, otherwise add
            const list = document.getElementById("badge-list");
            if (list) {
                const boxId = getBoxPrefix(el.id);
                const existingBadges = findBadgesByBoxId(boxId);
                if (existingBadges.length > 0) {
                    existingBadges.forEach((b) => b.remove());
                    // NEW: refit remaining badges after removal
                    requestAnimationFrame(refitAllBadges);
                } else {
                    const badge = createBadgeForShape(el);
                    list.insertBefore(badge, list.firstChild); // prepend
                    // NEW: fit newly added badge
                    requestAnimationFrame(() => fitBadgeLabel(badge));
                }
            }

            // NEW: push undo state only if state changed
            const newState = getCurrentBadgeState();
            if (!statesEqual(prevState, newState)) {
                undoStack.push(prevState);
            }

            // Replace displayed SVG using createSvg; use YAML input when loaded
            (async () => {
                try {
                    if (typeof createSvg !== "function") {
                        console.error("createSvg is not available.");
                        return;
                    }
                    const canvas = document.getElementById("canvas");
                    if (!canvas) return;

                    // Wait for YAML load if promise exists
                    if (
                        window.inputLoaded &&
                        typeof window.inputLoaded.then === "function"
                    ) {
                        await window.inputLoaded;
                    }

                    // Fallback: toggle arg if no YAML loaded
                    const count =
                        document.querySelectorAll("#badge-list .badge").length;
                    const fallbackArg = count % 2 === 1 ? "1" : "2";
                    const arg =
                        typeof window.input === "string" &&
                        window.input.length > 0
                            ? window.input
                            : fallbackArg;

                    // UPDATED: pass all badge captions currently in the clicked shapes box

                    const filterTexts = getAllBadgeCaptions("badge-list");
                    // Extract ids from blacklisted badges (elements in blacklist)
                    const blacklistIds = blacklist
                        .map((boxId) => {
                            // Try to get the element and its id
                            const el = document.getElementById(boxId);
                            return el ? el.id : boxId;
                        })
                        .filter(Boolean);
                    console.log(
                        "Refreshing SVG: ",
                        filterTexts,
                        "blacklist ids: ",
                        blacklistIds,
                        "comments hidden: ",
                        window.hideCommentsEnabled
                    );
                    let svgStr = await callCreateSvgWithMode(
                        arg,
                        filterTexts,
                        blacklistIds
                    );
                    svgStr =
                        svgStr && typeof svgStr.then === "function"
                            ? await svgStr
                            : svgStr;

                    if (
                        typeof svgStr !== "string" ||
                        !svgStr.trim().startsWith("<svg")
                    ) {
                        console.error(
                            "createSvg did not return a valid SVG string."
                        );
                        console.error(svgStr);
                        return;
                    }

                    // Swap SVG and trigger existing initialization
                    canvas.innerHTML = svgStr;
                    const evtSwap = new Event("htmx:afterSwap", {
                        bubbles: true,
                    });
                    canvas.dispatchEvent(evtSwap);

                    // NEW: after click, observe caption/text changes and refresh when filled
                    observeCaptionAndRefresh(el);
                } catch (e) {
                    console.error("Error updating SVG via createSvg:", e);
                }
            })();
        };

        // NEW: click a badge to remove it and reload the SVG
        window.attachBadgeRemoval();

        window.svgControls = {
            // Ensure toggleBlacklist is always available on window.svgControls
            toggleBlacklist() {
                blacklistMode = !blacklistMode;
                const blist = document.getElementById("blacklist-collector");
                if (blist) {
                    if (blacklistMode) {
                        blist.classList.remove("hidden");
                        blist.setAttribute("aria-hidden", "false");
                        blist.scrollIntoView({
                            behavior: "smooth",
                            block: "nearest",
                        });
                        updateBlacklistUI(); // Always update UI from window.blacklist when showing
                    } else {
                        blist.classList.add("hidden");
                        blist.setAttribute("aria-hidden", "true");
                    }
                }
                updateToolButtons();
                positionBlacklistCollector();
            },
            zoom(factor) {
                state.scale *= factor;
                applyTransform();
            },
            pan(dx, dy) {
                state.tx += dx;
                state.ty += dy;
                applyTransform();
            },
            save() {
                const svg = getSvg();
                if (!svg) return;
                // Clone to avoid mutating live DOM
                const clone = svg.cloneNode(true);
                // Ensure xmlns attributes present for a standalone file
                if (!clone.getAttribute("xmlns")) {
                    clone.setAttribute("xmlns", "http://www.w3.org/2000/svg");
                }
                if (!clone.getAttribute("xmlns:xlink")) {
                    clone.setAttribute(
                        "xmlns:xlink",
                        "http://www.w3.org/1999/xlink"
                    );
                }
                // Preserve viewBox if available; fallback to width/height
                const vb = clone.viewBox && clone.viewBox.baseVal;
                if (!vb) {
                    const w = clone.getAttribute("width") || "800";
                    const h = clone.getAttribute("height") || "600";
                    clone.setAttribute("viewBox", `0 0 ${w} ${h}`);
                }
                const xml = new XMLSerializer().serializeToString(clone);
                const blob = new Blob([xml], {
                    type: "image/svg+xml;charset=utf-8",
                });
                const url = URL.createObjectURL(blob);
                const a = document.createElement("a");
                a.href = url;
                a.download = "canvas.svg";
                document.body.appendChild(a);
                a.click();
                document.body.removeChild(a);
                URL.revokeObjectURL(url);
            },
            toggleMinimap() {
                minimapVisible = !minimapVisible;
                const mmWrap = document.getElementById("minimap");
                if (!mmWrap) return;
                mmWrap.style.display = minimapVisible ? "flex" : "none";
                mmWrap.setAttribute(
                    "aria-hidden",
                    minimapVisible ? "false" : "true"
                );
                if (minimapVisible) {
                    initMinimap();
                    populateMinimapPreview();
                    updateMinimap();
                }
                // Always update the toolbar button state
                updateToolButtons();
            },
            toggleCollector() {
                setCollectorPanelVisibility(!collectorPanelVisible);
            },
            toggleCommentLegendPanel() {
                setCommentLegendPanelVisibility(!commentLegendPanelVisible);
            },
            toggleSelectedCollector() {
                const panel = document.getElementById("selected-collector");
                if (!panel) return;
                const hidden = panel.classList.toggle("hidden");
                panel.setAttribute("aria-hidden", hidden ? "true" : "false");
                updateToolbarComboSelectedList();
                updateToolButtons();
                positionBlacklistCollector();
            },
            // Toggle pan tool
            togglePanTool() {
                panToolActive = !panToolActive;
                applyTransform(); // updates cursor class
                updateToolButtons();
            },
            toggleConnectionAnimation() {
                connectionAnimationEnabled = !connectionAnimationEnabled;
                window.connectionAnimationEnabled = connectionAnimationEnabled;
                if (connectionAnimationEnabled) {
                    installConnectionRunners();
                } else {
                    clearConnectionRunners();
                }
                updateToolButtons();
            },
        };

        // After SVG loads, wrap and initialize sizes
        document.body.addEventListener("htmx:afterSwap", function (evt) {
            if (evt.target && evt.target.id === "canvas") {
                ensureStageWrapped();
                // NEW: only reset state if it's not already preserved
                if (state.scale === 1 && state.tx === 0 && state.ty === 0) {
                    // This is likely an initial load, keep the reset
                    state = { scale: 1, tx: 0, ty: 0 };
                }
                // If state was preserved (non-default values), keep it as-is
                computeBaseSize();
                applyTransform();
                initMinimap();
                populateMinimapPreview();
                updateMinimap();

                // Attach pan handlers to stage
                const stage = getStage();
                if (stage) {
                    stage.addEventListener("mousedown", onStageMouseDown);
                }
                updateToolButtons();
                // NEW: refit badges in case layout changed
                requestAnimationFrame(refitAllBadges);

                // Add pointer cursor to SVG elements with click handlers
                const svg = getSvg();
                if (svg) {
                    // Remove svg-clickable from elements that have an onclick handler
                    svg.querySelectorAll(".svg-clickable").forEach((el) => {
                        if (el.hasAttribute("onclick")) {
                            el.classList.remove("svg-clickable");
                        }
                    });
                    // Add onclick handler to all elements with svg-clickable class
                    svg.querySelectorAll(".svg-clickable").forEach((el) => {
                        if (!el.hasAttribute("onclick")) {
                            el.setAttribute(
                                "onclick",
                                "window.shapeClick(event)"
                            );
                        }
                    });
                    // Also attach handlers for elements that should open external links
                    attachAdditionalLinkHandlers();
                }
                updateCommentLegendPanel();
                installConnectionRunners();
            }
        });

// Attach click handlers to .additionalLink elements inside the current SVG
function attachAdditionalLinkHandlers() {
    const svg = getSvg();
    if (!svg) return;
    const nodes = svg.querySelectorAll(".additionalLink");
    nodes.forEach((node) => {
        try {
            node.style.cursor = "pointer";
        } catch {}
        if (node.__additionalLinkBound) return;
        node.__additionalLinkBound = true;
        node.addEventListener(
            "click",
            function (evt) {
                try {
                    let el = evt.target;
                    while (
                        el &&
                        el !== svg &&
                        !(el.classList && el.classList.contains("additionalLink"))
                    ) {
                        el = el.parentNode;
                    }
                    const url = el && el.getAttribute ? el.getAttribute("data-link") : null;
                    if (url && typeof url === "string") {
                        window.open(url, "_blank", "noopener,noreferrer");
                    } else {
                        console.warn("additionalLink clicked without data-link URL", el);
                    }
                } catch (err) {
                    console.error("Failed to handle additionalLink click:", err);
                } finally {
                    if (evt) {
                        evt.preventDefault();
                        evt.stopPropagation();
                    }
                }
            },
            { capture: true }
        );
    });
}

        // Keep centered and minimap in sync on resize
        window.addEventListener("resize", function () {
            applyTransform(); // recenter and resize stage
            if (minimapVisible) {
                initMinimap();
                populateMinimapPreview();
                updateMinimap();
            }
            // NEW: badges may need refitting on resize
            requestAnimationFrame(refitAllBadges);
        });
        getCanvas().addEventListener("scroll", function () {
            updateMinimap();
        });

        // Global keyboard listeners for Space pan
        window.addEventListener("keydown", onKeyDown, {
            passive: false,
        });
        window.addEventListener("keyup", onKeyUp);
    })();
}

async function loadYamlInput() {
    try {
        if (location.protocol === "file:") {
            console.error(
                "YAML must be served over HTTP(S). Use a local web server to host this page."
            );
            window.input = "";
            window.inputLoaded = Promise.resolve();
            window.currentYamlFile = undefined;
            return;
        }
        const yamlFile = window.queryInput;
        const p = (async () => {
            const resp = await fetch(window.getBasePath() + "/data/" + yamlFile, {
                cache: "no-cache",
            });
            if (!resp.ok) throw new Error("HTTP " + resp.status);
            window.input = await resp.text();
            window.currentYamlFile = yamlFile;
        })().catch((e) => {
            console.error("Failed to fetch " + window.queryInput + ":", e);
            window.input = "";
            window.currentYamlFile = undefined;
        });
        window.inputLoaded = p.then(() => undefined);
    } catch (e) {
        console.error("Unexpected error loading " + window.queryInput + ":", e);
        window.input = "";
        window.inputLoaded = Promise.resolve();
        window.currentYamlFile = undefined;
    }
}

async function loadSVGFromWasm() {
    const canvas = document.getElementById("canvas");
    if (!canvas) return;

    if (typeof Go !== "function") {
        console.error(
            "Go runtime not available. Ensure ./wasm_exec.js is loaded before this script."
        );
        return;
    }
    // Warn if running from file:// which cannot fetch WASM
    if (location.protocol === "file:") {
        console.error(
            "WASM must be served over HTTP(S). Use a local web server to host this page."
        );
        return;
    }

    const go = new Go();

    // Fetch boxes.wasm with streaming fallback
    let resp;
    try {
        resp = await fetch(window.getBasePath() + "/wasm/boxes_1.1.2.wasm", {
            cache: "no-cache",
        });
        if (!resp.ok) throw new Error("HTTP " + resp.status);
    } catch (e) {
        console.error("Failed to fetch boxes.wasm:", e);
        return;
    }

    let instance;
    try {
        if (WebAssembly.instantiateStreaming) {
            try {
                ({ instance } = await WebAssembly.instantiateStreaming(
                    resp,
                    go.importObject
                ));
            } catch {
                const buf = await resp.arrayBuffer();
                ({ instance } = await WebAssembly.instantiate(
                    buf,
                    go.importObject
                ));
            }
        } else {
            const buf = await resp.arrayBuffer();
            ({ instance } = await WebAssembly.instantiate(
                buf,
                go.importObject
            ));
        }
    } catch (e) {
        console.error("Error instantiating boxes.wasm:", e);
        return;
    }

    // Start the runtime, then wait for renderer functions to appear.
    go.run(instance).catch((e) => console.warn("go.run finished/failed:", e));

    // Wait until createSvgExt or createSvgForConnected is exposed by the Go code (poll with timeout)
    const timeoutMs = 5000;
    const start = Date.now();
    while (
        typeof window.createSvgExt !== "function" &&
        typeof window.createSvgForConnected !== "function"
    ) {
        if (Date.now() - start > timeoutMs) {
            console.error(
                'createSvgExt/createSvgForConnected not exposed by boxes.wasm within timeout. Ensure js.Global().Set("createSvgExt", fn) or js.Global().Set("createSvgForConnected", fn).'
            );
            return;
        }
        await new Promise((r) => setTimeout(r, 50));
    }

    // If the URL had a compressed ?z= permalink, decompress and apply now that WASM is ready.
    // Then re-run loadYamlInput and loadComboOptionsFromYaml so all state is correctly restored.
    if (window.queryZ) {
        applyCompressedQueryParam();
        await loadYamlInput();
        if (typeof window.loadComboOptionsFromYaml === "function") {
            await window.loadComboOptionsFromYaml();
        }
    }

    // Install spinner wrapper once renderers are available
    installCreateSvgSpinnerWrapper("createSvgExt");
    installCreateSvgSpinnerWrapper("createSvgForConnected");

    let svgStr;
    try {
        // Wait for YAML load if available, then pass it to the active renderer
        if (
            window.inputLoaded &&
            typeof window.inputLoaded.then === "function"
        ) {
            await window.inputLoaded;
        }
        // Get current expanded and blacklisted IDs
        const badgeList = document.getElementById("badge-list");
        const filterTexts = badgeList ? Array.from(badgeList.querySelectorAll(".badge")).map(b => b.dataset.hid).filter(Boolean) : [];
        // Use window.blacklist if set, otherwise fallback to DOM
        const blacklistIds = (window.blacklist && Array.isArray(window.blacklist)) ? window.blacklist : (document.getElementById("blacklist-list") ? Array.from(document.getElementById("blacklist-list").querySelectorAll(".badge")).map(b => b.dataset.hid).filter(Boolean) : []);
        const initialArg =
            typeof window.input === "string" && window.input.length > 0
                ? window.input
                : "";
        console.log(
            "Refreshing SVG: ",
            filterTexts,
            "blacklist ids: ",
            blacklistIds,
            "comments hidden: ",
            window.hideCommentsEnabled
        );
        const res = await callCreateSvgWithMode(
            initialArg,
            filterTexts,
            blacklistIds
        );
        const resolved = res && typeof res.then === "function" ? await res : res;
        const normalized = normalizeCreateSvgResult(
            resolved,
            filterTexts,
            blacklistIds
        );
        svgStr = normalized.svgStr;
        applyExpandedAndBlacklistState(
            normalized.expanded,
            normalized.blacklisted
        );
    } catch (e) {
        console.error("Error calling createSvg:", e);
        return;
    }

    if (typeof svgStr !== "string" || !svgStr.trim().startsWith("<svg")) {
        console.error("createSvg did not return a valid SVG string.");
        console.error(svgStr);
        return;
    }

    canvas.innerHTML = svgStr;
    const evt = new Event("htmx:afterSwap", { bubbles: true });
    canvas.dispatchEvent(evt);
}

function applyQueryParams(params) {
    // Store input (existing)
    if (params.has("input")) {
        const raw = params.get("input") || "";
        const content = raw.replace(/\+/g, " ");
        window.queryInput = content;
    }
    // parse 'options' query param for dynamic combo source
    if (params.has("options")) {
        const rawOptions = params.get("options") || "";
        window.queryOptions = rawOptions.replace(/\+/g, " ");
    } else {
        window.queryOptions = undefined;
    }
    // parse 'debug' query param and store as global boolean
    const rawDebug = params.get("debug");
    const truthy = ["true", "1", "yes", "on"];
    const falsy = ["false", "0", "no", "off"];
    let val = false;
    if (rawDebug != null) {
        const normalized = String(rawDebug).toLowerCase();
        if (!normalized) {
            val = true;
        } else if (falsy.includes(normalized)) {
            val = false;
        } else if (truthy.includes(normalized)) {
            val = true;
        } else {
            val = true;
        }
    }
    window.debug = val;

    // Parse combo, expandedIds, blacklistedIds
    // Combo box (allow comma-separated or repeated params)
    const comboParams = params.getAll("combo");
    if (comboParams.length) {
        const comboList = [];
        const seen = new Set();
        comboParams.forEach((chunk) => {
            String(chunk || "")
                .split(",")
                .map((part) => part.trim())
                .filter(Boolean)
                .forEach((val) => {
                    if (!seen.has(val)) {
                        seen.add(val);
                        comboList.push(val);
                    }
                });
        });
        window.queryComboValues = comboList;
        window.queryCombo = comboList.length ? comboList[0] : undefined;
    } else {
        window.queryComboValues = [];
        window.queryCombo = undefined;
    }

    // Expanded IDs (badges) — works before and after DOMContentLoaded
    if (params.has("expandedIds")) {
        const expandedIds = params.get("expandedIds").split(",").map(s => s.trim()).filter(Boolean);
        const applyExpanded = function () {
            const list = document.getElementById("badge-list");
            if (list && expandedIds.length) {
                list.innerHTML = "";
                expandedIds.forEach(hid => {
                    const span = document.createElement("span");
                    span.className = "badge";
                    span.dataset.hid = hid;
                    const label = document.createElement("span");
                    label.textContent = (window.getCaptionForId ? window.getCaptionForId(hid) : hid);
                    span.appendChild(label);
                    list.appendChild(span);
                });
                requestAnimationFrame(window.refitAllBadges || (()=>{}));
            }
        };
        if (document.readyState === "loading") {
            document.addEventListener("DOMContentLoaded", applyExpanded);
        } else {
            applyExpanded();
        }
    }

    // Blacklisted IDs
    if (params.has("blacklistedIds")) {
        const blacklistedIds = params.get("blacklistedIds").split(",").map(s => s.trim()).filter(Boolean);
        if (document.readyState === "loading") {
            document.addEventListener("DOMContentLoaded", function () {
                window.blacklist = blacklistedIds;
            });
        } else {
            window.blacklist = blacklistedIds;
        }
    }

    // Search IDs
    if (params.has("search_ids")) {
        window.querySearchIds = params.get("search_ids").split(",").map(s => s.trim()).filter(Boolean);
    } else {
        window.querySearchIds = [];
    }
}

function handleInputQueryParam() {
    try {
        const params = new URLSearchParams(window.location.search);
        // If a compressed permalink is present, defer processing until WASM is ready
        if (params.has("z")) {
            window.queryZ = params.get("z");
            return;
        }
        applyQueryParams(params);
    } catch (e) {
        console.error("Failed to read query params:", e);
        if (typeof window.queryInput === "undefined") {
            window.queryInput = undefined;
        }
        window.debug = false;
    }
}

function applyCompressedQueryParam() {
    if (!window.queryZ) return;
    if (typeof window.decompressString !== "function") {
        console.error("decompressString not available from WASM");
        return;
    }
    try {
        const decompressed = window.decompressString(window.queryZ);
        if (!decompressed) {
            console.error("Failed to decompress permalink z param");
            return;
        }
        console.log("[permalink] decompressed z param:", decompressed);
        applyQueryParams(new URLSearchParams(decompressed));
    } catch (e) {
        console.error("Failed to apply compressed query param:", e);
    }
}

// NEW: Load combo-box options from a YAML mapping specified by the 'options' query param
window.loadComboOptionsFromYaml = async function () {
    try {
        // Only proceed if an 'options' param was provided
        const src = window.queryOptions;
        if (!src) {
            // Even if no remote options source, still ensure upload entries are present
            const { select: selNoSrc, root: comboRoot } = getToolbarComboElements();
            if (comboRoot) comboRoot.style.display = "";
            if (selNoSrc) {
                ensureUploadOptionsInCombo(selNoSrc);
            }
            return;
        }
        if (location.protocol === "file:") {
            console.error("Options YAML must be served over HTTP(S).");
            return;
        }
        const resp = await fetch(window.getBasePath() + "/data/" + src, { cache: "no-cache" });
        if (!resp.ok) throw new Error("HTTP " + resp.status);
        const text = await resp.text();
        // Parse a simple YAML mapping: label: value (labels may be quoted)
        const entries = parseSimpleYamlMapping(text);
        const { select: sel, root: comboRoot } = getToolbarComboElements();
        if (!sel) return;
        // Ensure it's visible when options are available
        if (comboRoot) comboRoot.style.display = "";
        // Replace existing options with dynamically loaded ones
        sel.innerHTML = "";
        for (const [label, value] of entries) {
            const opt = document.createElement("option");
            opt.value = value;
            opt.textContent = label;
            sel.appendChild(opt);
        }
        // Add uploaded mixins and the upload sentinel
        ensureUploadOptionsInCombo(sel);
        const desiredSelection = getDefaultComboSelectionValues();
        setToolbarComboSelections(desiredSelection, { silent: true });
        updateToolbarComboSelectedList();
        await applySelectedMixins();
        refreshToolbarComboUI();
    } catch (e) {
        console.error("Failed to load combo options from YAML:", e);
    }
};

// Helper: very small YAML mapping parser for lines like: "Label": value
function parseSimpleYamlMapping(text) {
    const lines = String(text).split(/\r?\n/);
    const out = [];
    for (let line of lines) {
        // Strip comments and trim
        line = line.replace(/#.*/, "").trim();
        if (!line) continue;
        // Match key: value allowing quoted keys and values
        const m = line.match(/^([^:]+):\s*(.*)$/);
        if (!m) continue;
        let key = m[1].trim();
        let val = m[2].trim();
        // Remove surrounding quotes from key
        if ((key.startsWith('"') && key.endsWith('"')) || (key.startsWith("'") && key.endsWith("'"))) {
            key = key.slice(1, -1);
        }
        // Remove surrounding quotes from value
        if ((val.startsWith('"') && val.endsWith('"')) || (val.startsWith("'") && val.endsWith("'"))) {
            val = val.slice(1, -1);
        }
        out.push([key, val]);
    }
    return out;
}

function attachBadgeRemoval() {
    const list = document.getElementById("badge-list");
    if (!list) return;
    list.addEventListener("click", (evt) => {
        const badge = evt.target.closest(".badge");
        if (!badge) return;

        // Snapshot current badges BEFORE removal
        const prevState = getCurrentBadgeState();

        // Deselect related SVG element(s) for visual consistency
        if (badge.dataset.hid) deselectElementByHid(badge.dataset.hid);
        if (badge.dataset.id) {
            const el = document.getElementById(badge.dataset.id);
            if (el && el.getAttribute("data-selected") === "true") {
                el.setAttribute("data-selected", "false");
                el.setAttribute(
                    "stroke-width",
                    el.getAttribute("data-original-stroke-width") || "3"
                );
                el.removeAttribute("filter");
            }
        }

        badge.remove();
        // NEW: refit after removal
        requestAnimationFrame(refitAllBadges);

        // Push undo entry and reload
        undoStack.push(prevState);
        reloadSvgFromBadges();
    });
}

function getAllBadgeCaptions(boxId) {
    const list = document.getElementById(boxId);
    if (!list) return [];
    const badges = list.querySelectorAll(".badge");
    const out = new Set();
    badges.forEach((b) => {
        const id = b.dataset.id;
        if (id) out.add(id);
        const hid = b.dataset.hid;
        if (hid) out.add(hid);
    });
    return Array.from(out);
}

// NEW: snapshot current badge state (ordered)
function getCurrentBadgeState() {
    const list = document.getElementById("badge-list");
    if (!list) return [];
    const badges = Array.from(list.querySelectorAll(".badge"));
    return badges.map((b) => ({
        id: b.dataset.id || "",
        hid: b.dataset.hid || "",
    }));
}

// NEW: shallow compare two badge states
function statesEqual(a, b) {
    if (!Array.isArray(a) || !Array.isArray(b)) return false;
    if (a.length !== b.length) return false;
    for (let i = 0; i < a.length; i++) {
        if (a[i].id !== b[i].id || a[i].hid !== b[i].hid) return false;
    }
    return true;
}

// NEW: create a badge from an id/hid snapshot
function createBadgeFromIdOrHid(item) {
    const { id, hid } = item || {};
    let el = null;
    // Prefer exact id
    if (id) el = document.getElementById(id);
    // Else try the hierarchical container id
    if (!el && hid) {
        el = document.getElementById(hid);
        if (!el) {
            // fallback: any descendant starting with hid
            const svg = getSvg();
            if (svg) el = svg.querySelector(`[id^="${CSS.escape(hid)}"]`);
        }
    }
    if (el) {
        const badge = createBadgeForShape(el);
        if (id) badge.dataset.id = id;
        if (hid) badge.dataset.hid = hid;
        return badge;
    }
    // Fallback minimal badge if element not found (e.g., filtered out)
    const span = document.createElement("span");
    span.className = "badge";
    if (id) span.dataset.id = id;
    if (hid) span.dataset.hid = hid;
    const label = document.createElement("span");
    label.textContent = getCaptionForId(hid || id || "item");
    span.appendChild(label);
    return span;
}

// NEW: apply a badge state and reload SVG
function applyBadgeState(state) {
    const list = document.getElementById("badge-list");
    if (!list) return;
    list.innerHTML = "";
    state.forEach((item) => {
        const badge = createBadgeFromIdOrHid(item);
        list.appendChild(badge);
    });
    // NEW: fit after re-creating badges
    requestAnimationFrame(refitAllBadges);
    reloadSvgFromBadges();
}

// NEW: perform undo (Ctrl/Cmd+Z)
function performUndo() {
    if (!undoStack.length) return;
    const prev = undoStack.pop();
    applyBadgeState(prev);
}

function deselectElementByHid(hid) {
    if (!hid) return;
    const el = document.getElementById(hid);
    if (!el) return;
    if (el.getAttribute("data-selected") === "true") {
        el.setAttribute("data-selected", "false");
        el.setAttribute(
            "stroke-width",
            el.getAttribute("data-original-stroke-width") || "3"
        );
        el.removeAttribute("filter");
    }
}

// NEW: refit all badges currently in the collector
function refitAllBadges() {
    const list = document.getElementById("badge-list");
    if (!list) return;
    const badges = list.querySelectorAll(".badge");
    badges.forEach((b) => fitBadgeLabel(b));
}

// NEW: truncate badge label from the left so it fits the available width
function fitBadgeLabel(badge) {
    if (!badge) return;
    const label = badge.querySelector("span:first-child");
    if (!label) return;

    const full = badge.dataset.fullLabel || label.textContent || "";
    badge.dataset.fullLabel = full;
    label.textContent = full;

    // If not measurable (e.g., hidden), skip
    const cw = label.clientWidth;
    if (!cw) return;

    // Fits already
    if (label.scrollWidth <= cw) return;

    const sep = " > ";
    const parts = full
        .split(sep)
        .map((s) => s.trim())
        .filter(Boolean);

    // Remove left-most breadcrumb segments first
    let cutSegments = 0;
    while (cutSegments < parts.length - 1) {
        const txt = "… > " + parts.slice(cutSegments + 1).join(sep);
        label.textContent = txt;
        if (label.scrollWidth <= cw) return;
        cutSegments++;
    }

    // Only last segment remains; chop characters from the left using binary search
    const last = parts.length ? parts[parts.length - 1] : full;
    let lo = 0,
        hi = last.length,
        best = "";
    while (lo <= hi) {
        const mid = Math.floor((lo + hi) / 2);
        const candidate = "… " + last.slice(mid);
        label.textContent = candidate;
        if (label.scrollWidth <= cw) {
            best = last.slice(mid);
            hi = mid - 1; // try to remove fewer chars
        } else {
            lo = mid + 1; // need to remove more
        }
    }
    label.textContent = best ? "… " + best : "…";
}

// NEW: reload SVG using current badges
function reloadSvgFromBadges(forceAllExpanded = false) {
    return reloadSvgFromBadgesImpl(forceAllExpanded);
}

async function reloadSvgFromBadgesImpl(forceAllExpanded) {
    try {
        if (!getActiveCreateSvgFunction()) return;
        const canvas = document.getElementById("canvas");
        if (!canvas) return;

        // NEW: preserve current zoom/pan state before reload
        const preservedState = { ...state };

        if (window.inputLoaded && typeof window.inputLoaded.then === "function") {
            await window.inputLoaded;
        }

        let expandedIds = [];
        let maxDepth = window.defaultDepth;
        if (forceAllExpanded) {
            // Instead of collecting all box IDs, set maxDepth to 100 to expand all
            maxDepth = 100;
        }
        else {
            // Use badges in the collector
            const list = document.getElementById("badge-list");
            if (list) {
                expandedIds = Array.from(list.querySelectorAll(".badge")).map(b => b.dataset.hid).filter(Boolean);
            }
        }
        // Extract ids from blacklisted badges (elements in blacklist)
        const blacklistIds = getAllBadgeCaptions("blacklist-list");
        // Use input YAML filename or fallback
        const arg = (typeof window.input === "string" && window.input.length > 0) ? window.input : "1";
        console.log(
            "Refreshing SVG: ",
            expandedIds,
            "blacklist ids: ",
            blacklistIds,
            "comments hidden: ",
            window.hideCommentsEnabled
        );
        let svgStr = await callCreateSvgWithMode(
            arg,
            expandedIds,
            blacklistIds,
            maxDepth
        );
        svgStr = svgStr && typeof svgStr.then === "function" ? await svgStr : svgStr;
        if (typeof svgStr !== "string" || !svgStr.trim().startsWith("<svg")) return;

        canvas.innerHTML = svgStr;

        // NEW: restore zoom/pan state after DOM update but before event
        state = preservedState;

        const evtSwap = new Event("htmx:afterSwap", { bubbles: true });
        canvas.dispatchEvent(evtSwap);
    } catch (e) {
        console.error("Error updating SVG via createSvg:", e);
    }
}
window.reloadSvgFromBadges = reloadSvgFromBadges;

// Helper: create a badge element from a shape
function createBadgeForShape(el) {
    // Build breadcrumb label from clicked and its parents
    const breadcrumb = buildBreadcrumbForId(el.id || "unnamed");

    const badge = document.createElement("span");
    badge.className = "badge";

    // Store stable identifier for toggle behavior
    const boxId = getBoxPrefix(el.id || "");
    if (boxId) badge.dataset.hid = boxId;

    // NEW: store full clicked element id for filtering
    if (el.id) badge.dataset.id = el.id;

    // Label
    const label = document.createElement("span");
    label.textContent = breadcrumb;

    // NEW: show full breadcrumb on hover
    badge.title = breadcrumb;

    // Compose
    badge.appendChild(label);

    // Take color from the clicked object's effective fill
    el =
        document.getElementById(boxId) ||
        document.querySelector(`[data-hid='${boxId}']`);

    const fillColor = getEffectiveFill(el);
    if (fillColor) {
        badge.dataset.fill = fillColor;
        badge.style.background = fillColor;
        badge.style.color = pickTextColor(fillColor);
    } else {
        // ...existing fallback styling via CSS .badge...
    }

    // Preserve any stroke info for future use
    const stroke = el.getAttribute("stroke");
    if (stroke) badge.dataset.stroke = stroke;

    return badge;
}

// NEW: resolve the effective fill color of an SVG element (attribute, style, or inherited/computed)
function getEffectiveFill(el) {
    if (!el || !(el instanceof Element)) return null;

    // Direct attribute first
    let fill = el.getAttribute && el.getAttribute("fill");
    if (fill && fill.toLowerCase() !== "none") return fill;

    // Computed style on this element
    try {
        const cs = getComputedStyle(el);
        if (cs && cs.fill && cs.fill !== "none") return cs.fill;
    } catch {
        /* ignore */
    }

    // Walk up to find inherited/computed non-none fill
    let cur = el.parentElement;
    while (cur && cur.tagName && cur.tagName.toLowerCase() !== "svg") {
        try {
            const cs = getComputedStyle(cur);
            if (cs && cs.fill && cs.fill !== "none") return cs.fill;
        } catch {
            /* ignore */
        }
        const attrFill = cur.getAttribute && cur.getAttribute("fill");
        if (attrFill && attrFill.toLowerCase() !== "none") return attrFill;
        cur = cur.parentElement;
    }

    return null;
}

// NEW: pick contrasting text color based on background
function pickTextColor(bg) {
    if (!bg) return "#000";

    // Normalize the color to rgb(r,g,b)
    let rgb = bg;
    if (!/^rgba?\(/i.test(bg)) {
        const tmp = document.createElement("span");
        tmp.style.color = bg;
        document.body.appendChild(tmp);
        const resolved = getComputedStyle(tmp).color;
        document.body.removeChild(tmp);
        rgb = resolved || bg;
    }
    const m = rgb.match(/rgba?\((\d+),\s*(\d+),\s*(\d+)/i);
    if (!m) return "#000";

    const r = parseInt(m[1], 10);
    const g = parseInt(m[2], 10);
    const b = parseInt(m[3], 10);
    // Perceived brightness
    const brightness = (r * 299 + g * 587 + b * 114) / 1000;
    return brightness > 140 ? "#000" : "#fff";
}

function getStrokeWidth(el) {
    if (!el || !(el instanceof Element)) return 1;
    const attr = parseFloat(el.getAttribute && el.getAttribute("stroke-width"));
    if (Number.isFinite(attr) && attr > 0) return attr;
    try {
        const cs = getComputedStyle(el);
        const parsed = parseFloat(cs && cs.strokeWidth);
        if (Number.isFinite(parsed) && parsed > 0) return parsed;
    } catch {}
    return 1;
}

function getConnectionPathData(conn) {
    if (!conn || !conn.tagName) return "";
    const tag = conn.tagName.toLowerCase();
    if (tag === "path") {
        return conn.getAttribute("d") || "";
    }
    if (tag === "line") {
        const x1 = conn.getAttribute("x1") || "0";
        const y1 = conn.getAttribute("y1") || "0";
        const x2 = conn.getAttribute("x2") || "0";
        const y2 = conn.getAttribute("y2") || "0";
        return `M ${x1} ${y1} L ${x2} ${y2}`;
    }
    if (tag === "polyline") {
        return polylinePointsToPath(conn.getAttribute("points"));
    }
    return "";
}

function polylinePointsToPath(pointsAttr) {
    if (!pointsAttr) return "";
    const tokens = pointsAttr
        .trim()
        .split(/[\s,]+/)
        .filter(Boolean);
    if (tokens.length < 4) return "";
    const pairs = [];
    for (let i = 0; i < tokens.length - 1; i += 2) {
        const x = tokens[i];
        const y = tokens[i + 1];
        if (typeof y === "undefined") break;
        pairs.push([x, y]);
    }
    if (!pairs.length) return "";
    let d = `M ${pairs[0][0]} ${pairs[0][1]}`;
    for (let i = 1; i < pairs.length; i++) {
        d += ` L ${pairs[i][0]} ${pairs[i][1]}`;
    }
    return d;
}

function getConnectionAnimationDuration(conn) {
    const length = getConnectionLength(conn);
    if (!Number.isFinite(length) || length <= 0) return 1.2;
    const speed = 140; // pixels per second (slower for calmer motion)
    const minDuration = 0.8;
    const maxDuration = 5.5;
    const raw = length / speed;
    return Math.min(Math.max(raw, minDuration), maxDuration);
}

function getConnectionLength(conn) {
    if (!conn) return NaN;
    if (typeof conn.getTotalLength === "function") {
        try {
            const len = conn.getTotalLength();
            if (Number.isFinite(len)) return len;
        } catch {}
    }
    if (conn.tagName && conn.tagName.toLowerCase() === "line") {
        const x1 = parseFloat(conn.getAttribute("x1"));
        const y1 = parseFloat(conn.getAttribute("y1"));
        const x2 = parseFloat(conn.getAttribute("x2"));
        const y2 = parseFloat(conn.getAttribute("y2"));
        if ([x1, y1, x2, y2].every(Number.isFinite)) {
            return Math.hypot(x2 - x1, y2 - y1);
        }
    }
    return NaN;
}

function getStrokeColor(el) {
    if (!el || !(el instanceof Element)) return null;
    const attr = el.getAttribute && el.getAttribute("stroke");
    if (attr && attr.toLowerCase() !== "none") return attr;
    try {
        const cs = getComputedStyle(el);
        if (
            cs &&
            cs.stroke &&
            cs.stroke !== "none" &&
            cs.stroke !== "rgba(0, 0, 0, 0)"
        ) {
            return cs.stroke;
        }
    } catch {
        /* ignore */
    }
    return null;
}

function clearConnectionRunners(targetSvg) {
    const svg = targetSvg || getSvg();
    if (!svg) return;
    svg.querySelectorAll(".connection-runner-dot").forEach((runner) => runner.remove());
}

function installConnectionRunners() {
    const svg = getSvg();
    if (!svg) return;

    clearConnectionRunners(svg);

    if (!connectionAnimationEnabled) return;

    const reduceMotion =
        typeof window !== "undefined" &&
        typeof window.matchMedia === "function" &&
        window.matchMedia("(prefers-reduced-motion: reduce)").matches;
    if (reduceMotion) {
        return;
    }

    const connections = svg.querySelectorAll(
        "line.connection, path.connection, polyline.connection"
    );
    connections.forEach((conn) => {
        if (!conn.parentNode) return;
        const pathData = getConnectionPathData(conn);
        if (!pathData) return;

        const runner = document.createElementNS(SVG_NS, "circle");
        runner.classList.add("connection-runner-dot");
        (conn.classList ? Array.from(conn.classList) : []).forEach((cls) => {
            if (cls) runner.classList.add(cls);
        });

        const stroke = getStrokeColor(conn) || "#0a84ff";
        runner.setAttribute("fill", stroke);
        runner.setAttribute("stroke", "none");

        const radius = Math.max(getStrokeWidth(conn), 1);
        runner.setAttribute("r", radius);
        runner.setAttribute("opacity", "0.95");
        runner.setAttribute("pointer-events", "none");

        if (conn.id) runner.dataset.connectionId = conn.id;

        const duration = getConnectionAnimationDuration(conn);
        runner.dataset.runnerDuration = `${duration}`;

        const motion = document.createElementNS(SVG_NS, "animateMotion");
        motion.setAttribute("dur", `${duration.toFixed(2)}s`);
        motion.setAttribute("repeatCount", "indefinite");
        motion.setAttribute("path", pathData);
        motion.setAttribute("rotate", "auto");
        if (Math.random() > 0.5) {
            motion.setAttribute("begin", `${(-Math.random() * duration).toFixed(2)}s`);
        }
        runner.appendChild(motion);

        const parent = conn.parentNode;
        if (!parent) return;
        const next = conn.nextSibling;
        if (next) {
            parent.insertBefore(runner, next);
        } else {
            parent.appendChild(runner);
        }
    });
}

function resetCommentLegendState() {
    commentLegendState.byNodeId.clear();
    commentLegendState.byGroupClass.clear();
    commentLegendState.itemActiveCount.forEach((_, el) => {
        if (el && el.classList) {
            el.classList.remove("comment-legend-item-active");
        }
    });
    commentLegendState.itemActiveCount.clear();
}

function registerCommentLegendEntry(entry, element) {
    if (!entry || !element) return;
    const sourceIds = Array.isArray(entry.sourceIds) ? entry.sourceIds : [];
    sourceIds.forEach((id) => {
        if (!id) return;
        commentLegendState.byNodeId.set(id, element);
    });
    const groups = Array.isArray(entry.groupClasses) ? entry.groupClasses : [];
    groups.forEach((cls) => {
        if (!cls) return;
        if (!commentLegendState.byGroupClass.has(cls)) {
            commentLegendState.byGroupClass.set(cls, new Set());
        }
        commentLegendState.byGroupClass.get(cls).add(element);
    });
}

function attachCommentLegendHover(element, entry) {
    if (!element || !entry) return;
    const groups = Array.isArray(entry.groupClasses)
        ? entry.groupClasses.filter(Boolean)
        : [];
    if (!groups.length) return;
    const handleEnter = () => {
        groups.forEach((grp) => {
            if (typeof window.highlightConnectionGroup === "function") {
                window.highlightConnectionGroup(grp);
            }
        });
    };
    const handleLeave = () => {
        groups.forEach((grp) => {
            if (typeof window.unhighlightConnectionGroup === "function") {
                window.unhighlightConnectionGroup(grp);
            }
        });
    };
    element.addEventListener("mouseenter", handleEnter);
    element.addEventListener("mouseleave", handleLeave);
}

function incrementLegendHighlight(el) {
    if (!el) return;
    const prev = commentLegendState.itemActiveCount.get(el) || 0;
    if (prev === 0) {
        el.classList.add("comment-legend-item-active");
    }
    commentLegendState.itemActiveCount.set(el, prev + 1);
}

function decrementLegendHighlight(el) {
    if (!el) return;
    const prev = commentLegendState.itemActiveCount.get(el) || 0;
    if (prev <= 1) {
        el.classList.remove("comment-legend-item-active");
        commentLegendState.itemActiveCount.delete(el);
    } else {
        commentLegendState.itemActiveCount.set(el, prev - 1);
    }
}

function highlightCommentLegendByGroupClass(groupClass) {
    if (!groupClass) return;
    const items = commentLegendState.byGroupClass.get(groupClass);
    if (!items || !items.size) return;
    items.forEach((el) => incrementLegendHighlight(el));
}

function unhighlightCommentLegendByGroupClass(groupClass) {
    if (!groupClass) return;
    const items = commentLegendState.byGroupClass.get(groupClass);
    if (!items || !items.size) return;
    items.forEach((el) => decrementLegendHighlight(el));
}

function hasCommentClass(el) {
    if (!el) return false;
    if (el.classList && el.classList.length) {
        for (const cls of el.classList) {
            if (String(cls || "").toLowerCase().includes("comment")) {
                return true;
            }
        }
    }
    const rawClass =
        (typeof el.className === "string"
            ? el.className
            : el.className && typeof el.className.baseVal === "string"
              ? el.className.baseVal
              : el.getAttribute && el.getAttribute("class")) || "";
    return rawClass
        .split(/\s+/)
        .filter(Boolean)
        .some((cls) => cls.toLowerCase().includes("comment"));
}

function isSvgTextElement(el) {
    return !!el && typeof el.tagName === "string" && el.tagName.toLowerCase() === "text";
}

function getTrimmedSvgText(el) {
    if (!el) return "";
    return (el.textContent || "").replace(/\s+/g, " ").trim();
}

function findNextLegendSibling(list, startIndex, predicate, processed, maxDistance) {
    if (!Array.isArray(list)) return null;
    let steps = 0;
    for (let i = startIndex; i < list.length; i++) {
        if (typeof maxDistance === "number" && maxDistance >= 0 && steps > maxDistance) {
            break;
        }
        const candidate = list[i];
        if (!candidate) continue;
        if (processed && processed.has(candidate)) continue;
        if (predicate(candidate)) {
            return { node: candidate, index: i };
        }
        steps++;
    }
    return null;
}

function resolveLegendTripletFromList(circle, list, processed, maxDistance) {
    if (!circle || !Array.isArray(list) || !list.length) return null;
    const idx = list.indexOf(circle);
    if (idx === -1) return null;
    const marker = findNextLegendSibling(
        list,
        idx + 1,
        (el) => isSvgTextElement(el) && hasCommentClass(el),
        processed,
        maxDistance
    );
    if (!marker) return null;
    const description = findNextLegendSibling(
        list,
        marker.index + 1,
        (el) => isSvgTextElement(el),
        processed,
        maxDistance
    );
    if (!description) return null;
    return { marker: marker.node, description: description.node };
}

// Locate sequential comment legend triplets inside the rendered SVG.
function collectCommentLegendEntries() {
    const svg = getSvg();
    if (!svg) return [];
    const entries = [];
    const processed = new WeakSet();
    const orderedNodes = Array.from(svg.querySelectorAll("circle, text"));
    const circles = svg.querySelectorAll("circle");
    circles.forEach((circle) => {
        if (!circle || processed.has(circle) || !hasCommentClass(circle)) return;
        const parentList = circle.parentElement
            ? Array.from(circle.parentElement.children || [])
            : [];
        const triplet =
            resolveLegendTripletFromList(circle, parentList, processed) ||
            resolveLegendTripletFromList(circle, orderedNodes, processed, 12);
        if (!triplet) return;
        processed.add(circle);
        processed.add(triplet.marker);
        processed.add(triplet.description);
        const strokeColor = circle.getAttribute("stroke");
        const fallbackColor =
            strokeColor && typeof strokeColor === "string" && strokeColor.toLowerCase() !== "none"
                ? strokeColor
                : "#6c7a89";
        const groupClassSet = new Set();
        [circle, triplet.marker, triplet.description].forEach((node) => {
            if (!node || !node.classList) return;
            node.classList.forEach((cls) => {
                if (/^conLine_\d+$/.test(cls)) {
                    groupClassSet.add(cls);
                }
            });
        });
        const groupClasses = Array.from(groupClassSet);
        const sourceIds = [circle.id, triplet.marker.id, triplet.description.id].filter(Boolean);
        entries.push({
            id:
                circle.id ||
                triplet.description.id ||
                triplet.marker.id ||
                `comment-${entries.length}`,
            markerLabel: getTrimmedSvgText(triplet.marker) || "Comment",
            body: getTrimmedSvgText(triplet.description) || "",
            color: getEffectiveFill(circle) || fallbackColor,
            sourceIds,
            groupClasses,
        });
    });
    return entries;
}

function createCommentLegendItem(entry) {
    const item = document.createElement("div");
    item.className = "comment-legend-item";
    item.dataset.commentLegendId = entry.id;
    const markerRow = document.createElement("div");
    markerRow.className = "comment-legend-marker";
    const label = document.createElement("span");
    label.className = "comment-legend-label";
    label.textContent = entry.markerLabel || "Comment";
    markerRow.appendChild(label);
    const body = document.createElement("div");
    body.className = "comment-legend-body";
    body.textContent = entry.body || "";
    item.appendChild(markerRow);
    item.appendChild(body);
    return item;
}

function updateCommentLegendPanel() {
    const list = document.getElementById("comment-legend-list");
    const panel = document.getElementById("comment-legend-panel");
    if (!list || !panel) return;
    resetCommentLegendState();
    const entries = collectCommentLegendEntries();
    list.innerHTML = "";
    if (!entries.length) {
        const empty = document.createElement("div");
        empty.className = "comment-legend-empty";
        empty.textContent = "No comment legend entries.";
        list.appendChild(empty);
        return;
    }
    entries.forEach((entry) => {
        const item = createCommentLegendItem(entry);
        list.appendChild(item);
        registerCommentLegendEntry(entry, item);
        attachCommentLegendHover(item, entry);
    });
}

// NEW: extract text content of the clicked box as a string array
function getTextContentArray(el) {
    const texts = [];
    try {
        const captionEl = document.getElementById(`${el.id}_capt`);
        if (captionEl && String(captionEl.tagName).toLowerCase() === "text") {
            const t = (captionEl.textContent || "").trim();
            if (t) {
                t.split(/\r?\n/)
                    .map((s) => s.trim())
                    .filter(Boolean)
                    .forEach((x) => texts.push(x));
            }
        }
    } catch {
        /* ignore */
    }

    if (texts.length === 0) {
        const tc = (el.textContent || "").trim();
        if (tc) {
            tc.split(/\r?\n/)
                .map((s) => s.trim())
                .filter(Boolean)
                .forEach((x) => texts.push(x));
        }
    }
    if (texts.length === 0 && el.id) {
        texts.push(el.id);
    }
    return texts;
}

// NEW: observe caption/text changes and trigger a refresh
function observeCaptionAndRefresh(el) {
    const target = document.getElementById(`${el.id}_capt`) || el; // fallback to the element itself

    if (!target) return;

    const getCurrentTexts = () => getTextContentArray(el);

    // Initial snapshot
    let lastTexts = getCurrentTexts();
    // If already has text, nothing to wait for
    if (lastTexts.length > 0) return;

    const obs = new MutationObserver(async (mutations) => {
        const texts = getCurrentTexts();
        // Trigger when text becomes non-empty
        if (texts.length > 0) {
            obs.disconnect();
            try {
                if (typeof createSvg !== "function") return;

                // Wait for YAML load if promise exists
                if (
                    window.inputLoaded &&
                    typeof window.inputLoaded.then === "function"
                ) {
                    await window.inputLoaded;
                }

                const canvas = document.getElementById("canvas");
                if (!canvas) return;

                const arg =
                    typeof window.input === "string" && window.input.length > 0
                        ? window.input
                        : "";

                const filterTexts = getAllBadgeCaptions("badge-list");
                // Extract ids from blacklisted badges (elements in blacklist)
                const blacklistIds = blacklist
                    .map((boxId) => {
                        // Try to get the element and its id
                        const el = document.getElementById(boxId);
                        return el ? el.id : boxId;
                    })
                    .filter(Boolean);
                console.log(
                    "Refreshing SVG: ",
                    filterTexts,
                    "blacklist ids: ",
                    blacklistIds,
                    "comments hidden: ",
                    window.hideCommentsEnabled
                );
                let svgStr = await callCreateSvgWithMode(
                    arg,
                    filterTexts,
                    blacklistIds
                );

                svgStr =
                    svgStr && typeof svgStr.then === "function"
                        ? await svgStr
                        : svgStr;
                if (
                    typeof svgStr !== "string" ||
                    !svgStr.trim().startsWith("<svg")
                )
                    return;

                canvas.innerHTML = svgStr;
                const evtSwap = new Event("htmx:afterSwap", {
                    bubbles: true,
                });
                canvas.dispatchEvent(evtSwap);
            } catch (e) {
                console.error("Error refreshing SVG after fill:", e);
            }
        }
    });

    obs.observe(target, {
        subtree: true,
        characterData: true,
        childList: true,
    });
}

// Helper: return caption (text) for a given hierarchical id via its companion <id>_capt <text> element
function getCaptionForId(hid) {
    try {
        const captionEl = document.getElementById(`${hid}_capt`);
        if (captionEl && String(captionEl.tagName).toLowerCase() === "text") {
            const text = (captionEl.textContent || "").trim();
            if (text) return text;
        }
    } catch {
        /* ignore */
    }
    return hid; // fallback to id
}

// Helper: build breadcrumb text for a shape id
function buildBreadcrumbForId(id) {
    const chain = getHierarchyIds(id);
    if (chain.length === 0) return getCaptionForId(id);
    const captions = chain.map(getCaptionForId);
    return captions.join(" > ");
}

// Helper: split an id into hierarchical ids from top to deepest
// Example: "id_1_2_3" -> ["id_1", "id_1_2", "id_1_2_3"]
function getHierarchyIds(id) {
    if (!id) return [];
    const parts = id.split("_");
    const isNumeric = (s) => /^\d+$/.test(s);

    // Find first numeric segment index
    let firstNumIdx = -1;
    for (let i = 0; i < parts.length; i++) {
        if (isNumeric(parts[i])) {
            firstNumIdx = i;
            break;
        }
    }
    if (firstNumIdx === -1) return []; // no hierarchy

    const basePrefix = parts.slice(0, firstNumIdx).join("_");
    const nums = parts.slice(firstNumIdx);
    // Only include numeric segments (ignore trailing non-numeric if any)
    const pureNums = nums.filter(isNumeric);

    const ids = [];
    for (let i = 0; i < pureNums.length && i < 5; i++) {
        // cap at five levels
        const suffix = pureNums.slice(0, i + 1).join("_");
        ids.push(basePrefix ? `${basePrefix}_${suffix}` : suffix);
    }
    return ids;
}

function getBoxPrefix(id) {
    const parts = id.split("_");
    const isNumeric = (s) => /^\d+$/.test(s);
    let firstNumIdx = parts.findIndex(isNumeric);
    if (firstNumIdx < 0) return id;
    const base = parts.slice(0, firstNumIdx).join("_");
    const nums = parts.slice(firstNumIdx).filter(isNumeric);
    // Build deepest id available as the box id
    const deep = nums.join("_");
    return base ? `${base}_${deep}` : deep;
}

// Helper: find existing badges by hierarchical box id
function findBadgesByBoxId(boxId) {
    const list = document.getElementById("badge-list");
    if (!list || !boxId) return [];
    // Use [data-hid] to uniquely identify badges per box
    return Array.from(
        list.querySelectorAll(`.badge[data-hid="${CSS.escape(boxId)}"]`)
    );
}

// NEW: check hierarchical relationship between hids (strict descendant)
function isDescendantHid(childHid, parentHid) {
    if (!childHid || !parentHid) return false;
    if (childHid === parentHid) return false;
    return childHid.startsWith(parentHid + "_");
}

// NEW: does the collector already contain a child of the given parent hid?
function anyBadgeIsChildOf(parentHid) {
    const list = document.getElementById("badge-list");
    if (!list || !parentHid) return false;
    const badges = list.querySelectorAll(".badge");
    for (const b of badges) {
        // Prefer dataset.hid; fallback to getBoxPrefix of dataset.id
        const hid =
            b.dataset.hid || (b.dataset.id ? getBoxPrefix(b.dataset.id) : "");
        if (hid && isDescendantHid(hid, parentHid)) return true;
    }
    return false;
}

// NEW: find badges that are ancestors of a given child hid
function findAncestorBadgesOf(childHid) {
    const list = document.getElementById("badge-list");
    if (!list || !childHid) return [];
    const badges = list.querySelectorAll(".badge");
    return Array.from(badges).filter((b) => {
        const hid =
            b.dataset.hid || (b.dataset.id ? getBoxPrefix(b.dataset.id) : "");
        return hid && isDescendantHid(childHid, hid); // badge.hid is an ancestor if child startsWith parent + "_"
    });
}

function setCommentLegendPanelVisibility(visible) {
    commentLegendPanelVisible = !!visible;
    const panel = document.getElementById("comment-legend-panel");
    if (!panel) return;
    const shouldHide = !commentLegendPanelVisible;
    panel.classList.toggle("hidden", shouldHide);
    panel.setAttribute("aria-hidden", shouldHide ? "true" : "false");
    if (!shouldHide) {
        updateCommentLegendPanel();
    }
    updateToolButtons();
}

function setCollectorPanelVisibility(visible) {
    collectorPanelVisible = !!visible;
    const box = document.getElementById("collector");
    if (!box) return;
    const shouldHide = !collectorPanelVisible;
    box.classList.toggle("hidden", shouldHide);
    box.setAttribute("aria-hidden", shouldHide ? "true" : "false");
    if (!shouldHide) {
        requestAnimationFrame(refitAllBadges);
    }
    positionBlacklistCollector();
    updateToolButtons();
}

function initCollectorVisibilityGuard() {
    if (collectorVisibilityGuardAttached) return;
    const box = document.getElementById("collector");
    if (!box) return;
    const observer = new MutationObserver(() => {
        if (!collectorPanelVisible && !box.classList.contains("hidden")) {
            box.classList.add("hidden");
            box.setAttribute("aria-hidden", "true");
            positionBlacklistCollector();
            updateToolButtons();
        }
    });
    observer.observe(box, { attributes: true, attributeFilter: ["class"] });
    collectorVisibilityGuardAttached = true;
}

// Pan/zoom via CSS transform on an HTML wrapper (#svg-stage) to avoid SVG viewport clipping.
// Ensure left-side panels (selected, expanded, blacklist) stack without overlapping
function positionBlacklistCollector() {
    const gap = 6;
    let currentTop = 56; // base below the toolbar
    const panels = [
        document.getElementById("selected-collector"),
        document.getElementById("collector"),
        document.getElementById("blacklist-collector"),
    ];
    panels.forEach((panel) => {
        if (!panel) return;
        if (panel.classList && panel.classList.contains("hidden")) {
            panel.style.top = "";
            return;
        }
        panel.style.top = currentTop + "px";
        const height = panel.offsetHeight || panel.getBoundingClientRect().height || 0;
        currentTop += (height || 0) + gap;
    });
}

// Drag to pan: start on mousedown if pan tool or space is active
function onStageMouseDown(e) {
    // Only left-button drag
    if (e.button !== 0) return;
    // Require pan tool or Space
    if (!(panToolActive || spacePressed)) return;
    ensureStageWrapped();
    const stage = getStage();
    if (!stage) return;

    isDragging = true;
    dragStart.x = e.clientX;
    dragStart.y = e.clientY;
    dragStart.tx = state.tx;
    dragStart.ty = state.ty;

    // Update cursor immediately
    applyTransform();

    // Prevent text selection during drag
    e.preventDefault();

    // Listen on window to keep drag even if cursor leaves stage
    window.addEventListener("mousemove", onWindowMouseMove);
    window.addEventListener("mouseup", onWindowMouseUp);
}

function onWindowMouseMove(e) {
    if (!isDragging) return;
    const dx = e.clientX - dragStart.x;
    const dy = e.clientY - dragStart.y;
    // Pan in screen pixels; translation is before scale, so apply directly
    state.tx = dragStart.tx + dx;
    state.ty = dragStart.ty + dy;
    applyTransform();
}

function onWindowMouseUp() {
    if (!isDragging) return;
    isDragging = false;
    window.removeEventListener("mousemove", onWindowMouseMove);
    window.removeEventListener("mouseup", onWindowMouseUp);
    applyTransform();
}

function ensureStageWrapped() {
    const canvas = getCanvas();
    if (!canvas) return;
    // If the SVG is directly under #canvas, wrap it in #svg-stage
    const existingStage = getStage();
    if (existingStage) return;
    const svg = canvas.querySelector("svg");
    if (!svg) return;
    const stage = document.createElement("div");
    stage.id = "svg-stage";
    canvas.innerHTML = "";
    stage.appendChild(svg);
    canvas.appendChild(stage);
}

function computeBaseSize() {
    const svg = getSvg();
    if (!svg) return;
    // from viewBox if present
    const vb = svg.viewBox && svg.viewBox.baseVal;
    if (vb && vb.width && vb.height) {
        baseSize.width = vb.width;
        baseSize.height = vb.height;
        return;
    }
    // fallback bbox/client size
    try {
        const bbox = svg.getBBox();
        baseSize.width = bbox.width || svg.clientWidth || 800;
        baseSize.height = bbox.height || svg.clientHeight || 600;
    } catch {
        baseSize.width = svg.clientWidth || 800;
        baseSize.height = svg.clientHeight || 600;
    }
}

// Centering helpers: only horizontal centering
function getCenterOffset() {
    const canvas = getCanvas();
    if (!canvas || !baseSize.width || !baseSize.height) return { cx: 0, cy: 0 };
    const cx = (canvas.clientWidth - baseSize.width * state.scale) / 2;
    const cy = 0; // vertical centering disabled
    return { cx, cy };
}
function getEffectiveTxTy() {
    const { cx, cy } = getCenterOffset();
    return { tx: cx + state.tx, ty: cy + state.ty };
}

function resizeStage() {
    // Keep scroll area equal to the visible canvas viewport, not scaling with zoom
    const stage = getStage();
    const canvas = getCanvas();
    if (!stage || !canvas) return;
    const w = canvas.clientWidth;
    const h = canvas.clientHeight;
    stage.style.width = w + "px";
    stage.style.height = h + "px";
}

function applyTransform() {
    ensureStageWrapped();
    const stage = getStage();
    if (!stage) return;
    // Use centered translation plus user pan offsets
    const { tx, ty } = getEffectiveTxTy();
    stage.style.transform = `translate(${tx}px, ${ty}px) scale(${state.scale})`;
    // Update cursor classes based on pan capability and drag
    stage.classList.toggle("pan-enabled", panToolActive || spacePressed);
    stage.classList.toggle("pan-dragging", isDragging);
    resizeStage();
    updateMinimap();
    updateToolButtons();
}

function initMinimap() {
    const mm = getMinimap();
    const scene = getMinimapScene();
    if (!mm || !scene || !baseSize.width || !baseSize.height) return;
    // Fit scene into minimap viewBox 0..100 with uniform scale
    mm.setAttribute("viewBox", `0 0 100 100`);
    // Compute scale to map scene -> 0..100 space while preserving aspect
    const sx = 100 / baseSize.width;
    const sy = 100 / baseSize.height;
    const s = Math.min(sx, sy);
    const sceneW = baseSize.width * s;
    const sceneH = baseSize.height * s;
    const offsetX = (100 - sceneW) / 2;
    const offsetY = (100 - sceneH) / 2;
    scene.setAttribute("x", String(offsetX));
    scene.setAttribute("y", String(offsetY));
    scene.setAttribute("width", String(sceneW));
    scene.setAttribute("height", String(sceneH));
    // Store mapping for quick updates
    mm._map = { s, offsetX, offsetY };
    updateMinimap(); // initial vp
    // reflect visibility state
    const mmWrap = document.getElementById("minimap");
    if (mmWrap) {
        mmWrap.style.display = minimapVisible ? "flex" : "none";
        mmWrap.setAttribute("aria-hidden", minimapVisible ? "false" : "true");
        populateMinimapPreview();
    }
}

function populateMinimapPreview() {
    const mm = getMinimap();
    const svg = getSvg();
    const content = getMinimapContent();
    if (!mm || !svg || !content || !baseSize.width || !baseSize.height) return;
    // Clear previous preview
    while (content.firstChild) content.removeChild(content.firstChild);

    // Clone visible scene (children of the root SVG). Avoid copying width/height/viewBox, styles are inline.
    const frag = document.createDocumentFragment();
    const children = svg.cloneNode(true).children; // shallow clone then use its children
    for (let i = 0; i < children.length; i++) {
        frag.appendChild(children[i].cloneNode(true));
    }
    content.appendChild(frag);

    // Scale and center cloned content to fit minimap scene using stored map
    const map = mm._map;
    if (!map) return;
    const tx = map.offsetX;
    const ty = map.offsetY;
    const s = map.s;
    content.setAttribute("transform", `translate(${tx},${ty}) scale(${s})`);
}

function updateMinimap() {
    const mmWrap = document.getElementById("minimap");
    if (!minimapVisible || !mmWrap || mmWrap.style.display === "none") return;
    const mm = getMinimap();
    const vp = getMinimapVp();
    const scene = getMinimapScene();
    const canvas = getCanvas();
    const svg = getSvg();
    if (
        !mm ||
        !vp ||
        !scene ||
        !canvas ||
        !svg ||
        !baseSize.width ||
        !baseSize.height
    )
        return;
    const map = mm._map;
    if (!map) return;

    // Visible viewport size in scene coordinates:
    // canvas shows a window of size (canvas.clientWidth, canvas.clientHeight) over the transformed stage.
    // Since we apply CSS scale to the stage, effective scene pixels per screen pixel = scale.
    const viewWScene = canvas.clientWidth / state.scale;
    const viewHScene = canvas.clientHeight / state.scale;

    // Top-left scene coordinate of current view considering pan and canvas scroll
    // Pan translates the stage by (tx, ty), positive moves content right/down => view moves left/up in scene.
    // Canvas scroll offsets move the view window within the stage.
    // Use effective (centered) translation for viewport origin
    const { tx: effTx, ty: effTy } = getEffectiveTxTy();
    const sx = (-effTx + canvas.scrollLeft) / state.scale;
    const sy = (-effTy + canvas.scrollTop) / state.scale;

    // Map to minimap coordinates
    const x = map.offsetX + sx * map.s;
    const y = map.offsetY + sy * map.s;
    const w = viewWScene * map.s;
    const h = viewHScene * map.s;

    // Clamp viewport inside scene rect
    const maxX =
        parseFloat(scene.getAttribute("x")) +
        parseFloat(scene.getAttribute("width")) -
        w;
    const maxY =
        parseFloat(scene.getAttribute("y")) +
        parseFloat(scene.getAttribute("height")) -
        h;
    vp.setAttribute(
        "x",
        String(Math.max(parseFloat(scene.getAttribute("x")), Math.min(x, maxX)))
    );
    vp.setAttribute(
        "y",
        String(Math.max(parseFloat(scene.getAttribute("y")), Math.min(y, maxY)))
    );
    vp.setAttribute("width", String(Math.max(1, w)));
    vp.setAttribute("height", String(Math.max(1, h)));
}

function getCanvas() {
    return document.getElementById("canvas");
}
function getStage() {
    return document.getElementById("svg-stage");
}
function getSvg() {
    return (
        document.querySelector("#svg-stage svg") ||
        document.querySelector("#canvas svg")
    );
}
function getMinimap() {
    return document.getElementById("minimap-svg");
}
function getMinimapScene() {
    return document.getElementById("minimap-scene");
}
function getMinimapVp() {
    return document.getElementById("minimap-vp");
}
function getMinimapContent() {
    return document.getElementById("minimap-content");
}

function updateToolButtons() {
    const btnPan = document.getElementById("btn-pan");
    const btnMinimap = document.getElementById("btn-minimap");
    const btnCollector = document.getElementById("btn-collector");
    const btnSelectedPanel = document.getElementById("btn-selected-panel");
    const btnBlacklist = document.getElementById("btn-blacklist");
    const btnHideComments = document.getElementById("btn-hide-comments");
    const btnCommentLegend = document.getElementById("btn-comment-legend");
    const btnConnectionAnim = document.getElementById("btn-connection-anim");
    if (btnPan)
        btnPan.classList.toggle("active", panToolActive || spacePressed);
    if (btnMinimap) btnMinimap.classList.toggle("active", minimapVisible);
    // Collector is active when visible (not hidden)
    const collector = document.getElementById("collector");
    const collectorVisible =
        collectorPanelVisible &&
        collector &&
        !collector.classList.contains("hidden");
    if (btnCollector)
        btnCollector.classList.toggle("active", !!collectorVisible);
    const selectedPanel = document.getElementById("selected-collector");
    const selectedVisible =
        selectedPanel && !selectedPanel.classList.contains("hidden");
    if (btnSelectedPanel)
        btnSelectedPanel.classList.toggle("active", !!selectedVisible);
    // Blacklist is active when blacklistMode is true
    if (btnBlacklist) btnBlacklist.classList.toggle("active", blacklistMode);
    // Reflect Hide Comments toggle
    if (btnHideComments) {
        const enabled = !!window.hideCommentsEnabled;
        btnHideComments.classList.toggle("active", enabled);
    }
    if (btnCommentLegend) {
        btnCommentLegend.classList.toggle("active", !!commentLegendPanelVisible);
    }
    if (btnConnectionAnim) {
        btnConnectionAnim.classList.toggle("active", !!connectionAnimationEnabled);
    }
}

// Blacklist tool toggle (moved to global svgControls below)

// Add to blacklist and update UI
function addToBlacklist(el) {
    if (!el || !el.id) return;
    const boxId = getBoxPrefix(el.id);
    // Always update both window.blacklist and local blacklist
    if (window.blacklist && Array.isArray(window.blacklist)) {
        if (window.blacklist.includes(boxId)) return;
        window.blacklist.push(boxId);
    }
    if (!blacklist.includes(boxId)) {
        blacklist.push(boxId);
    }
    updateBlacklistUI();
    // Reload SVG after adding to blacklist
    if (typeof reloadSvgFromBadges === "function") reloadSvgFromBadges();
}

// Remove from blacklist
function removeFromBlacklist(boxId) {
    blacklist = blacklist.filter((id) => id !== boxId);
    updateBlacklistUI();
    // Reload SVG after removing from blacklist
    if (typeof reloadSvgFromBadges === "function") reloadSvgFromBadges();
}

// Update blacklist collector UI
function updateBlacklistUI() {
    const list = document.getElementById("blacklist-list");
    if (!list) return;
    list.innerHTML = "";
    // Use window.blacklist if set, otherwise fallback to local blacklist
    const ids = (window.blacklist && Array.isArray(window.blacklist)) ? window.blacklist : blacklist;
    ids.forEach((boxId) => {
        // Try to find the SVG element for color extraction
        const el =
            document.getElementById(boxId) ||
            document.querySelector(`[data-hid='${boxId}']`);
        const breadcrumb = buildBreadcrumbForId(boxId);
        const badge = document.createElement("span");
        badge.className = "badge blacklist-badge";
        badge.dataset.hid = boxId;
        // Label span for truncation, etc.
        const label = document.createElement("span");
        label.textContent = breadcrumb;
        badge.appendChild(label);
        badge.title = breadcrumb;
        // Color from shape if possible
        if (el) {
            const fillColor = getEffectiveFill(el);
            if (fillColor) {
                badge.dataset.fill = fillColor;
                badge.style.background = fillColor;
                badge.style.color = pickTextColor(fillColor);
            }
        }
        badge.onclick = function () {
            // Remove from window.blacklist if present, else from local blacklist
            if (window.blacklist && Array.isArray(window.blacklist)) {
                window.blacklist = window.blacklist.filter(id => id !== boxId);
            } else {
                blacklist = blacklist.filter(id => id !== boxId);
            }
            updateBlacklistUI();
            // Also reload the SVG to reflect the change
            if (typeof reloadSvgFromBadges === "function") reloadSvgFromBadges();
        };
        list.appendChild(badge);
    });
}

// Space key enables temporary pan + NEW: Ctrl/Cmd+Z undo
function onKeyDown(e) {
    // NEW: Ctrl/Cmd+Z -> undo last click action
    if (
        (e.ctrlKey || e.metaKey) &&
        !e.shiftKey &&
        (e.key === "z" || e.code === "KeyZ")
    ) {
        performUndo();
        e.preventDefault();
        return;
    }
    if (e.code === "Space" && !spacePressed) {
        spacePressed = true;
        applyTransform();
        // Prevent page scroll when space is pressed
        e.preventDefault();
    }
}
function onKeyUp(e) {
    if (e.code === "Space") {
        spacePressed = false;
        // If releasing Space while dragging, end drag
        if (isDragging) onWindowMouseUp();
        applyTransform();
    }
}

// Stub: invoked when the Hide Comments toolbar toggle changes state.
// Later, this can call into WASM or update rendering as needed.
window.onHideCommentsChanged = function (enabled) {
    try {
        console.log("onHideCommentsChanged:", enabled);
        // Keep toolbar state in sync
        updateToolButtons();
        // Re-render SVG with updated showComments flag
        if (typeof window.reloadSvgFromBadges === "function") {
            window.reloadSvgFromBadges();
        }
    } catch (err) {
        // no-op
    }
};
