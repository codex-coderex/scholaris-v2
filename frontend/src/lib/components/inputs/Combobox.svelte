<script lang="ts">
  type Option = {
    code: string
    name: string
    college_code?: string
  }

  type Props = {
    label: string
    placeholder: string
    options: Option[]
    selectedCode: string
    selectedLabel: string
    onSelect: (option: Option | null) => void
    disabled?: boolean
  }

  let {
    label,
    placeholder,
    options,
    selectedCode,
    selectedLabel,
    onSelect,
    disabled = false
  }: Props = $props()

  let open = $state(false)
  let query = $state('')
  let inputElement: HTMLInputElement | null = null

  let filtered = $derived(
    options.filter((option) => {
      const search = query.trim().toLowerCase()
      if (!search) return true
      return `${option.code}${option.name}`.toLowerCase().includes(search)
    })
  )

  $effect(() => {
    if (!open) {
      query = selectedLabel
    }
  })

  function choose(option: Option) {
    query = `${option.code} — ${option.name}`
    onSelect(option)
    open = false
  }

  function clearSelection() {
    query = ''
    onSelect(null)
    open = true
    inputElement?.focus()
  }

  function handleBlur() {
    window.setTimeout(() => {
      open = false
      if (!selectedCode) {
        query = ''
      }
    }, 120)
  }
</script>

<label class="form-control w-full gap-2 overflow-visible">
  <span class="label-text text-[0.78rem] font-semibold uppercase tracking-[0.08em] text-(--text-muted)">{label}</span>
  <div class="relative overflow-visible">
    <input
      bind:this={inputElement}
      class="input input-bordered w-full bg-(--bg-elevated) pr-10 text-sm text-(--text-primary) border-(--border) placeholder:text-(--text-muted) focus:border-(--accent) focus:outline-none"
      {placeholder}
      value={query}
      autocomplete="off"
      disabled={disabled}
      onfocus={() => {
        open = true
      }}
      onblur={handleBlur}
      oninput={(event) => {
        const input = event.currentTarget as HTMLInputElement
        query = input.value
        open = true
        if (!input.value.trim()) clearSelection()
      }}
    />

    {#if selectedCode}
      <button
        type="button"
        class="btn btn-ghost btn-xs absolute right-2 top-2 rounded-md px-2 text-(--text-muted) hover:bg-(--bg-hover) hover:text-(--text-primary)"
        onclick={() => {
          clearSelection()
          open = true
        }}
      >
        Clear
      </button>
    {/if}

    {#if open}
      <div class="ssis-combobox-menu absolute left-0 right-0 top-[calc(100%+0.45rem)] z-90 max-h-72 overflow-y-auto overflow-x-hidden rounded-xl border border-(--border) bg-(--bg-surface) shadow-[0_22px_60px_rgba(6,12,24,0.22)]">
        {#if filtered.length === 0}
          <div class="px-4 py-3 text-sm italic text-(--text-muted)">No results found</div>
        {:else}
          {#each filtered as option}
            <button
              type="button"
              class="block w-full border-b border-(--border-subtle) px-4 py-3 text-left text-sm text-(--text-secondary) last:border-b-0 hover:bg-(--bg-hover) hover:text-(--text-primary)"
              onclick={() => choose(option)}
            >
              <span class="block whitespace-normal wrap-break-word font-medium text-(--text-primary)">{option.code} <span class="font-normal text-(--text-muted)">— {option.name}</span></span>
            </button>
          {/each}
        {/if}
      </div>
    {/if}
  </div>
</label>