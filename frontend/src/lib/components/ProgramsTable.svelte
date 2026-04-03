<script lang="ts">
  import {
    GetPrograms,
    DeleteProgram,
    GetColleges
  } from '../../../bindings/scholaris-v2/appservice'
  
  import { Program as ProgramModel } from '../../../bindings/scholaris-v2/internal/models/models.js'
  import SearchToolbar from './shared/SearchToolbar.svelte'
  import SortHeader from './shared/SortHeader.svelte'
  import Pagination from './shared/Pagination.svelte'
  import TableShell from './shared/TableShell.svelte'

  // state

  let programs    = $state<any[]>([])
  let colleges    = $state<any[]>([])
  let total       = $state(0)
  let search      = $state('')
  let sortBy      = $state('p.code')
  let order       = $state('ASC')
  let page        = $state(1)
  let pageSize    = $state(10)
  let loading     = $state(false)
  let error       = $state('')

  // derived

  let totalPages  = $derived(Math.ceil(total / pageSize))

  // load

  async function load() {
    loading = true
    error   = ''
    try {
      const [data, count] = await GetPrograms(search, sortBy, order, page, pageSize)
      programs = data ?? []
      total    = count
    } catch (e: any) {
      error = e.message ?? 'Failed to load programs'
    } finally {
      loading = false
    }
  }

  async function loadColleges() {
    try {
      const [data] = await GetColleges('', 'code', 'ASC', 1, 100)
      colleges = data ?? []
    } catch {}
  }

  // sort

  function setSort(col: string) {
    if (sortBy === col) {
      order = order === 'ASC' ? 'DESC' : 'ASC'
    } else {
      sortBy = col
      order  = 'ASC'
    }
    page = 1
    load()
  }

  function handleSearch(value: string) {
    search = value
    page = 1
    load()
  }

  function previousPage() {
    page--
    load()
  }

  function nextPage() {
    page++
    load()
  }

  async function remove(code: string) {
    if (!confirm(`Delete program ${code}?`)) return
    try {
      await DeleteProgram(code)
      load()
    } catch (e: any) {
      error = e.message ?? 'Failed to delete'
    }
  }

  // init

  $effect(() => {
    load()
    loadColleges()
  })
</script>

<div class="space-y-4">

    <SearchToolbar
      search={search}
      placeholder="Search programs..."
      onSearch={handleSearch}
    />

    <!-- error -->
    {#if error}
      <div class="alert alert-error mb-4">{error}</div>
    {/if}

    <!-- table -->
    <div class="overflow-x-auto">
      <table class="table w-full border border-slate-300 text-sm">
        <thead>
          <tr>
            <SortHeader label="Code" sortBy={sortBy} order={order} sortKey="p.code" onSort={setSort} />
            <SortHeader label="Name" sortBy={sortBy} order={order} sortKey="p.name" onSort={setSort} />
            <SortHeader label="College" sortBy={sortBy} order={order} sortKey="c.name" onSort={setSort} />
            <th class="px-3 py-2 text-left text-xs font-medium text-slate-600">Actions</th>
          </tr>
        </thead>
        <tbody>
          <TableShell loading={loading} hasData={programs.length > 0} colspan={4} emptyMessage="No programs found">
            {#each programs as p}
              <tr>
                <td>{p.code}</td>
                <td>{p.name}</td>
                <td>
                  {p.college?.code} {p.college?.name}
                </td>
                <td>
                  <button class="rounded border border-slate-300 px-2 py-1 text-xs hover:bg-slate-50" onclick={() => remove(p.code)}>Delete</button>
                </td>
              </tr>
            {/each}
          </TableShell>
        </tbody>
      </table>
    </div>

    <!-- pagination -->
    <Pagination
      page={page}
      totalPages={totalPages}
      total={total}
      onPrevious={previousPage}
      onNext={nextPage}
    />

</div>
