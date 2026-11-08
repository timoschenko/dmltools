function escapeHtml(text) {
  return text
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#039;");
}

const KEYWORDS = [
  "attribute",
  "bank",
  "bitfield",
  "break",
  "case",
  "connect",
  "const",
  "constant",
  "continue",
  "default",
  "delete",
  "device",
  "dml",
  "do",
  "else",
  "export",
  "extern",
  "field",
  "for",
  "foreach",
  "goto",
  "group",
  "if",
  "implement",
  "import",
  "in each",
  "inline",
  "interface",
  "layout",
  "local",
  "method",
  "new",
  "param",
  "port",
  "provisional",
  "register",
  "restrict",
  "return",
  "saved",
  "select",
  "session",
  "sizeof",
  "static",
  "struct",
  "subdevice",
  "switch",
  "template",
  "typedef",
  "typeof",
  "volatile",
  "while",
];

const TYPES = [
  "void",
  "char",
  "short",
  "int",
  "long",
  "float",
  "double",
  "signed",
  "unsigned",
  "bool",
  "auto",
  "size_t",
  "int8",
  "int16",
  "int32",
  "int64",
  "uint8",
  "uint16",
  "uint32",
  "uint64",
];

const keywordsPattern = `\\b(?:${KEYWORDS.join("|")})\\b`;
const typesPattern = `\\b(?:${TYPES.join("|")})\\b`;

const TOKEN_REGEX = new RegExp(
  [
    "(/\\/[^\\n]*|/\\*[\\s\\S]*?\\*/)", // 1: comments (//... or /*...*/)
    "(\"(?:\\\\.|[^\"\\\\])*\"|'(?:\\\\.|[^'\\\\])*')", // 2: strings ("..." or '...')
    "(^|\\n)(\\s*#\\s*[a-zA-Z_]\\w*)", // 3, 4: preprocessor directives (#include, etc.)
    `(${keywordsPattern})`, // 5: keywords
    `(${typesPattern})`, // 6: types
    "(\\b0x[0-9a-fA-F]+\\b|\\b\\d+(?:\\.\\d+)?(?:[eE][+-]?\\d+)?\\b)", // 7: numbers
  ].join("|"),
  "g",
);

function formatText(rawText, showWhitespace) {
  let escaped = escapeHtml(rawText);
  if (showWhitespace) {
    escaped = escaped
      .replace(/ /g, '<span class="hl-space">·</span>')
      .replace(/\t/g, '<span class="hl-tab">\t</span>');
  }
  return escaped;
}

function getErrorCharRange(code, errorInfo) {
  if (
    !errorInfo ||
    typeof errorInfo.line !== "number" ||
    typeof errorInfo.column !== "number"
  ) {
    return null;
  }
  const targetLineNum = errorInfo.line; // 1-based
  const targetCol = Math.max(0, errorInfo.column); // 0-based
  const lines = code.split("\n");
  if (targetLineNum < 1 || targetLineNum > lines.length) {
    return null;
  }

  let lineStartIdx = 0;
  for (let i = 0; i < targetLineNum - 1; i++) {
    lineStartIdx += lines[i].length + 1; // +1 for \n
  }
  const currentLine = lines[targetLineNum - 1];

  if (targetCol >= currentLine.length) {
    // Error is at end of line (e.g. missing token)
    const errorIdx = lineStartIdx + currentLine.length;
    return {
      start: errorIdx,
      end: errorIdx,
      isEol: true,
      line: targetLineNum,
      inserted: false,
    };
  }

  const start = lineStartIdx + targetCol;
  let len = errorInfo.length;
  if (!len || len < 1) {
    const match = currentLine
      .slice(targetCol)
      .match(/^[a-zA-Z0-9_]+|[^\s\w]|\s+/);
    len = match ? match[0].length : 1;
  }
  const end = Math.min(start + len, lineStartIdx + currentLine.length);
  return {
    start,
    end,
    isEol: false,
    line: targetLineNum,
    inserted: false,
  };
}

function renderPiece(text, cssClass, showWhitespace) {
  const formatted = formatText(text, showWhitespace);
  if (!cssClass) return formatted;
  return `<span class="${cssClass}">${formatted}</span>`;
}

function renderChunk(
  rawText,
  chunkStart,
  cssClass,
  showWhitespace,
  errorRange,
) {
  if (!rawText) return "";
  if (!errorRange) {
    return renderPiece(rawText, cssClass, showWhitespace);
  }

  const chunkEnd = chunkStart + rawText.length;

  if (errorRange.isEol) {
    if (
      !errorRange.inserted &&
      chunkStart <= errorRange.start &&
      errorRange.start <= chunkEnd
    ) {
      errorRange.inserted = true;
      const offset = errorRange.start - chunkStart;
      const before = rawText.slice(0, offset);
      const after = rawText.slice(offset);
      const eolMarker = `<mark class="syntax-error-token syntax-error-eol" title="Syntax error at line ${errorRange.line}">&nbsp;</mark>`;
      return (
        renderPiece(before, cssClass, showWhitespace) +
        eolMarker +
        renderPiece(after, cssClass, showWhitespace)
      );
    }
    return renderPiece(rawText, cssClass, showWhitespace);
  }

  const overlapStart = Math.max(chunkStart, errorRange.start);
  const overlapEnd = Math.min(chunkEnd, errorRange.end);

  if (overlapStart < overlapEnd) {
    let out = "";
    const prefix = rawText.slice(0, overlapStart - chunkStart);
    const errText = rawText.slice(
      overlapStart - chunkStart,
      overlapEnd - chunkStart,
    );
    const suffix = rawText.slice(overlapEnd - chunkStart);

    if (prefix) {
      out += renderPiece(prefix, cssClass, showWhitespace);
    }
    if (errText) {
      const inner = renderPiece(errText, cssClass, showWhitespace);
      out += `<mark class="syntax-error-token" title="Syntax error at line ${errorRange.line}">${inner}</mark>`;
    }
    if (suffix) {
      out += renderPiece(suffix, cssClass, showWhitespace);
    }
    return out;
  }

  return renderPiece(rawText, cssClass, showWhitespace);
}

