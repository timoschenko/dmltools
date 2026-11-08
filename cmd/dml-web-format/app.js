function updateEditor(
  textarea,
  highlightEl,
  gutterEl,
  countEl,
  guideEl,
  showWhitespace = false,
  errorInfo = null,
  showLineLength = false,
  wrapLimit = 80,
  tokens = null,
) {
  const text = textarea.value;
  highlightEl.firstElementChild.innerHTML = highlightCode(
    text,
    showWhitespace,
    errorInfo,
    tokens,
  );

  const lines = text.length === 0 ? [""] : text.split("\n");
  const lineCount = lines.length;

  let maxLen = 0;
  if (showLineLength) {
    for (let i = 0; i < lineCount; i++) {
      if (lines[i].length > maxLen) {
        maxLen = lines[i].length;
      }
    }
  }

  const maxNumDigits = Math.max(1, lineCount.toString().length);
  const maxLenDigits = Math.max(2, maxLen.toString().length);

  let nums = "";
  for (let i = 1; i <= lineCount; i++) {
    const lineText = lines[i - 1];
    const numStr = i.toString().padStart(maxNumDigits, " ");

    let lineContent = `<span class="gutter-num">${numStr}</span>`;
    if (showLineLength) {
      const len = lineText.length;
      const lenStr = len.toString().padStart(maxLenDigits, " ");
      const isOverflow = len > wrapLimit;
      const title = isOverflow
        ? `Line ${i}: ${len} characters (exceeds wrap limit ${wrapLimit})`
        : `Line ${i}: ${len} characters`;
      lineContent += `<span class="gutter-sep"> │ </span><span class="gutter-len${isOverflow ? " len-overflow" : ""}" title="${title}">${lenStr}</span>`;
    }

    if (errorInfo && errorInfo.line === i) {
      nums += `<span class="gutter-error" title="Syntax error at line ${i}">${lineContent}</span>\n`;
    } else {
      nums += lineContent + "\n";
    }
  }
  gutterEl.innerHTML = nums;

  if (countEl) {
    const count = text.length === 0 ? 0 : lineCount;
    countEl.textContent = `${count} line${count === 1 ? "" : "s"}`;
  }

  syncScroll(textarea, highlightEl, gutterEl, guideEl);
}

function syncScroll(textarea, highlightEl, gutterEl, guideEl) {
  highlightEl.scrollTop = textarea.scrollTop;
  highlightEl.scrollLeft = textarea.scrollLeft;
  gutterEl.scrollTop = textarea.scrollTop;
  if (guideEl) {
    guideEl.style.transform = `translateX(-${textarea.scrollLeft}px)`;
  }
}

const inputArea = document.getElementById("input");
const inputHighlight = document.getElementById("inputHighlight");
const inputLineNumbers = document.getElementById("inputLineNumbers");
const inputLineCount = document.getElementById("inputLineCount");
const inputCursorInfo = document.getElementById("inputCursorInfo");
const inputWrapGuide = document.getElementById("inputWrapGuide");

const outputArea = document.getElementById("output");
const outputHighlight = document.getElementById("outputHighlight");
const outputLineNumbers = document.getElementById("outputLineNumbers");
const outputLineCount = document.getElementById("outputLineCount");
const outputCursorInfo = document.getElementById("outputCursorInfo");
const outputWrapGuide = document.getElementById("outputWrapGuide");
const outputWhitespaceToggle = document.getElementById(
  "outputWhitespaceToggle",
);
const outputLineLengthToggle = document.getElementById(
  "outputLineLengthToggle",
);
const lineWrapInput = document.getElementById("lineWrapInput");

const formatBtn = document.getElementById("formatBtn");
const statusDiv = document.getElementById("status");

const downloadTraceBtn = document.getElementById("downloadTraceBtn");
const openPerfettoBtn = document.getElementById("openPerfettoBtn");
const openChromeTracingBtn = document.getElementById("openChromeTracingBtn");

const chromeTracingModal = document.getElementById("chromeTracingModal");
const closeChromeTracingModal = document.getElementById(
  "closeChromeTracingModal",
);
const dismissChromeTracingModal = document.getElementById(
  "dismissChromeTracingModal",
);
const copyChromeTracingUrlBtn = document.getElementById(
  "copyChromeTracingUrlBtn",
);
const chromeTracingUrlInput = document.getElementById("chromeTracingUrlInput");
const modalDownloadTraceBtn = document.getElementById("modalDownloadTraceBtn");

const STORAGE_KEY_INPUT = "dml_last_input";
const STORAGE_KEY_OUTPUT = "dml_last_output";
const STORAGE_KEY_WHITESPACE = "dml_show_whitespace";
const STORAGE_KEY_LINE_LENGTH = "dml_show_line_length";
const STORAGE_KEY_LINE_WRAP = "dml_line_wrap";
const STORAGE_KEY_INPUT_COLLAPSED = "dml_input_panel_collapsed";
const STORAGE_KEY_STATS_COLLAPSED = "dml_stats_panel_collapsed";
const STORAGE_KEY_INPUT_WIDTH = "dml_input_panel_width";
const STORAGE_KEY_STATS_WIDTH = "dml_stats_panel_width";

const inputPanel = document.getElementById("inputPanel");
const outputPanel = document.getElementById("outputPanel");
const statsPanel = document.getElementById("statsPanel");

const resizerLeft = document.getElementById("resizerLeft");
const resizerRight = document.getElementById("resizerRight");

const collapseInputBtn = document.getElementById("collapseInputBtn");
const expandInputBtn = document.getElementById("expandInputBtn");
const inputPanelRail = document.getElementById("inputPanelRail");

const collapseStatsBtn = document.getElementById("collapseStatsBtn");
const expandStatsBtn = document.getElementById("expandStatsBtn");
const statsPanelRail = document.getElementById("statsPanelRail");

const headerToggleInputBtn = document.getElementById("headerToggleInputBtn");
const headerToggleStatsBtn = document.getElementById("headerToggleStatsBtn");

const statsTokenCount = document.getElementById("statsTokenCount");
const statsTraceEvents = document.getElementById("statsTraceEvents");
const statsLinesFormatted = document.getElementById("statsLinesFormatted");
const statsDuration = document.getElementById("statsDuration");

let currentError = null;
let currentTrace = null;
let currentTraceCount = 0;
let currentTokens = null;
let tokenTypeToName = new Map();
let parsedTraceCache = null;
let winLossStatsCache = null;

