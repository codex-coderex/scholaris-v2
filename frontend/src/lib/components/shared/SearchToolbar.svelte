<script lang="ts">
  import type { AppMode } from '$lib/types/schema'

  type Props = {
    search: string
    placeholder: string
    mode?: AppMode
    onSearch: (value: string) => void
  }

  let { search, placeholder, mode = 'light', onSearch }: Props = $props()
</script>

<div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
  <div class="ssis-search-wrap">
    <span class="pointer-events-none absolute inset-y-0 left-3 flex items-center text-(--text-muted)">
      <svg class="h-4 w-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
        <circle cx="11" cy="11" r="8"/><path d="m21 21-4.35-4.35"/>
      </svg>
    </span>
    <input
      class="ssis-search-input"
      {placeholder}
      value={search}
      oninput={(e) => onSearch((e.currentTarget as HTMLInputElement).value)}
    />
  </div>

</div>