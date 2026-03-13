<script lang="ts">
  import type { KV } from "../../../../bindings/changeme";
  import { onMount } from "svelte";

  type Mode = "isHeader" | "isQParam";

  let { content = $bindable<KV[]>([]), inputMode = $bindable<Mode>() } = $props<{
    inputMode: Mode;
    content: KV[];
  }>();

  let raw: string = $state("");
  let numOfLines = $derived(Math.max(1, raw.split("\n").length));

  let ta: HTMLTextAreaElement;
  let gutter: HTMLDivElement;

  function syncScroll() {
    if (!ta || !gutter) return;
    // move numbers up as textarea scrolls down (no gutter scrollbar)
    gutter.style.transform = `translateY(-${ta.scrollTop}px)`;
  }

  function contentToRaw(content: KV[]) {
    let r: string = "";

    for (const pair of content) {
      let row: string = "";
      let key: string = pair.key.trim();
      let value: string = pair.value.trim();

      if (key === "" && value === "") continue;

      if (inputMode === "isHeader") {
        row = key + ":" + value + "\n";
      } else if (inputMode === "isQParam") {
        row = key + "=" + value + "\n";
      }
      r = r + row;
    }
    raw = r;
  }

  function rawToContent() {
    if (raw === "") {
      content = [];
      return;
    }

    const r: string[] = raw.split("\n");
    let idx: number = 0;
    let tempContent: KV[] = [];

    for (const row of r) {
      if (row === "") continue;

      if (inputMode === "isHeader") {
        idx = row.indexOf(":");
      }

      if (inputMode === "isQParam") {
        idx = row.indexOf("=");
      }

      const key: string = (idx === -1 ? row : row.slice(0, idx)).trim();
      const val: string = (idx === -1 ? "" : row.slice(idx + 1)).trim();

      const newRow: KV = { key: key, value: val };
      tempContent.push(newRow);
    }
    content = tempContent;
  }

  $effect(() => {
    contentToRaw(content);
  });

  onMount(() => {
    contentToRaw(content);
  });
</script>

<div class="bg-accent m-3 flex h-64 overflow-hidden rounded-md border text-sm">
  <div class="relative w-10 shrink-0 overflow-hidden p-2 text-right select-none">
    <div bind:this={gutter} class="absolute top-2 right-0 left-0 will-change-transform">
      {#each Array(numOfLines) as _, i}
        <div class="h-5 pr-2 leading-5">{i + 1}</div>
      {/each}
    </div>
  </div>

  <textarea
    bind:this={ta}
    class="h-full flex-1 resize-none overflow-y-auto p-2 leading-5 text-green-300 outline-none"
    wrap="off"
    bind:value={raw}
    onchange={rawToContent}
    rows={10}
    onscroll={syncScroll}
    oninput={syncScroll}
    maxlength={5000}
  ></textarea>
</div>