function updateTraceButtons(traceJson, traceCount = 0) {
  parsedTraceCache = null;
  winLossStatsCache = null;

  currentTrace = traceJson;
  currentTraceCount = traceCount;

  const hasTrace = typeof traceJson === "string" && traceJson.length > 0;
  if (downloadTraceBtn) {
    downloadTraceBtn.disabled = !hasTrace;
    downloadTraceBtn.title = hasTrace
      ? `Download collected trace (${traceCount} events)`
      : "Download collected trace as JSON (Format DML first)";
  }
  if (openPerfettoBtn) {
    openPerfettoBtn.disabled = !hasTrace;
    openPerfettoBtn.title = hasTrace
      ? `Open trace in ui.perfetto.dev (${traceCount} events)`
      : "Open trace directly in ui.perfetto.dev (Format DML first)";
  }
  if (openChromeTracingBtn) {
    openChromeTracingBtn.disabled = !hasTrace;
    openChromeTracingBtn.title = hasTrace
      ? `Open chrome://tracing guidance and load trace (${traceCount} events)`
      : "Open chrome://tracing guidance and load trace (Format DML first)";
  }
}

function downloadTraceFile(traceJson, fileName = "dml_trace.json") {
  if (!traceJson) return;
  const blob = new Blob([traceJson], { type: "application/json" });
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = fileName;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  setTimeout(() => URL.revokeObjectURL(url), 1000);
}

function openInPerfetto(traceJson) {
  if (!traceJson) return;

  const perfettoUrl = "https://ui.perfetto.dev";
  const win = window.open(perfettoUrl, "_blank");
  if (!win) {
    alert(
      "Pop-up blocked. Please allow pop-ups for this site to open ui.perfetto.dev.",
    );
    return;
  }

  const labelSpan = openPerfettoBtn
    ? openPerfettoBtn.querySelector("span")
    : null;
  const origBtnText = labelSpan ? labelSpan.textContent : null;
  if (openPerfettoBtn && labelSpan) {
    labelSpan.textContent = "Opening...";
    openPerfettoBtn.disabled = true;
  }

  const resetBtn = () => {
    if (openPerfettoBtn && labelSpan && origBtnText) {
      labelSpan.textContent = origBtnText;
      openPerfettoBtn.disabled = false;
    }
  };

  const buffer = new TextEncoder().encode(traceJson).buffer;

  // Handshake protocol as specified by Perfetto deep-linking docs
  const pingInterval = setInterval(() => {
    try {
      win.postMessage("PING", perfettoUrl);
    } catch (e) {
      clearInterval(pingInterval);
      resetBtn();
    }
  }, 50);

  const timeoutId = setTimeout(() => {
    clearInterval(pingInterval);
    window.removeEventListener("message", onMessage);
    resetBtn();
  }, 15000);

  function onMessage(evt) {
    if (evt.origin !== perfettoUrl) return;
    if (evt.data === "PONG") {
      clearInterval(pingInterval);
      clearTimeout(timeoutId);
      window.removeEventListener("message", onMessage);

      win.postMessage(
        {
          perfetto: {
            buffer: buffer,
            title: "DML Formatter Trace",
            fileName: "dml_trace.json",
          },
        },
        perfettoUrl,
      );
      resetBtn();
    }
  }

  window.addEventListener("message", onMessage);
}

function showChromeTracingModal() {
  if (chromeTracingModal) {
    chromeTracingModal.classList.remove("hidden");
  }
}

function hideChromeTracingModal() {
  if (chromeTracingModal) {
    chromeTracingModal.classList.add("hidden");
  }
}

function saveToLocalStorage(key, val) {
  try {
    localStorage.setItem(key, val);
  } catch (e) {
    console.warn("Unable to save to localStorage:", e);
  }
}

function loadFromLocalStorage(key) {
  try {
    return localStorage.getItem(key);
  } catch (e) {
    console.warn("Unable to read from localStorage:", e);
    return null;
  }
}

function getLineWrap() {
  const val = parseInt(lineWrapInput.value, 10);
  return isNaN(val) || val < 1 ? 80 : val;
}

function updateLineWrap() {
  const col = getLineWrap();
  document.documentElement.style.setProperty("--wrap-guide-column", col);
  saveToLocalStorage(STORAGE_KEY_LINE_WRAP, col.toString());
  updateOutputView();
}

function handleInputChange() {
  saveToLocalStorage(STORAGE_KEY_INPUT, inputArea.value);
  updateEditor(
    inputArea,
    inputHighlight,
    inputLineNumbers,
    inputLineCount,
    inputWrapGuide,
    false,
    null,
    false,
  );
  updateInputCursorInfo();
}

function updateEditorCursorInfo(textarea, infoEl) {
  if (!infoEl || !textarea) return;
  const text = textarea.value;
  const start = textarea.selectionStart || 0;
  const end = textarea.selectionEnd || 0;
  const selectedCount = Math.abs(end - start);

  const pos = textarea.selectionDirection === "backward" ? start : end;
  const textUpToCursor = text.slice(0, pos);
  const lines = textUpToCursor.split("\n");
  const line = lines.length;
  const col = lines[lines.length - 1].length + 1;

  if (selectedCount > 0) {
    infoEl.innerHTML = `Ln ${line}, Col ${col} <span class="selection-count">(${selectedCount} selected)</span>`;
  } else {
    infoEl.textContent = `Ln ${line}, Col ${col}`;
  }
}

function updateInputCursorInfo() {
  updateEditorCursorInfo(inputArea, inputCursorInfo);
}

function updateOutputCursorInfo() {
  updateEditorCursorInfo(outputArea, outputCursorInfo);
}

function updateOutputView(errorInfo = currentError, tokens = null) {
  const showLength = outputLineLengthToggle
    ? outputLineLengthToggle.checked
    : true;

  const isStatsCollapsed =
    statsPanel && statsPanel.classList.contains("collapsed");
  const tokensToRender = isStatsCollapsed ? null : tokens;

  if (outputHighlight) {
    outputHighlight.classList.remove("has-selection");
  }

  updateEditor(
    outputArea,
    outputHighlight,
    outputLineNumbers,
    outputLineCount,
    outputWrapGuide,
    outputWhitespaceToggle.checked,
    errorInfo,
    showLength,
    getLineWrap(),
    tokensToRender,
  );

  const outputWrapper = outputArea.closest(".editor-wrapper");
  if (outputWrapper) {
    if (errorInfo) {
      outputWrapper.classList.add("has-error");
    } else {
      outputWrapper.classList.remove("has-error");
    }
  }

  updateOutputCursorInfo();
}

