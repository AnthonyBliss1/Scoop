<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import * as monaco from "monaco-editor";
  import { defineMyTheme } from "$lib/monaco/theme.ts";

  export let value = "";
  export let contentType = "";

  let editorElement: HTMLDivElement;
  let editor: monaco.editor.IStandaloneCodeEditor | null = null;
  let model: monaco.editor.ITextModel | null = null;

  function languageFromContentType(ct: string): string {
    switch (ct) {
      case "application/json":
      case "text/json":
        return "json";
      case "text/typescript":
      case "application/typescript":
        return "typescript";
      case "text/javascript":
      case "application/javascript":
        return "javascript";
      default:
        return "plaintext";
    }
  }

  function uriForLanguage(lang: string) {
    const ext =
      lang === "json"
        ? "json"
        : lang === "typescript"
          ? "ts"
          : lang === "javascript"
            ? "js"
            : "txt";

    return monaco.Uri.parse(`inmemory://viewer/${crypto.randomUUID()}.${ext}`);
  }

  async function createModel(value: string, lang: string) {
    if (!editor) return;
    model?.dispose();
    model = monaco.editor.createModel(value, lang, uriForLanguage(lang));
    model.updateOptions({
      insertSpaces: true,
      tabSize: 2,
    });

    editor.setModel(model);
  }

  onMount(async () => {
    const lang: string = languageFromContentType(contentType);

    defineMyTheme(monaco);
    monaco.editor.setTheme("neon");

    editor = monaco.editor.create(editorElement, {
      automaticLayout: true,
      minimap: { enabled: false },
      readOnly: true,
      wordWrap: "off",
      theme: "neon",
      bracketPairColorization: { enabled: false },
      guides: { bracketPairs: false },
    });

    await createModel(value, lang);
    console.log("lang:", model?.getLanguageId());
  });

  $: if (editor && model) {
    const lang = languageFromContentType(contentType);
    if (value !== model.getValue()) {
      createModel(value, lang);
      console.log("lang:", model?.getLanguageId());
    }
  }

  onDestroy(() => {
    editor?.dispose();
    model?.dispose();
    editor = null;
    model = null;
  });
</script>

<div class="min-h-0 flex-1 overflow-hidden" bind:this={editorElement}></div>
