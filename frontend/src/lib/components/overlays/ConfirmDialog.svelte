<script lang="ts">
  import Modal from './Modal.svelte'

  type Props = {
    open: boolean
    title: string
    message: string
    confirmLabel?: string
    cancelLabel?: string
    onConfirm: () => void
    onClose: () => void
  }

  let {
    open,
    title,
    message,
    confirmLabel = 'Delete',
    cancelLabel = 'Cancel',
    onConfirm,
    onClose
  }: Props = $props()
</script>

<Modal {open} {title} {onClose} size="sm">
  <div class="space-y-6">
    <p class="text-sm leading-6 text-(--text-secondary)">{message}</p>

    <div class="modal-action justify-start gap-3">
      <button type="button" class="btn btn-error" onclick={onConfirm}>{confirmLabel}</button>
      <button type="button" class="btn btn-ghost" onclick={onClose}>{cancelLabel}</button>
    </div>
  </div>
</Modal>