inputArea.addEventListener("input", handleInputChange);
inputArea.addEventListener("scroll", () =>
  syncScroll(inputArea, inputHighlight, inputLineNumbers, inputWrapGuide),
);

inputArea.addEventListener("select", updateInputCursorInfo);
inputArea.addEventListener("keyup", updateInputCursorInfo);
inputArea.addEventListener("mouseup", updateInputCursorInfo);
inputArea.addEventListener("click", updateInputCursorInfo);
inputArea.addEventListener("focus", updateInputCursorInfo);

// Tab key support for input textarea
inputArea.addEventListener("keydown", (e) => {
  if (e.key === "Tab") {
    e.preventDefault();
    const start = inputArea.selectionStart;
    const end = inputArea.selectionEnd;
    inputArea.setRangeText("    ", start, end, "end");
    handleInputChange();
  }
});

outputArea.addEventListener("scroll", () =>
  syncScroll(outputArea, outputHighlight, outputLineNumbers, outputWrapGuide),
);

outputArea.addEventListener("select", updateOutputCursorInfo);
outputArea.addEventListener("keyup", updateOutputCursorInfo);
outputArea.addEventListener("mouseup", updateOutputCursorInfo);
outputArea.addEventListener("click", updateOutputCursorInfo);
outputArea.addEventListener("focus", updateOutputCursorInfo);
document.addEventListener("selectionchange", () => {
  if (document.activeElement === outputArea) {
    updateOutputCursorInfo();
  } else if (document.activeElement === inputArea) {
    updateInputCursorInfo();
  }
});

outputWhitespaceToggle.addEventListener("change", () => {
  saveToLocalStorage(
    STORAGE_KEY_WHITESPACE,
    outputWhitespaceToggle.checked ? "1" : "0",
  );
  updateOutputView(currentError, currentTokens);
});

if (outputLineLengthToggle) {
  outputLineLengthToggle.addEventListener("change", () => {
    saveToLocalStorage(
      STORAGE_KEY_LINE_LENGTH,
      outputLineLengthToggle.checked ? "1" : "0",
    );
    updateOutputView(currentError, currentTokens);
  });
}

lineWrapInput.addEventListener("input", updateLineWrap);
lineWrapInput.addEventListener("change", updateLineWrap);

// Restore from localStorage
const savedInput = loadFromLocalStorage(STORAGE_KEY_INPUT);
if (savedInput !== null) {
  inputArea.value = savedInput;
}

const savedOutput = loadFromLocalStorage(STORAGE_KEY_OUTPUT);
if (savedOutput !== null) {
  outputArea.value = savedOutput;
}

const savedWhitespace = loadFromLocalStorage(STORAGE_KEY_WHITESPACE);
if (savedWhitespace !== null) {
  outputWhitespaceToggle.checked = savedWhitespace === "1";
}

if (outputLineLengthToggle) {
  const savedLineLength = loadFromLocalStorage(STORAGE_KEY_LINE_LENGTH);
  if (savedLineLength !== null) {
    outputLineLengthToggle.checked = savedLineLength === "1";
  }
}

const savedLineWrap = loadFromLocalStorage(STORAGE_KEY_LINE_WRAP);
if (savedLineWrap !== null) {
  lineWrapInput.value = savedLineWrap;
}
updateLineWrap();

// Initial render
handleInputChange();
updateOutputView(null);

function setInputCollapsed(collapsed, save = true) {
  if (!inputPanel) return;
  if (collapsed) {
    inputPanel.classList.add("collapsed");
    inputPanel.style.width = "";
    inputPanel.style.flex = "";
    if (resizerLeft) resizerLeft.classList.add("hidden");
  } else {
    inputPanel.classList.remove("collapsed");
    if (inputPanel.dataset.customWidth) {
      inputPanel.style.width = `${inputPanel.dataset.customWidth}px`;
      inputPanel.style.flex = `0 0 ${inputPanel.dataset.customWidth}px`;
    }
    if (resizerLeft) resizerLeft.classList.remove("hidden");
  }
  if (headerToggleInputBtn) {
    headerToggleInputBtn.classList.toggle("active", !collapsed);
  }
  if (save) {
    saveToLocalStorage(STORAGE_KEY_INPUT_COLLAPSED, collapsed ? "1" : "0");
  }
  setTimeout(() => {
    syncScroll(inputArea, inputHighlight, inputLineNumbers, inputWrapGuide);
    syncScroll(outputArea, outputHighlight, outputLineNumbers, outputWrapGuide);
  }, 220);
}

function setStatsCollapsed(collapsed, save = true) {
  if (!statsPanel) return;
  if (collapsed) {
    statsPanel.classList.add("collapsed");
    statsPanel.style.width = "";
    statsPanel.style.flex = "";
    if (resizerRight) resizerRight.classList.add("hidden");
  } else {
    statsPanel.classList.remove("collapsed");
    if (statsPanel.dataset.customWidth) {
      statsPanel.style.width = `${statsPanel.dataset.customWidth}px`;
      statsPanel.style.flex = `0 0 ${statsPanel.dataset.customWidth}px`;
    }
    if (resizerRight) resizerRight.classList.remove("hidden");
  }
  if (headerToggleStatsBtn) {
    headerToggleStatsBtn.classList.toggle("active", !collapsed);
  }
  if (save) {
    saveToLocalStorage(STORAGE_KEY_STATS_COLLAPSED, collapsed ? "1" : "0");
  }
  updateOutputView(currentError, currentTokens);
  setTimeout(() => {
    syncScroll(outputArea, outputHighlight, outputLineNumbers, outputWrapGuide);
  }, 220);
}

function updateStatsPanel(tokenCount, traceCount, lineCount, durationMs) {
  if (statsTokenCount) {
    statsTokenCount.textContent =
      typeof tokenCount === "number" ? tokenCount.toLocaleString() : "-";
  }
  if (statsTraceEvents) {
    statsTraceEvents.textContent =
      typeof traceCount === "number" ? traceCount.toLocaleString() : "-";
  }
  if (statsLinesFormatted) {
    statsLinesFormatted.textContent =
      typeof lineCount === "number" ? lineCount.toLocaleString() : "-";
  }
  if (statsDuration) {
    statsDuration.textContent =
      typeof durationMs === "number" ? `${durationMs.toFixed(1)} ms` : "-";
  }
}

