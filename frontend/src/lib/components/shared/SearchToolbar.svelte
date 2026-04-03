<script lang="ts">
  type Props = {
    search: string
    placeholder: string
    onSearch: (value: string) => void
  }

  let { search, placeholder, onSearch }: Props = $props()
</script>

<div class="mb-4 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
  <input
    class="w-full rounded border border-slate-300 bg-white px-3 py-2 text-sm outline-none sm:max-w-md"
    placeholder={placeholder}
    value={search}
    oninput={(event) => onSearch((event.currentTarget as HTMLInputElement).value)}
  />
</div>
