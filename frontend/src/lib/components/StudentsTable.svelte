<script lang="ts">
  import {
    GetStudents,
    DeleteStudent
  } from '../../../bindings/scholaris-v2/appservice'
  
  import { GetPrograms } from '../../../bindings/scholaris-v2/appservice'
  import SearchToolbar from './shared/SearchToolbar.svelte'
  import SortHeader from './shared/SortHeader.svelte'
  import Pagination from './shared/Pagination.svelte'
  import TableShell from './shared/TableShell.svelte'

  // state

  let students   = $state<any[]>([])
  let programs   = $state<any[]>([])
  let total      = $state(0)
  let search     = $state('')
  let sortBy     = $state('s.id')
  let order      = $state('ASC')
  let page       = $state(1)
  let pageSize   = $state(10)
  let loading    = $state(false)
  let error      = $state('')

  // derived

  let totalPages = $derived(Math.ceil(total / pageSize))

  // load

  async function load() {
    loading = true
    error   = ''
    try {
      const [data, count] = await GetStudents(search, sortBy, order, page, pageSize)
      students = data ?? []
      total    = count
    } catch (e: any) {
      error = e.message ?? 'Failed to load students'
    } finally {
      loading = false
    }
  }

  async function loadPrograms() {
    try {
      const [data] = await GetPrograms('', 'p.code', 'ASC', 1, 200)
      programs = data ?? []
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

  async function remove(id: string) {
    if (!confirm(`Delete student ${id}?`)) return
    try {
      await DeleteStudent(id)
      load()
    } catch (e: any) {
      error = e.message ?? 'Failed to delete'
    }
  }

  // init

  $effect(() => {
    load()
    loadPrograms()
  })
</script>

<div class="space-y-4">

    <SearchToolbar
      search={search}
      placeholder="Search students..."
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
            <SortHeader label="ID" sortBy={sortBy} order={order} sortKey="s.id" onSort={setSort} />
            <SortHeader label="First Name" sortBy={sortBy} order={order} sortKey="s.first_name" onSort={setSort} />
            <SortHeader label="Last Name" sortBy={sortBy} order={order} sortKey="s.last_name" onSort={setSort} />
            <th class="px-3 py-2 text-left text-xs font-medium text-slate-600">Year</th>
            <th class="px-3 py-2 text-left text-xs font-medium text-slate-600">Gender</th>
            <SortHeader label="Program" sortBy={sortBy} order={order} sortKey="p.code" onSort={setSort} />
            <th class="px-3 py-2 text-left text-xs font-medium text-slate-600">College</th>
            <th class="px-3 py-2 text-left text-xs font-medium text-slate-600">Actions</th>
          </tr>
        </thead>
        <tbody>
          <TableShell loading={loading} hasData={students.length > 0} colspan={8} emptyMessage="No students found">
            {#each students as s}
              <tr>
                <td class="font-mono">{s.id}</td>
                <td>{s.first_name}</td>
                <td>{s.last_name}</td>
                <td>{s.year}</td>
                <td>{s.gender}</td>
                <td>
                  {s.program?.code}
                </td>
                <td>
                  {s.program?.college?.code}
                </td>
                <td>
                  <button class="rounded border border-slate-300 px-2 py-1 text-xs hover:bg-slate-50" onclick={() => remove(s.id)}>Delete</button>
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