// Draggable Resizers Controller
function initResizers() {
  const container = document.querySelector(".container");
  if (!container) return;

  // Left resizer (Input Panel)
  if (resizerLeft && inputPanel && outputPanel) {
    let isDraggingLeft = false;
    let startX = 0;
    let startWidth = 0;

    resizerLeft.addEventListener("pointerdown", (e) => {
      if (e.button !== 0) return;
      if (inputPanel.classList.contains("collapsed")) return;
      isDraggingLeft = true;
      startX = e.clientX;
      startWidth = inputPanel.getBoundingClientRect().width;
      resizerLeft.setPointerCapture(e.pointerId);
      resizerLeft.classList.add("resizing");
      document.body.classList.add("is-resizing");
      e.preventDefault();
    });

    resizerLeft.addEventListener("pointermove", (e) => {
      if (!isDraggingLeft) return;
      const deltaX = e.clientX - startX;
      const containerWidth = container.getBoundingClientRect().width;
      const statsWidth = statsPanel
        ? statsPanel.getBoundingClientRect().width
        : 300;
      const maxInputWidth = Math.max(200, containerWidth - statsWidth - 280);
      const newWidth = Math.max(
        200,
        Math.min(maxInputWidth, Math.round(startWidth + deltaX)),
      );

      inputPanel.style.width = `${newWidth}px`;
      inputPanel.style.flex = `0 0 ${newWidth}px`;
      inputPanel.dataset.customWidth = newWidth;

      syncScroll(inputArea, inputHighlight, inputLineNumbers, inputWrapGuide);
      syncScroll(
        outputArea,
        outputHighlight,
        outputLineNumbers,
        outputWrapGuide,
      );
    });

    const stopDragLeft = (e) => {
      if (!isDraggingLeft) return;
      isDraggingLeft = false;
      try {
        resizerLeft.releasePointerCapture(e.pointerId);
      } catch (_) {}
      resizerLeft.classList.remove("resizing");
      document.body.classList.remove("is-resizing");

      if (inputPanel.dataset.customWidth) {
        saveToLocalStorage(
          STORAGE_KEY_INPUT_WIDTH,
          inputPanel.dataset.customWidth,
        );
      }
      syncScroll(inputArea, inputHighlight, inputLineNumbers, inputWrapGuide);
      syncScroll(
        outputArea,
        outputHighlight,
        outputLineNumbers,
        outputWrapGuide,
      );
    };

    resizerLeft.addEventListener("pointerup", stopDragLeft);
    resizerLeft.addEventListener("pointercancel", stopDragLeft);

    // Double-click to reset to default
    resizerLeft.addEventListener("dblclick", () => {
      inputPanel.style.width = "";
      inputPanel.style.flex = "1 1 0";
      delete inputPanel.dataset.customWidth;
      try {
        localStorage.removeItem(STORAGE_KEY_INPUT_WIDTH);
      } catch (_) {}
      syncScroll(inputArea, inputHighlight, inputLineNumbers, inputWrapGuide);
      syncScroll(
        outputArea,
        outputHighlight,
        outputLineNumbers,
        outputWrapGuide,
      );
    });
  }

  // Right resizer (Stats Panel)
  if (resizerRight && statsPanel && outputPanel) {
    let isDraggingRight = false;
    let startX = 0;
    let startWidth = 0;

    resizerRight.addEventListener("pointerdown", (e) => {
      if (e.button !== 0) return;
      if (statsPanel.classList.contains("collapsed")) return;
      isDraggingRight = true;
      startX = e.clientX;
      startWidth = statsPanel.getBoundingClientRect().width;
      resizerRight.setPointerCapture(e.pointerId);
      resizerRight.classList.add("resizing");
      document.body.classList.add("is-resizing");
      e.preventDefault();
    });

    resizerRight.addEventListener("pointermove", (e) => {
      if (!isDraggingRight) return;
      // Moving left (negative deltaX) increases right panel width
      const deltaX = startX - e.clientX;
      const containerWidth = container.getBoundingClientRect().width;
      const inputWidth = inputPanel
        ? inputPanel.getBoundingClientRect().width
        : 300;
      const maxStatsWidth = Math.max(220, containerWidth - inputWidth - 280);
      const newWidth = Math.max(
        220,
        Math.min(Math.min(550, maxStatsWidth), Math.round(startWidth + deltaX)),
      );

      statsPanel.style.width = `${newWidth}px`;
      statsPanel.style.flex = `0 0 ${newWidth}px`;
      statsPanel.dataset.customWidth = newWidth;

      syncScroll(
        outputArea,
        outputHighlight,
        outputLineNumbers,
        outputWrapGuide,
      );
    });

    const stopDragRight = (e) => {
      if (!isDraggingRight) return;
      isDraggingRight = false;
      try {
        resizerRight.releasePointerCapture(e.pointerId);
      } catch (_) {}
      resizerRight.classList.remove("resizing");
      document.body.classList.remove("is-resizing");

      if (statsPanel.dataset.customWidth) {
        saveToLocalStorage(
          STORAGE_KEY_STATS_WIDTH,
          statsPanel.dataset.customWidth,
        );
      }
      syncScroll(
        outputArea,
        outputHighlight,
        outputLineNumbers,
        outputWrapGuide,
      );
    };

    resizerRight.addEventListener("pointerup", stopDragRight);
    resizerRight.addEventListener("pointercancel", stopDragRight);

    // Double-click to reset to default
    resizerRight.addEventListener("dblclick", () => {
      statsPanel.style.width = "320px";
      statsPanel.style.flex = "0 0 320px";
      delete statsPanel.dataset.customWidth;
      try {
        localStorage.removeItem(STORAGE_KEY_STATS_WIDTH);
      } catch (_) {}
      syncScroll(
        outputArea,
        outputHighlight,
        outputLineNumbers,
        outputWrapGuide,
      );
    });
  }
}

// Collapsible Panel Listeners
if (collapseInputBtn) {
  collapseInputBtn.addEventListener("click", () => setInputCollapsed(true));
}
if (expandInputBtn) {
  expandInputBtn.addEventListener("click", (e) => {
    e.stopPropagation();
    setInputCollapsed(false);
  });
}
if (inputPanelRail) {
  inputPanelRail.addEventListener("click", () => setInputCollapsed(false));
}
if (headerToggleInputBtn) {
  headerToggleInputBtn.addEventListener("click", () => {
    const isCollapsed =
      inputPanel && inputPanel.classList.contains("collapsed");
    setInputCollapsed(!isCollapsed);
  });
}

