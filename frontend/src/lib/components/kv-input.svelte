<script lang="ts">
  import type { KV } from "bindings/changeme";
  import { onDestroy, onMount } from "svelte";

  let { rawContent = $bindable<string>("") } = $props<{ rawContent: string }>();
  let rows: KV[] = $state([]);

  let scroller: HTMLDivElement;

  function addRow() {
    rows = [...rows, { key: "", value: "" }];
  }

  function removeRow() {
    if (rows.length === 1) return;
    rows = rows.filter((_, idx) => idx !== 1);
  }

  function rawToKV() {
    if (rawContent === "") {
      rows.push({ key: "", value: "" });
      return;
    }

    const r: string[] = rawContent.split("\n");

    for (const row of r) {
      if (row === "") return;

      const item: string[] = row.split(":");

      const key: string = item[0];
      const val: string = item[1];

      const newRow: KV = { key: key, value: val };
      rows.push(newRow);
    }
  }

  function kvToRaw() {
    let newContent: string = "";

    for (const row of rows) {
      if (row.key === "" && row.value === "") return;

      const line: string = row.key + ":" + row.value + "\n";
      newContent = newContent + line;
    }

    rawContent = newContent;
  }

  onMount(() => {
    rawToKV();
  });

  onDestroy(() => {
    kvToRaw();
  });
</script>

<div class="bg-accent m-3 flex h-64 flex-col overflow-hidden rounded-md border text-sm">
  <div class="h-full flex-1 overflow-y-auto" bind:this={scroller}>
    {#each rows as row, i}
      <div class="flex h-10 items-center">
        <span class="w-10 shrink-0 p-2 text-right leading-5 select-none">{i + 1}</span>

        <div class="flex flex-1 gap-5 p-2">
          <input class="border-border w-full border px-1" placeholder="Key" bind:value={row.key} />
          <input
            class="border-border w-full border px-1"
            placeholder="Value"
            bind:value={row.value}
          />
        </div>
      </div>
    {/each}
  </div>
  <footer class="flex flex-row items-center justify-end gap-5 px-3 text-2xl">
    <div>
      <button onclick={addRow}>+</button>
    </div>
    <div>
      <button onclick={removeRow}>-</button>
    </div>
  </footer>
</div>
