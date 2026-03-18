<script lang="ts">
  import { onMount } from "svelte";
  import X from "@lucide/svelte/icons/x";
  import { toast } from "svelte-sonner";
  import { ScoopService, DNSOverride } from "../../../../bindings/changeme";

  let { cmd = $bindable("Create New Request") } = $props<{
    cmd: any;
  }>();

  let inputEl: HTMLInputElement | null = $state(null);

  let newVariable: string = $state("");
  let variableIPV4: string = $state("");
  let allOv: [DNSOverride[], string] | null = $state(null);

  let createNew: boolean = $state(true);
  let manageExisting: boolean = $state(false);

  async function createDNSOverride() {
    if (newVariable === "" || variableIPV4 === "") {
      toast.error("Please enter a valid variable and IPV4");
      return;
    }

    try {
      const ok = await ScoopService.CreateDNSOverride(
        new DNSOverride({ variable: newVariable, ipv4: variableIPV4 }),
      );

      if (ok) {
        toast.success(`Created DNS Override "${newVariable}"`);
        console.log(`Created DNS Override: ${newVariable}`);
      }
    } catch (error) {
      console.error(error);
    } finally {
      cmd = null;
    }
  }

  async function loadDNSOverrides() {
    try {
      allOv = await ScoopService.OpenDNSOverrides();
    } catch (error) {
      console.error(error);
    }
  }

  $effect(() => {
    if (cmd === "Configure DNS Override") {
      inputEl?.focus();
    }
  });

  onMount(() => {
    loadDNSOverrides();
  });
</script>

<div
  class={`border-border bg-background flex max-h-[40vh] min-h-[20vh] flex-col gap-5 rounded-sm border ${manageExisting ? `py-5` : `p-5`}`}
>
  <!-->Tab Buttons<-->
  <div class="mb-5 flex flex-row items-center justify-center gap-5">
    <button
      class={`focus:ring-offset-background h-9 w-45 rounded-sm border border-none focus:ring-2 
              focus:ring-green-400/20 focus:ring-offset-2 focus:outline-none ${createNew ? `bg-green-300/15` : `bg-muted/70`}`}
      onclick={() => {
        if (createNew) return;
        createNew = !createNew;
        manageExisting = !manageExisting;
      }}
    >
      New Variable
    </button>

    <button
      class={`focus:ring-offset-background h-9 w-45 rounded-sm border border-none focus:ring-2 
              focus:ring-green-400/20 focus:ring-offset-2 focus:outline-none ${manageExisting ? `bg-green-300/15` : `bg-muted/70`}`}
      onclick={() => {
        if (manageExisting) return;
        manageExisting = !manageExisting;
        createNew = !createNew;
      }}
    >
      Manage Existing
    </button>
  </div>

  <!-->Create Tab<-->
  {#if createNew}
    <div class="flex h-full flex-row items-center justify-center gap-5">
      <input
        class="focus:ring-offset-background bg-background border-border h-8 w-full
    min-w-0 rounded-sm border px-2 text-green-300 shadow-md
    focus:ring-2 focus:ring-green-400/20 focus:ring-offset-2 focus:outline-none"
        bind:value={newVariable}
        bind:this={inputEl}
        placeholder="Enter variable ..."
      />
      <input
        class="focus:ring-offset-background bg-background border-border h-8 w-full
    min-w-0 rounded-sm border px-2 text-green-300 shadow-md
    focus:ring-2 focus:ring-green-400/20 focus:ring-offset-2 focus:outline-none"
        bind:value={variableIPV4}
        placeholder="Enter IPV4 ..."
      />
    </div>

    <div class="flex w-full flex-row items-center justify-center gap-5">
      <button
        class="border-border bg-accent text-foreground focus:ring-offset-background inline-flex h-9 items-center
        justify-center rounded-sm border px-3
        text-sm hover:bg-green-400 hover:text-black focus:ring-2
        focus:ring-green-400/20 focus:ring-offset-2 focus:outline-none
        disabled:cursor-not-allowed disabled:opacity-50 disabled:hover:bg-transparent disabled:hover:text-green-500"
        onclick={createDNSOverride}>Create</button
      >

      <button
        class="border-border bg-accent text-foreground focus:ring-offset-background inline-flex h-9 items-center
        justify-center rounded-sm border px-3
        text-sm hover:bg-green-400 hover:text-black focus:ring-2
        focus:ring-green-400/20 focus:ring-offset-2 focus:outline-none
        disabled:cursor-not-allowed disabled:opacity-50 disabled:hover:bg-transparent disabled:hover:text-green-500"
        onclick={() => {
          cmd = null;
        }}>Cancel</button
      >
    </div>
  {:else if manageExisting}
    <div class="flex h-full w-full flex-col justify-center gap-5 overflow-y-auto p-5">
      {#each allOv ? allOv[0] : [] as ov}
        <div class="flex flex-row">
          <input
            class="focus:ring-offset-background bg-background border-border mr-5 h-8
    w-full min-w-0 rounded-sm border px-2 text-green-300 shadow-md
    focus:ring-2 focus:ring-green-400/20 focus:ring-offset-2 focus:outline-none"
            value={ov.variable}
            bind:this={inputEl}
            readonly
            placeholder="Enter variable ..."
          />
          <input
            class="focus:ring-offset-background bg-background border-border mr-2 h-8
    w-full min-w-0 rounded-sm border px-2 text-green-300 shadow-md
    focus:ring-2 focus:ring-green-400/20 focus:ring-offset-2 focus:outline-none"
            value={ov.ipv4}
            readonly
            placeholder="Enter IPV4 ..."
          />
          <div class="flex min-w-[2.5vh] items-center justify-center">
            <X class="flex h-[2.5vh] w-[2.5vw] text-red-800" />
          </div>
        </div>
      {/each}
    </div>
  {/if}
</div>
