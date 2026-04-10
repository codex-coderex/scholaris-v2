<script lang="ts">
  import { onMount } from 'svelte'

  import { setupResponsivePageSize } from '$lib/shared/page-size'
  import { ProgramsViewModel } from '$lib/viewmodels/programs.vm.svelte'

  import SearchToolbar from '$lib/components/shared/SearchToolbar.svelte'
  import Pagination from '$lib/components/shared/Pagination.svelte'
  import ConfirmDialog from '$lib/components/overlays/ConfirmDialog.svelte'

  import ProgramsTable from './ProgramsTable.svelte'
  import ProgramFormModal from './ProgramFormModal.svelte'

  const vm = new ProgramsViewModel()

  export function openAdd() {
    vm.openAdd()
  }

  onMount(() => {
    const cleanup = setupResponsivePageSize(
      (next) => vm.query.setPageSize(next),
      () => void vm.load()
    )

    void vm.loadOptions().finally(() => vm.load())
    return cleanup
  })
</script>

<div class="space-y-5 font-dmsans">
  <div class="flex flex-col gap-3 lg:flex-row lg:items-end lg:justify-between">
    <div class="flex-1 min-w-0">
      <SearchToolbar
        search={vm.query.search}
        placeholder="Search programs..."
        onSearch={(value) => vm.handleSearch(value)}
      />
    </div>

    <label class="flex items-center gap-2 text-sm text-(--text-muted) lg:shrink-0">
      <span class="text-xs font-semibold uppercase tracking-[0.08em]">College</span>
      <select
        class="select select-bordered select-sm min-w-56 rounded-md border-(--border) bg-(--bg-base) text-(--text-primary) focus:border-(--accent) focus:outline-none"
        value={vm.filterCollegeCode}
        onchange={(event) => vm.handleCollegeFilter((event.currentTarget as HTMLSelectElement).value)}
      >
        <option value="">All Colleges</option>
        <option value={vm.noCollegeFilterCode}>(No College)</option>
        {#each vm.colleges as college}
          <option value={college.code}>{college.code} — {college.name}</option>
        {/each}
      </select>
    </label>
  </div>

  {#if vm.error}
    <div class="alert alert-error rounded-box border border-(--danger-dim) bg-(--danger-dim) text-(--danger)">{vm.error}</div>
  {/if}

  <div class="ssis-table-panel">
    <ProgramsTable
      rows={vm.programs}
      loading={vm.loading}
      emptyMessage={vm.error.toLowerCase().includes('no database') ? 'No database found. Start PostgreSQL and reload the app.' : 'No programs found'}
      sortBy={vm.query.sortBy}
      order={vm.query.order}
      onSort={(column) => vm.setSort(column)}
      onEdit={(program) => vm.openEdit(program)}
      onDelete={(program) => vm.askDelete(program)}
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

  <ProgramFormModal
    open={vm.formOpen}
    title={vm.formMode === 'create' ? 'Add Program' : 'Edit Program'}
    code={vm.code}
    name={vm.name}
    colleges={vm.colleges}
    selectedCollegeCode={vm.selectedCollegeCode}
    selectedCollegeLabel={vm.selectedCollegeLabel}
    formError={vm.formError}
    onCodeInput={(value) => {
      vm.code = value
    }}
    onNameInput={(value) => {
      vm.name = value
    }}
    onSelectCollege={(option) => vm.handleCollegeSelect(option)}
    onSave={() => void vm.saveProgram()}
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