function highlightCode(
  code,
  showWhitespace = false,
  errorInfo = null,
  tokens = null,
) {
  let html = "";
  let lastIndex = 0;

  const errorRange = getErrorCharRange(code, errorInfo);

  if (tokens && tokens.length > 0) {
    const hashString = (str) => {
      let hash = 0;
      for (let i = 0; i < str.length; i++) {
        hash = str.charCodeAt(i) + ((hash << 5) - hash);
      }
      return hash;
    };

    const getColor = (typeStr) => {
      const colors = [
        "#ff7b72",
        "#79c0ff",
        "#a5d6ff",
        "#ffa657",
        "#3fb950",
        "#d2a8ff",
        "#8b949e",
        "#f0883e",
        "#ff9bce",
        "#7ee787",
      ];
      return colors[Math.abs(hashString(typeStr)) % colors.length];
    };

    for (let i = 0; i < tokens.length; i++) {
      const token = tokens[i];
      if (token.start > lastIndex) {
        html += renderChunk(
          code.slice(lastIndex, token.start),
          lastIndex,
          null,
          showWhitespace,
          errorRange,
        );
      }
      // token.stop is inclusive index
      const tokenText = code.slice(token.start, token.stop + 1);
      const color = getColor(token.name);

      const innerHtml = renderChunk(
        tokenText,
        token.start,
        null,
        showWhitespace,
        errorRange,
      );
      html += `<span class="token-outline token-type-${token.type}" data-token-index="${i}" style="border-color: ${color};">${innerHtml}</span>`;
      lastIndex = token.stop + 1;
    }

    if (lastIndex < code.length) {
      html += renderChunk(
        code.slice(lastIndex),
        lastIndex,
        null,
        showWhitespace,
        errorRange,
      );
    }

    if (code.endsWith("\n")) {
      html += "\n";
    }
    return html;
  }

  let match;

  TOKEN_REGEX.lastIndex = 0;

  while ((match = TOKEN_REGEX.exec(code)) !== null) {
    if (match.index > lastIndex) {
      html += renderChunk(
        code.slice(lastIndex, match.index),
        lastIndex,
        null,
        showWhitespace,
        errorRange,
      );
    }

    const [
      full,
      comment,
      str,
      prepPrefix,
      prepDirective,
      keyword,
      type,
      number,
    ] = match;

    if (comment) {
      html += renderChunk(
        comment,
        match.index,
        "hl-comment",
        showWhitespace,
        errorRange,
      );
    } else if (str) {
      html += renderChunk(
        str,
        match.index,
        "hl-string",
        showWhitespace,
        errorRange,
      );
    } else if (prepDirective) {
      html += renderChunk(
        prepPrefix,
        match.index,
        null,
        showWhitespace,
        errorRange,
      );
      html += renderChunk(
        prepDirective,
        match.index + prepPrefix.length,
        "hl-preprocessor",
        showWhitespace,
        errorRange,
      );
    } else if (keyword) {
      html += renderChunk(
        keyword,
        match.index,
        "hl-keyword",
        showWhitespace,
        errorRange,
      );
    } else if (type) {
      html += renderChunk(
        type,
        match.index,
        "hl-type",
        showWhitespace,
        errorRange,
      );
    } else if (number) {
      html += renderChunk(
        number,
        match.index,
        "hl-number",
        showWhitespace,
        errorRange,
      );
    } else {
      html += renderChunk(full, match.index, null, showWhitespace, errorRange);
    }

    lastIndex = TOKEN_REGEX.lastIndex;
  }

  if (lastIndex < code.length) {
    html += renderChunk(
      code.slice(lastIndex),
      lastIndex,
      null,
      showWhitespace,
      errorRange,
    );
  }

  // If EOL error was not inserted yet (e.g. at EOF)
  if (errorRange && errorRange.isEol && !errorRange.inserted) {
    html += `<mark class="syntax-error-token syntax-error-eol" title="Syntax error at line ${errorRange.line}">&nbsp;</mark>`;
    errorRange.inserted = true;
  }

  if (code.endsWith("\n")) {
    html += "\n";
  }

  return html;
}

if (typeof module !== "undefined" && module.exports) {
  module.exports = { highlightCode, getErrorCharRange };
}
