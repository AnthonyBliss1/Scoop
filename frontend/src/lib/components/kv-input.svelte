<script lang="ts">
  import type { KV } from "bindings/changeme";
  import Plus from "@lucide/svelte/icons/plus";
  import Minus from "@lucide/svelte/icons/minus";

  type Mode = "isHeader" | "isQParam" | undefined;

  let { inputMode = $bindable<Mode>(), content = $bindable<KV[]>([]) } = $props<{
    inputMode: Mode;
    content: KV[];
  }>();

  let scroller: HTMLDivElement;

  function addRow() {
    // max 100 rows
    if (content.length > 100) return;

    content = [...content, { key: "", value: "" }];
  }

  function removeRow() {
    if (content.length === 1) return;
    content = content.slice(0, -1);
  }

  $effect(() => {
    if (content.length === 0) content.push({ key: "", value: "" });
  });
</script>

<div class="bg-accent m-3 flex h-64 flex-col overflow-hidden rounded-md border text-sm">
  <div class="h-full flex-1 overflow-y-auto" bind:this={scroller}>
    {#each content as row, i}
      <div class="flex h-10 items-center">
        <span class="w-10 shrink-0 p-2 text-right leading-5 select-none">{i + 1}</span>

        <div class="flex flex-1 gap-3 p-2 pr-10">
          <input
            class="border-border focus:ring-offset-background w-full border px-1 text-green-300 focus:ring-2 focus:ring-green-400/20 focus:ring-offset-2 focus:outline-none"
            placeholder="Key"
            bind:value={row.key}
          />
          <input
            class="border-border focus:ring-offset-background w-full border px-1 text-green-300 focus:ring-2 focus:ring-green-400/20 focus:ring-offset-2 focus:outline-none"
            placeholder="Value"
            bind:value={row.value}
          />
        </div>
      </div>
    {/each}
  </div>
  <footer class="flex flex-row items-center justify-end gap-1 px-3 text-2xl">
    <div>
      <button
        class="flex h-7 w-7 items-center justify-center p-1 focus:outline-none"
        onclick={addRow}
      >
        <Plus size={18} strokeWidth={3} />
      </button>
    </div>
    <div>
      <button
        class="flex h-7 w-7 items-center justify-center p-1 focus:outline-none"
        onclick={removeRow}
      >
        <Minus size={18} strokeWidth={3} />
      </button>
    </div>
  </footer>
</div>
