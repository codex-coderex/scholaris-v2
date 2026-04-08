<script lang="ts">
  import type { AppMode, StudentRow } from '$lib/types/schema'

  import SortHeader from '$lib/components/shared/SortHeader.svelte'
  import TableShell from '$lib/components/shared/TableShell.svelte'
  import RowActions from '$lib/components/shared/RowActions.svelte'
  import ActionsHeader from '$lib/components/shared/ActionsHeader.svelte'

  type Props = {
    mode?: AppMode
    rows: StudentRow[]
    loading: boolean
    sortBy: string
    order: 'ASC' | 'DESC'
    onSort: (column: string) => void
    onEdit: (student: StudentRow) => void
    onDelete: (student: StudentRow) => void
  }

  let {
    mode = 'light',
    rows,
    loading,
    sortBy,
    order,
    onSort,
    onEdit,
    onDelete
  }: Props = $props()
</script>

<div class="ssis-table-wrap overflow-x-auto">
  <table class="table ssis-data-table w-full">
    <thead>
      <tr class="ssis-thead border-b border-(--border)">
        <SortHeader {mode} label="ID" sortBy={sortBy} order={order} sortKey="s.id" onSort={onSort} />
        <SortHeader {mode} label="First Name" sortBy={sortBy} order={order} sortKey="s.first_name" onSort={onSort} />
        <SortHeader {mode} label="Last Name" sortBy={sortBy} order={order} sortKey="s.last_name" onSort={onSort} />
        <SortHeader {mode} label="Program" sortBy={sortBy} order={order} sortKey="p.code" onSort={onSort} />
        <SortHeader {mode} label="Year" sortBy={sortBy} order={order} sortKey="s.year" onSort={onSort} />
        <SortHeader {mode} label="Gender" sortBy={sortBy} order={order} sortKey="s.gender" onSort={onSort} />
        <ActionsHeader />
      </tr>
    </thead>

    <tbody class="divide-y divide-(--border-subtle)">
      <TableShell loading={loading} hasData={rows.length > 0} colspan={7} emptyMessage="No students found">
        {#each rows as student}
          <tr class="ssis-row transition-colors hover:bg-(--bg-hover)">
            <td class="ssis-cell-primary px-4 py-3.5"><span class="ssis-code-pill">{student.id}</span></td>
            <td class="ssis-cell-muted ssis-cell-clip px-4 py-3.5">{student.first_name}</td>
            <td class="ssis-cell-muted ssis-cell-clip px-4 py-3.5">{student.last_name}</td>
            <td class="ssis-cell-muted px-4 py-3.5"><span class="ssis-code-pill">{student.program?.code}</span></td>
            <td class="ssis-cell-muted px-4 py-3.5">{student.year}</td>
            <td class="ssis-cell-muted px-4 py-3.5">{student.gender}</td>
            <td class="ssis-action-cell px-4 py-3.5 text-center">
              <RowActions mode={mode} onEdit={() => onEdit(student)} onDelete={() => onDelete(student)} />
            </td>
          </tr>
        {/each}
      </TableShell>
    </tbody>
  </table>
</div>