if (collapseStatsBtn) {
  collapseStatsBtn.addEventListener("click", () => setStatsCollapsed(true));
}
if (expandStatsBtn) {
  expandStatsBtn.addEventListener("click", (e) => {
    e.stopPropagation();
    setStatsCollapsed(false);
  });
}
if (statsPanelRail) {
  statsPanelRail.addEventListener("click", () => setStatsCollapsed(false));
}
if (headerToggleStatsBtn) {
  headerToggleStatsBtn.addEventListener("click", () => {
    const isCollapsed =
      statsPanel && statsPanel.classList.contains("collapsed");
    setStatsCollapsed(!isCollapsed);
  });
}

// Restore custom widths from localStorage
const savedInputWidth = loadFromLocalStorage(STORAGE_KEY_INPUT_WIDTH);
if (savedInputWidth && parseInt(savedInputWidth, 10) >= 200) {
  inputPanel.dataset.customWidth = savedInputWidth;
}

const savedStatsWidth = loadFromLocalStorage(STORAGE_KEY_STATS_WIDTH);
if (savedStatsWidth && parseInt(savedStatsWidth, 10) >= 220) {
  statsPanel.dataset.customWidth = savedStatsWidth;
}

// Restore panel collapsed states
const savedInputCollapsed = loadFromLocalStorage(STORAGE_KEY_INPUT_COLLAPSED);
if (savedInputCollapsed === "1") {
  setInputCollapsed(true, false);
} else {
  if (inputPanel && inputPanel.dataset.customWidth) {
    inputPanel.style.width = `${inputPanel.dataset.customWidth}px`;
    inputPanel.style.flex = `0 0 ${inputPanel.dataset.customWidth}px`;
  }
  if (headerToggleInputBtn) headerToggleInputBtn.classList.add("active");
}

const savedStatsCollapsed = loadFromLocalStorage(STORAGE_KEY_STATS_COLLAPSED);
if (savedStatsCollapsed === "1") {
  setStatsCollapsed(true, false);
} else {
  if (statsPanel && statsPanel.dataset.customWidth) {
    statsPanel.style.width = `${statsPanel.dataset.customWidth}px`;
    statsPanel.style.flex = `0 0 ${statsPanel.dataset.customWidth}px`;
  }
  if (headerToggleStatsBtn) headerToggleStatsBtn.classList.add("active");
}

initResizers();

window.addEventListener("resize", () => {
  syncScroll(inputArea, inputHighlight, inputLineNumbers, inputWrapGuide);
  syncScroll(outputArea, outputHighlight, outputLineNumbers, outputWrapGuide);
});

// Initial stats if output already populated
if (outputArea && outputArea.value) {
  updateStatsPanel(null, 0, outputArea.value.split("\n").length, null);
}

const go = new Go();
WebAssembly.instantiateStreaming(fetch("main.wasm?v=3"), go.importObject)
  .then((result) => {
    go.run(result.instance);

    formatBtn.textContent = "Format DML";
    formatBtn.disabled = false;

    formatBtn.addEventListener("click", () => {
      statusDiv.textContent = "";
      currentError = null;
      const dml = inputArea.value;
      const lineCap = getLineWrap();
      const startTime = performance.now();
      try {
        const response = formatDML(dml, lineCap);
        const durationMs = performance.now() - startTime;
        if (response && response.error) {
          statusDiv.textContent = response.error;
          outputArea.value = dml; // Fallback: unformatted snippet
          saveToLocalStorage(STORAGE_KEY_OUTPUT, dml);

          let errorInfo = null;
          if (
            response.syntax_error &&
            typeof response.line === "number" &&
            typeof response.column === "number"
          ) {
            errorInfo = {
              line: response.line,
              column: response.column,
              length: response.length || 1,
            };
          } else {
            // Also try to extract line & column from error message if available
            const match = response.error.match(
              /line\s+(\d+)[:\s]+(?:col(?:umn)?\s+)?(\d+)/i,
            );
            if (match) {
              errorInfo = {
                line: parseInt(match[1], 10),
                column: parseInt(match[2], 10),
                length: 1,
              };
            }
          }

          currentError = errorInfo;
          currentTokens = null;
          updateOutputView(currentError, currentTokens);
          updateTraceButtons(null, 0);
          updateStatsPanel(
            0,
            0,
            outputArea.value.split("\n").length,
            durationMs,
          );

          if (errorInfo && typeof errorInfo.line === "number") {
            const lineHeight = 20;
            const targetScrollTop = Math.max(
              0,
              (errorInfo.line - 3) * lineHeight,
            );
            outputArea.scrollTop = targetScrollTop;
            syncScroll(
              outputArea,
              outputHighlight,
              outputLineNumbers,
              outputWrapGuide,
            );

            const charRange = getErrorCharRange(outputArea.value, errorInfo);
            if (charRange) {
              outputArea.selectionStart = charRange.start;
              outputArea.selectionEnd = charRange.start;
            }
          }
        } else if (response && response.result !== undefined) {
          statusDiv.textContent = "";
          currentError = null;
          currentTokens = response.tokens || null;
          outputArea.value = response.result;
          saveToLocalStorage(STORAGE_KEY_INPUT, dml);
          saveToLocalStorage(STORAGE_KEY_OUTPUT, response.result);
          updateOutputView(null, currentTokens);
          updateTraceButtons(response.trace, response.trace_count || 0);

          const lines = response.result.split("\n").length;
          const traceCount = response.trace_count || 0;
          const tokenCount = response.tokens ? response.tokens.length : 0;
          updateStatsPanel(tokenCount, traceCount, lines, durationMs);

          if (currentTokens) {
            tokenTypeToName.clear();
            for (const t of currentTokens) {
              if (!tokenTypeToName.has(t.type)) {
                tokenTypeToName.set(t.type, t.name);
              }
            }
          }
        } else {
          statusDiv.textContent = "Unexpected response from WebAssembly.";
          outputArea.value = dml;
          currentTokens = null;
          updateOutputView(null, null);
          updateTraceButtons(null, 0);
          updateStatsPanel(
            0,
            0,
            outputArea.value.split("\n").length,
            durationMs,
          );
        }
      } catch (err) {
        const durationMs = performance.now() - startTime;
        statusDiv.textContent = "Execution error: " + err;
        outputArea.value = dml;
        currentTokens = null;
        updateOutputView(null, null);
        updateTraceButtons(null, 0);
        updateStatsPanel(0, 0, outputArea.value.split("\n").length, durationMs);
      }
    });

    if (downloadTraceBtn) {
      downloadTraceBtn.addEventListener("click", () => {
        downloadTraceFile(currentTrace);
      });
    }

    if (openPerfettoBtn) {
      openPerfettoBtn.addEventListener("click", () => {
        openInPerfetto(currentTrace);
      });
    }

    if (openChromeTracingBtn) {
      openChromeTracingBtn.addEventListener("click", () => {
        downloadTraceFile(currentTrace);
        try {
          window.open("chrome://tracing", "_blank");
        } catch (e) {
          // Browsers block chrome:// URLs directly
        }
        showChromeTracingModal();
      });
    }

    if (closeChromeTracingModal) {
      closeChromeTracingModal.addEventListener("click", hideChromeTracingModal);
    }

    if (dismissChromeTracingModal) {
      dismissChromeTracingModal.addEventListener(
        "click",
        hideChromeTracingModal,
      );
    }

    if (chromeTracingModal) {
      chromeTracingModal.addEventListener("click", (e) => {
        if (e.target === chromeTracingModal) {
          hideChromeTracingModal();
        }
      });
    }

    document.addEventListener("keydown", (e) => {
      if (
        e.key === "Escape" &&
        chromeTracingModal &&
        !chromeTracingModal.classList.contains("hidden")
      ) {
        hideChromeTracingModal();
      }
    });

    if (modalDownloadTraceBtn) {
      modalDownloadTraceBtn.addEventListener("click", () => {
        downloadTraceFile(currentTrace);
      });
    }

    if (copyChromeTracingUrlBtn && chromeTracingUrlInput) {
      copyChromeTracingUrlBtn.addEventListener("click", () => {
        const text = chromeTracingUrlInput.value;
        if (navigator.clipboard && navigator.clipboard.writeText) {
          navigator.clipboard
            .writeText(text)
            .then(() => {
              const origText = copyChromeTracingUrlBtn.textContent;
              copyChromeTracingUrlBtn.textContent = "Copied!";
              setTimeout(() => {
                copyChromeTracingUrlBtn.textContent = origText;
              }, 1500);
            })
            .catch(() => {
              chromeTracingUrlInput.select();
              document.execCommand("copy");
            });
        } else {
          chromeTracingUrlInput.select();
          document.execCommand("copy");
        }
      });
    }
  })
  .catch((err) => {
    console.error("Failed to load WebAssembly module:", err);
    document.getElementById("status").textContent =
      "Failed to load WebAssembly module. Check console for details.";
  });

