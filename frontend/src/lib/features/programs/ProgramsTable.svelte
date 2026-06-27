<script lang="ts">
  import type { ProgramRow } from '$lib/types/schema'

  import SortHeader from '$lib/components/shared/SortHeader.svelte'
  import TableShell from '$lib/components/shared/TableShell.svelte'
  import RowActions from '$lib/components/shared/RowActions.svelte'
  import ActionsHeader from '$lib/components/shared/ActionsHeader.svelte'

  type CollegeLookup = {
    code: string
    name?: string
  }

  type ProgramWithRefs = ProgramRow & {
    college_code?: string | null
    college?: {
      code?: string | null
    } | null
  }

  type Props = {
    rows: ProgramRow[]
    colleges: CollegeLookup[]
    loading: boolean
    emptyMessage?: string
    sortBy: string
    order: 'ASC' | 'DESC'
    onSort: (column: string) => void
    onEdit: (program: ProgramRow) => void
    onDelete: (program: ProgramRow) => void
  }

  let {
    rows,
    colleges,
    loading,
    emptyMessage = 'No programs found',
    sortBy,
    order,
    onSort,
    onEdit,
    onDelete
  }: Props = $props()

  function hasCollege(code?: string | null) {
    return !!code && colleges.some((college) => college.code === code)
  }

  function displayCollege(program: ProgramWithRefs) {
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
        <SortHeader
          label="Code"
          sortBy={sortBy}
          order={order}
          sortKey="code"
          headerClass="ssis-code-col-header"
          onSort={onSort}
        />
        <SortHeader
          label="Name"
          sortBy={sortBy}
          order={order}
          sortKey="name"
          headerClass="ssis-name-col-header"
          onSort={onSort}
        />
        <SortHeader
          label="College"
          sortBy={sortBy}
          order={order}
          sortKey="college"
          onSort={onSort}
        />
        <ActionsHeader />
      </tr>
    </thead>

    <tbody class="divide-y divide-(--border-subtle)">
      <TableShell loading={loading} hasData={rows.length > 0} colspan={4} {emptyMessage}>
        {#each rows as program}
          {@const collegeDisplay = displayCollege(program as ProgramWithRefs)}

          <tr class="ssis-row transition-colors hover:bg-(--bg-hover)">
            <td class="ssis-cell-primary ssis-code-col-cell px-4 py-3"><span class="ssis-code-pill">{program.code}</span></td>
            <td class="ssis-cell-muted ssis-cell-clip ssis-name-col-cell px-4 py-3">{program.name}</td>
            <td class="ssis-cell-muted px-4 py-3">
              <span
                class="ssis-code-pill"
                class:ssis-code-pill--missing={collegeDisplay.state === 'missing'}
              >{collegeDisplay.label}</span>
            </td>
            <td class="ssis-action-cell px-4 py-3 text-center">
              <RowActions onEdit={() => onEdit(program)} onDelete={() => onDelete(program)} />
            </td>
          </tr>
        {/each}
      </TableShell>
    </tbody>
  </table>
</div>
