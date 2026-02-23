<script lang="ts">
  import { toast } from "svelte-sonner";
  import { Backend, DNSOverride } from "../../../bindings/changeme";
  import { onMount } from "svelte";

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
      const ok = await Backend.CreateDNSOverride(
        new DNSOverride({ variable: newVariable, ipv4: variableIPV4 }),
      );

      if (ok) {
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
      allOv = await Backend.OpenDNSOverrides();
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

<div class="border-border bg-background flex min-h-[20vh] flex-col gap-5 rounded-sm border p-5">
  <!-->Tab Buttons<-->
  <div class="mb-5 flex flex-row items-center justify-center gap-5">
    <button
      class={`border-border focus:ring-offset-background h-9 rounded-sm border px-3 focus:ring-2 
              focus:ring-green-400/20 focus:ring-offset-2 focus:outline-none ${createNew ? `bg-green-300/10` : ``}`}
      onclick={() => {
        if (createNew) return;
        createNew = !createNew;
        manageExisting = !manageExisting;
      }}
    >
      Create New
    </button>

    <button
      class={`border-border focus:ring-offset-background h-9 rounded-sm border px-3 focus:ring-2 
              focus:ring-green-400/20 focus:ring-offset-2 focus:outline-none ${manageExisting ? `bg-green-300/10` : ``}`}
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

    <div class="flex w-full flex-row items-center justify-center gap-10">
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
    <div class="flex h-full flex-row items-center justify-center gap-5">
      {#each allOv ? allOv[0] : [] as ov}
        <input
          class="focus:ring-offset-background bg-background border-border h-8 w-full
    min-w-0 rounded-sm border px-2 text-green-300 shadow-md
    focus:ring-2 focus:ring-green-400/20 focus:ring-offset-2 focus:outline-none"
          bind:value={ov.variable}
          readonly
          placeholder="Enter variable ..."
        />
        <input
          class="focus:ring-offset-background bg-background border-border h-8 w-full
    min-w-0 rounded-sm border px-2 text-green-300 shadow-md
    focus:ring-2 focus:ring-green-400/20 focus:ring-offset-2 focus:outline-none"
          bind:value={ov.ipv4}
          readonly
          placeholder="Enter IPV4 ..."
        />
      {/each}
    </div>
  {/if}
</div>