// Token selection listener
function handleTokenSelection() {
  if (statsPanel && statsPanel.classList.contains("collapsed")) return;
  if (!currentTokens || currentTokens.length === 0) return;
  const cursor = outputArea.selectionStart;
  const tokenIndex = currentTokens.findIndex(
    (t) => t.start <= cursor && cursor <= t.stop + 1,
  );
  const token = tokenIndex >= 0 ? currentTokens[tokenIndex] : null;

  if (token) {
    if (outputHighlight) {
      outputHighlight.classList.add("has-selection");
      const tokenTypeClass = `token-type-${token.type}`;
      outputHighlight.querySelectorAll(".token-outline").forEach((el, idx) => {
        const isSameType = el.classList.contains(tokenTypeClass);
        if (isSameType) {
          el.classList.add("selected-type-highlight");
        } else {
          el.classList.remove("selected-type-highlight");
        }

        const isExactSelected =
          idx === tokenIndex || el.dataset.tokenIndex === String(tokenIndex);
        if (isExactSelected) {
          el.classList.add("selected-token-highlight");
        } else {
          el.classList.remove("selected-token-highlight");
        }
      });
    }

    const details = document.getElementById("tokenInfoDetails");
    const box = document.getElementById("tokenInfoBox");
    if (details && box) {
      box.querySelector(".stats-placeholder-title").style.display = "none";
      details.style.display = "block";

      const events = parseTrace();
      let groupVal = token.groupPath;
      if (!groupVal && events && events.length > 0) {
        if (
          typeof tokenIndex === "number" &&
          events[tokenIndex]?.args?.right?.groupPath
        ) {
          groupVal = events[tokenIndex].args.right.groupPath;
        } else if (
          typeof tokenIndex === "number" &&
          events[tokenIndex + 1]?.args?.left?.groupPath
        ) {
          groupVal = events[tokenIndex + 1].args.left.groupPath;
        } else {
          // Find in trace events by matching line, column, text, or type
          const match = events.find(
            (e) =>
              (e.args?.right?.line === token.line &&
                e.args?.right?.column === token.column) ||
              (e.args?.left?.line === token.line &&
                e.args?.left?.column === token.column) ||
              (e.args?.right?.text === token.text &&
                e.args?.right?.line === token.line),
          );
          if (match) {
            if (
              match.args?.right?.line === token.line &&
              match.args?.right?.column === token.column
            ) {
              groupVal = match.args.right.groupPath;
            } else if (
              match.args?.left?.line === token.line &&
              match.args?.left?.column === token.column
            ) {
              groupVal = match.args.left.groupPath;
            } else if (match.args?.right?.text === token.text) {
              groupVal = match.args.right.groupPath;
            }
          }
        }
      }
      if (!groupVal) {
        groupVal = "-";
      }

      const elements = [
        { id: "tokenInfoName", val: token.name },
        { id: "tokenInfoType", val: token.type },
        { id: "tokenInfoText", val: token.text },
        { id: "tokenInfoLine", val: token.line },
        { id: "tokenInfoCol", val: token.column },
        { id: "tokenInfoStart", val: token.start },
        { id: "tokenInfoStop", val: token.stop },
        { id: "tokenInfoGroup", val: groupVal },
      ];

      elements.forEach(({ id, val }) => {
        const el = document.getElementById(id);
        if (el) {
          el.textContent = val;
          el.style.color = ""; // reset color in case it was "Copied!"
        }
      });

      updateWinLossTables(token, tokenIndex);
    }
  } else {
    if (outputHighlight) {
      outputHighlight.classList.remove("has-selection");
      outputHighlight.querySelectorAll(".token-outline").forEach((el) => {
        el.classList.remove("selected-type-highlight");
        el.classList.remove("selected-token-highlight");
      });
    }
  }
}

if (outputArea) {
  outputArea.addEventListener("click", handleTokenSelection);
  outputArea.addEventListener("keyup", (e) => {
    if (
      [
        "ArrowLeft",
        "ArrowRight",
        "ArrowUp",
        "ArrowDown",
        "Home",
        "End",
        "PageUp",
        "PageDown",
      ].includes(e.key)
    ) {
      handleTokenSelection();
    }
  });
}

