import type * as Monaco from "monaco-editor";

export function defineMyTheme(monaco: typeof Monaco) {
  monaco.editor.defineTheme("neon", {
    base: "vs-dark",
    inherit: true,
    colors: {
      "editorBracketHighlight.foreground1": "#a855f7", // purple
      "editorBracketHighlight.foreground2": "#3b82f6", // blue
      "editorBracketHighlight.foreground3": "#f97316", // orange
      "editorBracketHighlight.foreground4": "#a855f7",
      "editorBracketHighlight.foreground5": "#3b82f6",
      "editorBracketHighlight.foreground6": "#f97316",
      "editorBracketHighlight.unexpectedBracket.foreground": "#ef4444",
      "editor.background": "#050907",
      "editor.foreground": "#c7f9d4",
      "editorLineNumber.foreground": "#1f7a46",
      "editorLineNumber.activeForeground": "#22c55e",
      "editorCursor.foreground": "#22c55e",
      "editorCursor.background": "#050907",
      "editor.selectionBackground": "#14532d66",
      "editor.inactiveSelectionBackground": "#14532d33",
      "editor.selectionHighlightBackground": "#16653455",
      "editor.lineHighlightBackground": "#0b1510",
      "editor.lineHighlightBorder": "#0b1510",
      "editor.findMatchBackground": "#f9731666",
      "editor.findMatchHighlightBackground": "#f9731633",
      "editor.findRangeHighlightBackground": "#a855f733",
      "editorBracketMatch.background": "#22c55e33", // green-500 with alpha
      "editorBracketMatch.border": "#22c55eaa",
      "editorError.foreground": "#ef4444",
      "editorWarning.foreground": "#f59e0b",
      "editorInfo.foreground": "#3b82f6",
      "editorGutter.background": "#050907",
      "minimap.background": "#050907",
    },
    rules: [
      { token: "string.key.json", foreground: "22c55e" }, // green keys
      { token: "string.value.json", foreground: "3b82f6" }, // blue string values
      { token: "number.json", foreground: "f97316" }, // orange numbers
      { token: "keyword.json", foreground: "a855f7" }, // true/false/null (sometimes)
      { token: "delimiter.bracket.json", foreground: "6ee7b7" },
      { token: "delimiter.array.json", foreground: "6ee7b7" },
      { token: "delimiter.comma.json", foreground: "6ee7b7" },
      { token: "delimiter.colon.json", foreground: "6ee7b7" },
      { token: "delimiter.json", foreground: "6ee7b7" },

      { token: "", foreground: "c7f9d4" },
      { token: "comment", foreground: "3a8f63", fontStyle: "italic" },
      { token: "string", foreground: "22c55e" },
      { token: "number", foreground: "f97316" },
      { token: "keyword", foreground: "a855f7" },
      { token: "type", foreground: "3b82f6" },
      { token: "entity.name.function", foreground: "60a5fa" },
      { token: "variable", foreground: "a7f3d0" },
      { token: "attribute.name", foreground: "3b82f6" },
      { token: "delimiter", foreground: "6ee7b7" },
    ],
  });
}
