<script lang="ts">
  type Props = {
    onEdit?: () => void
    onDelete: () => void
  }

  let { onEdit = () => {}, onDelete }: Props = $props()
</script>

<div class="flex items-center justify-center gap-2">
  <button
    type="button"
    onclick={onEdit}
    class="btn btn-xs rounded-md border border-(--border) bg-(--bg-elevated) px-2.5 normal-case text-xs font-medium text-(--text-secondary) shadow-none transition duration-200 hover:bg-(--accent-dim) hover:text-(--accent) focus:outline-none focus-visible:outline-none"
  >Edit</button>

  <button
    type="button"
    onclick={onDelete}
    class="btn btn-ghost btn-xs rounded-md px-1 normal-case text-xs font-medium text-(--text-muted) shadow-none transition duration-200 hover:bg-(--danger-dim) hover:text-(--danger) focus:outline-none focus-visible:outline-none"
  >Delete</button>
</div>
