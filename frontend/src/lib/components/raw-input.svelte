<script lang="ts">
  let content = $state("");
  let numOfLines = $derived(Math.max(1, content.split("\n").length));

  let ta: HTMLTextAreaElement;
  let gutter: HTMLDivElement;

  function syncScroll() {
    if (!ta || !gutter) return;
    // move numbers up as textarea scrolls down (no gutter scrollbar)
    gutter.style.transform = `translateY(-${ta.scrollTop}px)`;
  }
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
    class="text-muted-foreground h-full flex-1 resize-none overflow-y-auto p-2 leading-5 outline-none"
    bind:value={content}
    rows={10}
    onscroll={syncScroll}
    oninput={syncScroll}
  ></textarea>
</div>
