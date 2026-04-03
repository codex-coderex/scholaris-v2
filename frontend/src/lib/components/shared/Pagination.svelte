<script lang="ts">
  type Props = {
    page: number
    totalPages: number
    total: number
    onPrevious: () => void
    onNext: () => void
  }

  let { page, totalPages, total, onPrevious, onNext }: Props = $props()
</script>

<div class="mt-4 flex items-center justify-between gap-3 text-sm text-slate-600">
  <span>{total} total</span>
  <div class="flex items-center gap-2">
    <button
      type="button"
      class="rounded border border-slate-300 bg-white px-3 py-1.5 disabled:opacity-50"
      disabled={page === 1}
      onclick={(event) => {
        (event.currentTarget as HTMLButtonElement).blur()
        onPrevious()
      }}
    >Prev</button>
    <span>{page} / {totalPages}</span>
    <button
      type="button"
      class="rounded border border-slate-300 bg-white px-3 py-1.5 disabled:opacity-50"
      disabled={page >= totalPages}
      onclick={(event) => {
        (event.currentTarget as HTMLButtonElement).blur()
        onNext()
      }}
    >Next</button>
  </div>
</div>
