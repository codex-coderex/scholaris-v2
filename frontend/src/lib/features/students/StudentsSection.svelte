<script lang="ts">
  import { onMount } from 'svelte'

  import type { AppMode, StudentGender } from '$lib/types/schema'
  import { setupResponsivePageSize } from '$lib/shared/page-size'
  import { StudentsViewModel } from '$lib/viewmodels/students.vm.svelte'

  import SearchToolbar from '$lib/components/shared/SearchToolbar.svelte'
  import Pagination from '$lib/components/shared/Pagination.svelte'
  import ConfirmDialog from '$lib/components/overlays/ConfirmDialog.svelte'

  import StudentsTable from './StudentsTable.svelte'
  import StudentFormModal from './StudentFormModal.svelte'

  type Props = {
    mode?: AppMode
  }

  let { mode = 'light' }: Props = $props()

  const vm = new StudentsViewModel()

  export function openAdd() {
    vm.openAdd()
  }

  onMount(() => {
    const cleanup = setupResponsivePageSize(
      (next) => vm.query.setPageSize(next),
      () => void vm.load()
    )

    void vm.loadOptions().then(() => vm.load())
    return cleanup
  })
</script>

<div class="space-y-5 font-dmsans">
  <SearchToolbar
    {mode}
    search={vm.query.search}
    placeholder="Search students..."
    onSearch={(value) => vm.handleSearch(value)}
  />

  {#if vm.error}
    <div class="alert alert-error rounded-box border border-(--danger-dim) bg-(--danger-dim) text-(--danger)">{vm.error}</div>
  {/if}

  <div class="ssis-table-panel">
    <StudentsTable
      {mode}
      rows={vm.students}
      loading={vm.loading}
      sortBy={vm.query.sortBy}
      order={vm.query.order}
      onSort={(column) => vm.setSort(column)}
      onEdit={(student) => vm.openEdit(student)}
      onDelete={(student) => vm.askDelete(student)}
    />

    <Pagination
      {mode}
      page={vm.query.page}
      totalPages={vm.query.totalPages}
      total={vm.query.total}
      busy={vm.loading || vm.isPageTransitioning}
      onPrevious={() => vm.previousPage()}
      onNext={() => vm.nextPage()}
      onGoToPage={(pageNumber) => vm.goToPage(pageNumber)}
      onJump={(event) => vm.handlePageJump(event)}
    />
  </div>

  <StudentFormModal
    open={vm.formOpen}
    title={vm.formMode === 'create' ? 'Add Student' : 'Edit Student'}
    editStudentId={vm.editStudentId}
    studentIdYear={vm.studentIdYear}
    studentIdSeq={vm.studentIdSeq}
    firstName={vm.firstName}
    lastName={vm.lastName}
    year={vm.year}
    gender={vm.gender}
    collegeOptions={vm.collegeOptions}
    programOptions={vm.programOptions}
    selectedCollegeCode={vm.selectedCollegeCode}
    selectedCollegeLabel={vm.selectedCollegeLabel}
    selectedProgramCode={vm.selectedProgramCode}
    selectedProgramLabel={vm.selectedProgramLabel}
    formError={vm.formError}
    onStudentIdYearInput={(value) => {
      vm.studentIdYear = value
    }}
    onStudentIdSeqInput={(value) => {
      vm.studentIdSeq = value
    }}
    onFirstNameInput={(value) => {
      vm.firstName = value
    }}
    onLastNameInput={(value) => {
      vm.lastName = value
    }}
    onYearInput={(value) => {
      vm.year = value
    }}
    onGenderInput={(value) => {
      vm.gender = value as StudentGender | ''
    }}
    onSelectCollege={(option) => vm.handleCollegeSelect(option)}
    onSelectProgram={(option) => vm.handleProgramSelect(option)}
    onSave={() => void vm.saveStudent()}
    onClose={() => vm.closeForm()}
  />

  <ConfirmDialog
    open={vm.confirmOpen}
    title="Confirm Delete"
    message={vm.deleteTarget ? `Are you sure you want to delete ${vm.deleteTarget.first_name} ${vm.deleteTarget.last_name}? This action cannot be undone.` : ''}
    confirmLabel="Delete"
    onConfirm={() => void vm.remove()}
    onClose={() => vm.closeConfirm()}
  />
</div>