// Copyable token values click listener
const tokenInfoDetails = document.getElementById("tokenInfoDetails");
if (tokenInfoDetails) {
  tokenInfoDetails.addEventListener("click", (e) => {
    const target = e.target.closest(".copyable-value");
    if (!target) return;

    const text = target.textContent;
    if (text === "Copied!") return;

    const applyCopiedState = () => {
      const origText = target.textContent;
      target.textContent = "Copied!";
      target.style.color = "var(--button-hover)";
      setTimeout(() => {
        if (target.textContent === "Copied!") {
          target.textContent = origText;
          target.style.color = "";
        }
      }, 1000);
    };

    if (navigator.clipboard && navigator.clipboard.writeText) {
      navigator.clipboard
        .writeText(text)
        .then(applyCopiedState)
        .catch((err) => console.error("Failed to copy:", err));
    } else {
      // Fallback
      const textarea = document.createElement("textarea");
      textarea.value = text;
      textarea.style.position = "fixed";
      textarea.style.opacity = "0";
      document.body.appendChild(textarea);
      textarea.select();
      try {
        document.execCommand("copy");
        applyCopiedState();
      } catch (err) {
        console.error("Fallback copy failed:", err);
      }
      document.body.removeChild(textarea);
    }
  });
}

function parseTrace() {
  if (parsedTraceCache) return parsedTraceCache;
  if (!currentTrace) return [];
  try {
    const data = JSON.parse(currentTrace);
    // trace could be an array of entries directly, or { traceEvents: [...] }
    if (Array.isArray(data)) {
      parsedTraceCache = data;
    } else if (data.traceEvents && Array.isArray(data.traceEvents)) {
      parsedTraceCache = data.traceEvents;
    } else {
      parsedTraceCache = [];
    }
    return parsedTraceCache;
  } catch (e) {
    console.error("Failed to parse trace:", e);
    parsedTraceCache = [];
    return parsedTraceCache;
  }
}

function getWinLossStats(selectedTokenId) {
  if (!winLossStatsCache) {
    winLossStatsCache = new Map();
    const events = parseTrace();
    for (const e of events) {
      if (!e.args || !e.args.left || !e.args.right) continue;

      const leftId = e.args.left.tokenType;
      const rightId = e.args.right.tokenType;
      const winner = e.args.winner;
      const sep = e.args.resolvedSep || "-";

      if (!winLossStatsCache.has(leftId)) {
        winLossStatsCache.set(leftId, {
          before: new Map(),
          after: new Map(),
        });
      }
      if (!winLossStatsCache.has(rightId)) {
        winLossStatsCache.set(rightId, {
          before: new Map(),
          after: new Map(),
        });
      }

      const record = (map, sep, outcome, opponentId) => {
        const key = `${sep}:::${outcome}`;
        let item = map.get(key);
        if (!item) {
          item = {
            sep,
            outcome,
            uniqueTokenTypes: new Set(),
            totalCount: 0,
          };
          map.set(key, item);
        }
        item.uniqueTokenTypes.add(opponentId);
        item.totalCount++;
      };

      // rightId is the selected token (opponent is on the left)
      const rightStats = winLossStatsCache.get(rightId);
      if (winner === 2) {
        record(rightStats.before, sep, "Win", leftId);
      } else if (winner === 1) {
        record(rightStats.before, sep, "Loss", leftId);
      } else if (winner === 3) {
        record(rightStats.before, sep, "Tie", leftId);
      }

      // leftId is the selected token (opponent is on the right)
      const leftStats = winLossStatsCache.get(leftId);
      if (winner === 1) {
        record(leftStats.after, sep, "Win", rightId);
      } else if (winner === 2) {
        record(leftStats.after, sep, "Loss", rightId);
      } else if (winner === 3) {
        record(leftStats.after, sep, "Tie", rightId);
      }
    }
  }

  return (
    winLossStatsCache.get(selectedTokenId) || {
      before: new Map(),
      after: new Map(),
    }
  );
}

