<script lang="ts">
  type Props = {
    label: string
    sortBy: string
    order: string
    sortKey: string
    onSort: (key: string) => void
  }

  let { label, sortBy, order, sortKey, onSort }: Props = $props()
</script>

<th class="whitespace-nowrap px-3 py-2 text-left text-xs font-semibold text-slate-900">
  <button
    type="button"
    class="inline-flex items-center gap-1 text-left hover:text-slate-700"
    onclick={() => onSort(sortKey)}
  >
    <span>{label}</span>
    <span class="text-xs text-slate-500">{sortBy === sortKey ? (order === 'ASC' ? '↑' : '↓') : ''}</span>
  </button>
</th>
