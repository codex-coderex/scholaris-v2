<script lang="ts">
  type Props = {
    label: string
    sortBy: string
    order: string
    sortKey: string
    headerClass?: string
    onSort: (key: string) => void
  }

  let { label, sortBy, order, sortKey, headerClass = '', onSort }: Props = $props()

  let active = $derived(sortBy === sortKey)
</script>

<th class="ssis-col-header {headerClass} px-4 py-3 text-left">
  <button
    type="button"
    class="inline-flex items-center gap-1.5 text-xs font-semibold uppercase tracking-wider transition
      {active ? 'text-(--accent)' : 'text-(--text-muted) hover:text-(--text-primary)'}"
    onclick={() => onSort(sortKey)}
  >
    {label}
    {#if active}
      <span class="text-(--accent)">{order === 'ASC' ? '↑' : '↓'}</span>
    {:else}
      <span class="text-(--text-muted)">↕</span>
    {/if}
  </button>
</th>