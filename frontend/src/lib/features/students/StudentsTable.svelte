<script lang="ts">
  import type { StudentRow } from '$lib/types/schema'

  import SortHeader from '$lib/components/shared/SortHeader.svelte'
  import TableShell from '$lib/components/shared/TableShell.svelte'
  import RowActions from '$lib/components/shared/RowActions.svelte'
  import ActionsHeader from '$lib/components/shared/ActionsHeader.svelte'

  type CollegeLookup = {
    code: string
    name?: string
  }

  type StudentWithRefs = StudentRow & {
    program_code?: string | null
    program?: {
      code?: string | null
      college_code?: string | null
      college?: {
        code?: string | null
      } | null
    } | null
  }

  type Props = {
    rows: StudentRow[]
    colleges: CollegeLookup[]
    loading: boolean
    emptyMessage?: string
    sortBy: string
    order: 'ASC' | 'DESC'
    onSort: (column: string) => void
    onEdit: (student: StudentRow) => void
    onDelete: (student: StudentRow) => void
  }

  let {
    rows,
    colleges,
    loading,
    emptyMessage = 'No students found',
    sortBy,
    order,
    onSort,
    onEdit,
    onDelete
  }: Props = $props()

  function hasCollege(code?: string | null) {
    return !!code && colleges.some((college) => college.code === code)
  }

  function rawProgramCode(student: StudentWithRefs) {
    return student.program_code ?? student.program?.code ?? ''
  }

  function displayProgram(student: StudentWithRefs) {
    const code = rawProgramCode(student)

    if (!code) {
      return { label: 'N/A', state: 'unenrolled' }
    }

    if (!student.program?.code) {
      return { label: 'N/A', state: 'missing' }
    }

    return { label: student.program.code, state: 'ok' }
  }

  function displayCollege(student: StudentWithRefs) {
    const program = student.program

    if (!rawProgramCode(student) || !program?.code) {
      return { label: 'N/A', state: 'missing' }
    }

    const code = program.college?.code ?? program.college_code ?? ''

    if (!hasCollege(code)) {
      return { label: 'N/A', state: 'missing' }
    }

    return { label: code, state: 'ok' }
  }
</script>

<div class="ssis-table-wrap overflow-x-auto">
  <table class="table ssis-data-table w-full">
    <thead>
      <tr class="ssis-thead border-b border-(--border)">
        <SortHeader label="ID" sortBy={sortBy} order={order} sortKey="id" onSort={onSort} />
        <SortHeader label="First Name" sortBy={sortBy} order={order} sortKey="first_name" onSort={onSort} />
        <SortHeader label="Last Name" sortBy={sortBy} order={order} sortKey="last_name" onSort={onSort} />
        <SortHeader label="Program" sortBy={sortBy} order={order} sortKey="program" onSort={onSort} />
        <SortHeader label="College" sortBy={sortBy} order={order} sortKey="college" onSort={onSort} />
        <SortHeader label="Year" sortBy={sortBy} order={order} sortKey="year" onSort={onSort} />
        <SortHeader label="Gender" sortBy={sortBy} order={order} sortKey="gender" onSort={onSort} />
        <ActionsHeader />
      </tr>
    </thead>

    <tbody class="divide-y divide-(--border-subtle)">
      <TableShell loading={loading} hasData={rows.length > 0} colspan={8} {emptyMessage}>
        {#each rows as student}
          {@const programDisplay = displayProgram(student as StudentWithRefs)}
          {@const collegeDisplay = displayCollege(student as StudentWithRefs)}

          <tr class="ssis-row transition-colors hover:bg-(--bg-hover)">
            <td class="ssis-cell-primary px-4 py-3.5"><span class="ssis-code-pill">{student.id}</span></td>
            <td class="ssis-cell-muted ssis-cell-clip px-4 py-3.5">{student.first_name}</td>
            <td class="ssis-cell-muted ssis-cell-clip px-4 py-3.5">{student.last_name}</td>
            <td class="ssis-cell-muted px-4 py-3.5">
              <span
                class="ssis-code-pill"
                class:ssis-code-pill--missing={programDisplay.state === 'missing'}
              >{programDisplay.label}</span>
            </td>
            <td class="ssis-cell-muted px-4 py-3.5">
              <span
                class="ssis-code-pill"
                class:ssis-code-pill--missing={collegeDisplay.state === 'missing'}
              >{collegeDisplay.label}</span>
            </td>
            <td class="ssis-cell-muted px-4 py-3.5">{student.year}</td>
            <td class="ssis-cell-muted px-4 py-3.5">{student.gender}</td>
            <td class="ssis-action-cell px-4 py-3.5 text-center">
              <RowActions onEdit={() => onEdit(student)} onDelete={() => onDelete(student)} />
            </td>
          </tr>
        {/each}
      </TableShell>
    </tbody>
  </table>
</div>