function updateWinLossTables(tokenOrType, tokenIndex) {
  const tables = document.getElementById("tokenWinLossTables");
  if (!tables) return;
  tables.style.display = "block";

  const selectedTokenId =
    typeof tokenOrType === "object" && tokenOrType !== null
      ? tokenOrType.type
      : tokenOrType;

  const stats = getWinLossStats(selectedTokenId);

  const events = parseTrace();
  let idx = tokenIndex;
  if (typeof idx !== "number" || idx < 0) {
    if (
      typeof tokenOrType === "object" &&
      tokenOrType !== null &&
      currentTokens
    ) {
      idx = currentTokens.indexOf(tokenOrType);
    }
  }

  let decisionBefore = null;
  let decisionAfter = null;
  if (typeof idx === "number" && idx >= 0 && events && events.length > 0) {
    if (idx > 0 && idx < events.length) {
      decisionBefore = events[idx];
    }
    if (idx + 1 < events.length) {
      decisionAfter = events[idx + 1];
    }
  }

  const renderDecisionCard = (containerId, decision, role) => {
    const container = document.getElementById(containerId);
    if (!container) return;
    container.innerHTML = "";

    if (!decision || !decision.args) {
      const empty = document.createElement("div");
      empty.className = "decision-empty";
      empty.textContent =
        role === "before"
          ? "No prior token (start of file)"
          : "No next token (end of file)";
      container.appendChild(empty);
      return;
    }

    const { left, right, resolvedSep, breakPriority, winner } = decision.args;
    const opponent = role === "before" ? left : right;
    const isCurrentRight = role === "before";

    if (!opponent || (!opponent.tokenType && !opponent.text)) {
      const empty = document.createElement("div");
      empty.className = "decision-empty";
      empty.textContent =
        role === "before"
          ? "No prior token (start of file)"
          : "No next token (end of file)";
      container.appendChild(empty);
      return;
    }

    let outcome = "Tie";
    if (winner === (isCurrentRight ? 2 : 1)) {
      outcome = "Win";
    } else if (winner === (isCurrentRight ? 1 : 2)) {
      outcome = "Loss";
    }

    // Header bar with outcome and resolved properties
    const header = document.createElement("div");
    header.className = "decision-card-header";

    const outcomeGroup = document.createElement("div");
    outcomeGroup.className = "decision-card-outcome";

    const badge = document.createElement("span");
    badge.className = `badge-outcome ${
      outcome === "Win"
        ? "badge-win"
        : outcome === "Loss"
          ? "badge-loss"
          : "badge-tie"
    }`;
    badge.textContent = outcome;
    outcomeGroup.appendChild(badge);

    const sepPill = document.createElement("span");
    sepPill.className = "decision-detail-pill";
    sepPill.innerHTML = `Resolved: <code class="stats-table-sep">${
      resolvedSep || "-"
    }</code>`;
    outcomeGroup.appendChild(sepPill);

    const prioPill = document.createElement("span");
    prioPill.className = "decision-detail-pill";
    prioPill.innerHTML = `Break Prio: <code>${
      breakPriority !== undefined ? breakPriority : "-"
    }</code>`;
    outcomeGroup.appendChild(prioPill);

    header.appendChild(outcomeGroup);
    container.appendChild(header);

    // Grid with side-by-side token boxes
    const grid = document.createElement("div");
    grid.className = "decision-comparison-grid";

    const createTokenBox = (
      tokenInfo,
      title,
      isWinner,
      isLoser,
      sideName,
      sideFlag,
    ) => {
      const box = document.createElement("div");
      box.className = `decision-token-box${isWinner ? " is-winner" : ""}${
        isLoser ? " is-loser" : ""
      }`;

      const boxHeader = document.createElement("div");
      boxHeader.className = "decision-token-header";

      const titleEl = document.createElement("div");
      titleEl.className = "decision-token-title";
      titleEl.textContent = title;
      boxHeader.appendChild(titleEl);

      if (isWinner) {
        const winBadge = document.createElement("span");
        winBadge.className = "badge-outcome badge-win";
        winBadge.textContent = "Winner";
        boxHeader.appendChild(winBadge);
      }
      box.appendChild(boxHeader);

      const nameEl = document.createElement("div");
      nameEl.className = "decision-token-name";
      const typeName =
        tokenTypeToName.get(tokenInfo.tokenType) ||
        `Token(${tokenInfo.tokenType})`;
      const cleanText = tokenInfo.text
        ? tokenInfo.text.replace(/\r?\n/g, "\\n")
        : "";
      nameEl.textContent = `${typeName} `;
      if (cleanText) {
        const code = document.createElement("code");
        code.textContent = `"${cleanText}"`;
        nameEl.appendChild(code);
      }
      box.appendChild(nameEl);

      const props = document.createElement("div");
      props.className = "decision-flag-props";

      const addProp = (label, val, highlight = false) => {
        const row = document.createElement("div");
        row.className = "flag-prop";
        const lbl = document.createElement("span");
        lbl.className = "prop-label";
        lbl.textContent = label;
        const v = document.createElement("span");
        v.className = `prop-val${highlight ? " highlight-winner" : ""}`;
        v.textContent = val;
        row.appendChild(lbl);
        row.appendChild(v);
        props.appendChild(row);
      };

      addProp(`Sep (${sideName}):`, sideFlag?.separator || "-", isWinner);
      addProp(
        `Token Prio:`,
        sideFlag?.token_priority !== undefined ? sideFlag.token_priority : "-",
        isWinner,
      );
      addProp(
        `Break Prio:`,
        sideFlag?.break_priority !== undefined ? sideFlag.break_priority : "-",
      );
      addProp(`Indent:`, tokenInfo.formatFlag?.indent || "none");
      if (tokenInfo.formatFlag?.force_wrap) {
        addProp(`Force Wrap:`, "true");
      }
      addProp(`Group:`, tokenInfo.groupPath || "-");

      box.appendChild(props);
      return box;
    };

    // Left token
    const leftTitle =
      role === "before" ? "Opponent (Left)" : "Current Token (Left)";
    const leftIsWinner = winner === 1;
    const leftIsLoser = winner === 2;
    const leftBox = createTokenBox(
      left,
      leftTitle,
      leftIsWinner,
      leftIsLoser,
      "after",
      left.formatFlag?.after,
    );

    // Right token
    const rightTitle =
      role === "before" ? "Current Token (Right)" : "Opponent (Right)";
    const rightIsWinner = winner === 2;
    const rightIsLoser = winner === 1;
    const rightBox = createTokenBox(
      right,
      rightTitle,
      rightIsWinner,
      rightIsLoser,
      "before",
      right.formatFlag?.before,
    );

    grid.appendChild(leftBox);
    grid.appendChild(rightBox);
    container.appendChild(grid);
  };

  renderDecisionCard("currentBeforeDecisionCard", decisionBefore, "before");
  renderDecisionCard("currentAfterDecisionCard", decisionAfter, "after");

  const renderTable = (map, tableId) => {
    const tbody = document.querySelector(`#${tableId} tbody`);
    if (!tbody) return;
    tbody.innerHTML = "";

    const sorted = Array.from(map.values()).sort((a, b) => {
      if (a.sep !== b.sep) {
        return a.sep.localeCompare(b.sep);
      }
      const outcomeOrder = { Win: 1, Loss: 2, Tie: 3 };
      return (outcomeOrder[a.outcome] || 9) - (outcomeOrder[b.outcome] || 9);
    });

    if (sorted.length === 0) {
      const tr = document.createElement("tr");
      const td = document.createElement("td");
      td.colSpan = 4;
      td.className = "stats-table-empty";
      td.textContent = "No data";
      tr.appendChild(td);
      tbody.appendChild(tr);
      return;
    }

    for (const { sep, outcome, uniqueTokenTypes, totalCount } of sorted) {
      const tr = document.createElement("tr");

      const tdSep = document.createElement("td");
      tdSep.textContent = sep;
      tdSep.className = "stats-table-sep";

      const tdOutcome = document.createElement("td");
      const outcomeSpan = document.createElement("span");
      outcomeSpan.textContent = outcome;
      if (outcome === "Win") {
        outcomeSpan.className = "badge-outcome badge-win";
      } else if (outcome === "Loss") {
        outcomeSpan.className = "badge-outcome badge-loss";
      } else {
        outcomeSpan.className = "badge-outcome badge-tie";
      }
      tdOutcome.appendChild(outcomeSpan);

      const tdUnique = document.createElement("td");
      tdUnique.textContent = uniqueTokenTypes.size;
      const opponentNames = Array.from(uniqueTokenTypes)
        .map((id) => tokenTypeToName.get(id) || `Token(${id})`)
        .sort()
        .join(", ");
      tdUnique.title = opponentNames ? `Tokens: ${opponentNames}` : "";
      tdUnique.className = "stats-table-unique-tokens";

      const tdCount = document.createElement("td");
      tdCount.textContent = totalCount;

      tr.appendChild(tdSep);
      tr.appendChild(tdOutcome);
      tr.appendChild(tdUnique);
      tr.appendChild(tdCount);
      tbody.appendChild(tr);
    }
  };

  renderTable(stats.before, "beforeResolutionsTable");
  renderTable(stats.after, "afterResolutionsTable");
}
