<script lang="ts">
  import { onMount } from 'svelte'

  import { setupResponsivePageSize } from '$lib/shared/page-size'
  import { CollegesViewModel } from '$lib/viewmodels/colleges.vm.svelte'

  import SearchToolbar from '$lib/components/shared/SearchToolbar.svelte'
  import Pagination from '$lib/components/shared/Pagination.svelte'
  import ConfirmDialog from '$lib/components/overlays/ConfirmDialog.svelte'

  import CollegesTable from './CollegesTable.svelte'
  import CollegeFormModal from './CollegeFormModal.svelte'

  const vm = new CollegesViewModel()

  export function openAdd() {
    vm.openAdd()
  }

  onMount(() => {
    const cleanup = setupResponsivePageSize(
      (next) => vm.query.setPageSize(next),
      () => void vm.load()
    )

    void vm.load()
    return cleanup
  })
</script>

<div class="space-y-5 font-dmsans">
  <SearchToolbar
    search={vm.query.search}
    placeholder="Search colleges..."
    onSearch={(value) => vm.handleSearch(value)}
  />

  {#if vm.error}
    <div class="alert alert-error rounded-box border border-(--danger-dim) bg-(--danger-dim) text-(--danger)">{vm.error}</div>
  {/if}

  <div class="ssis-table-panel">
    <CollegesTable
      rows={vm.colleges}
      loading={vm.loading}
      emptyMessage={vm.error.toLowerCase().includes('no database') ? 'No database found. Start PostgreSQL and reload the app.' : 'No colleges found'}
      sortBy={vm.query.sortBy}
      order={vm.query.order}
      onSort={(column) => vm.setSort(column)}
      onEdit={(college) => vm.openEdit(college)}
      onDelete={(college) => vm.askDelete(college)}
    />

    <Pagination
      page={vm.query.page}
      totalPages={vm.query.totalPages}
      total={vm.query.total}
      busy={vm.loading || vm.isPageTransitioning}
      onPrevious={() => vm.previousPage()}
      onNext={() => vm.nextPage()}
    />
  </div>

  <CollegeFormModal
    open={vm.formOpen}
    title={vm.formMode === 'create' ? 'Add College' : 'Edit College'}
    code={vm.code}
    name={vm.name}
    formError={vm.formError}
    onCodeInput={(value) => {
      vm.code = value
    }}
    onNameInput={(value) => {
      vm.name = value
    }}
    onSave={() => void vm.saveCollege()}
    onClose={() => vm.closeForm()}
  />

  <ConfirmDialog
    open={vm.confirmOpen}
    title="Confirm Delete"
    message={vm.deleteTarget ? `Are you sure you want to delete ${vm.deleteTarget.code}? This action cannot be undone.` : ''}
    confirmLabel="Delete"
    onConfirm={() => void vm.remove()}
    onClose={() => vm.closeConfirm()}
  />
</div>