<script lang="ts">
  import type { CollegeRow } from '$lib/types/schema'

  import SortHeader from '$lib/components/shared/SortHeader.svelte'
  import TableShell from '$lib/components/shared/TableShell.svelte'
  import RowActions from '$lib/components/shared/RowActions.svelte'
  import ActionsHeader from '$lib/components/shared/ActionsHeader.svelte'

  type Props = {
    rows: CollegeRow[]
    loading: boolean
    emptyMessage?: string
    sortBy: string
    order: 'ASC' | 'DESC'
    onSort: (column: string) => void
    onEdit: (college: CollegeRow) => void
    onDelete: (college: CollegeRow) => void
  }

  let {
    rows,
    loading,
    emptyMessage = 'No colleges found',
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
        <ActionsHeader />
      </tr>
    </thead>

    <tbody class="divide-y divide-(--border-subtle)">
      <TableShell loading={loading} hasData={rows.length > 0} colspan={3} {emptyMessage}>
        {#each rows as college}
          <tr class="ssis-row transition-colors hover:bg-(--bg-hover)">
            <td class="ssis-cell-primary ssis-code-col-cell px-4 py-3"><span class="ssis-code-pill">{college.code}</span></td>
            <td class="ssis-cell-muted ssis-cell-clip ssis-name-col-cell px-4 py-3">{college.name}</td>
            <td class="ssis-action-cell px-4 py-3 text-center">
              <RowActions onEdit={() => onEdit(college)} onDelete={() => onDelete(college)} />
            </td>
          </tr>
        {/each}
      </TableShell>
    </tbody>
  </table>
</div>
