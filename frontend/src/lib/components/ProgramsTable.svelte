<script lang="ts">
  import {
    GetPrograms,
    CreateProgram,
    UpdateProgram,
    DeleteProgram,
    GetColleges
  } from '../../../bindings/scholaris-v2/appservice'
  import { Program as ProgramModel } from '../../../bindings/scholaris-v2/internal/models/models.js'

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

  // modal state

  let showModal   = $state(false)
  let editMode    = $state(false)

  function createEmptyForm() {
    return new ProgramModel({
      code: '',
      name: '',
      college_code: ''
    })
  }

  let form        = $state(createEmptyForm())

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

  // modal

  function openCreate() {
    form      = createEmptyForm()
    editMode  = false
    showModal = true
  }

  function openEdit(p: any) {
    form      = new ProgramModel(p)
    editMode  = true
    showModal = true
  }

  async function save() {
    try {
      if (editMode) {
        await UpdateProgram(form)
      } else {
        await CreateProgram(form)
      }
      showModal = false
      load()
    } catch (e: any) {
      error = e.message ?? 'Failed to save'
    }
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

<div class="card bg-base-100 shadow">
  <div class="card-body">

    <!-- toolbar -->
    <div class="flex justify-between items-center mb-4">
      <input
        class="input input-bordered w-72"
        placeholder="Search programs..."
        bind:value={search}
        oninput={() => { page = 1; load() }}
      />
      <button class="btn btn-primary" onclick={openCreate}>
        + Add Program
      </button>
    </div>

    <!-- error -->
    {#if error}
      <div class="alert alert-error mb-4">{error}</div>
    {/if}

    <!-- table -->
    <div class="overflow-x-auto">
      <table class="table table-zebra w-full">
        <thead>
          <tr>
            <th>
              <button onclick={() => setSort('p.code')}>
                Code {sortBy === 'p.code' ? (order === 'ASC' ? '↑' : '↓') : ''}
              </button>
            </th>
            <th>
              <button onclick={() => setSort('p.name')}>
                Name {sortBy === 'p.name' ? (order === 'ASC' ? '↑' : '↓') : ''}
              </button>
            </th>
            <th>
              <button onclick={() => setSort('c.name')}>
                College {sortBy === 'c.name' ? (order === 'ASC' ? '↑' : '↓') : ''}
              </button>
            </th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {#if loading}
            <tr><td colspan="4" class="text-center py-8">Loading...</td></tr>
          {:else if programs.length === 0}
            <tr><td colspan="4" class="text-center py-8">No programs found</td></tr>
          {:else}
            {#each programs as p}
              <tr>
                <td>{p.code}</td>
                <td>{p.name}</td>
                <td>
                  <span class="badge badge-ghost">{p.college?.code}</span>
                  {p.college?.name}
                </td>
                <td class="flex gap-2">
                  <button class="btn btn-sm btn-outline" onclick={() => openEdit(p)}>Edit</button>
                  <button class="btn btn-sm btn-error btn-outline" onclick={() => remove(p.code)}>Delete</button>
                </td>
              </tr>
            {/each}
          {/if}
        </tbody>
      </table>
    </div>

    <!-- pagination -->
    <div class="flex justify-between items-center mt-4">
      <span class="text-sm text-base-content/60">{total} total</span>
      <div class="join">
        <button
          class="join-item btn btn-sm"
          disabled={page === 1}
          onclick={() => { page--; load() }}
        >«</button>
        <button class="join-item btn btn-sm">{page} / {totalPages}</button>
        <button
          class="join-item btn btn-sm"
          disabled={page >= totalPages}
          onclick={() => { page++; load() }}
        >»</button>
      </div>
    </div>

  </div>
</div>

<!-- modal -->
{#if showModal}
  <dialog class="modal modal-open">
    <div class="modal-box">
      <h3 class="font-bold text-lg mb-4">
        {editMode ? 'Edit Program' : 'Add Program'}
      </h3>

      <fieldset class="fieldset">
        <label class="fieldset-label" for="program-code">Code</label>
        <input
          id="program-code"
          class="input input-bordered w-full"
          bind:value={form.code}
          disabled={editMode}
          placeholder="e.g. BSCS"
        />
      </fieldset>

      <fieldset class="fieldset mt-2">
        <label class="fieldset-label" for="program-name">Name</label>
        <input
          id="program-name"
          class="input input-bordered w-full"
          bind:value={form.name}
          placeholder="e.g. Bachelor of Science in Computer Science"
        />
      </fieldset>

      <fieldset class="fieldset mt-2">
        <label class="fieldset-label" for="program-college">College</label>
        <select id="program-college" class="select select-bordered w-full" bind:value={form.college_code}>
          <option value="">Select college</option>
          {#each colleges as c}
            <option value={c.code}>{c.code} — {c.name}</option>
          {/each}
        </select>
      </fieldset>

      <div class="modal-action">
        <button class="btn" onclick={() => showModal = false}>Cancel</button>
        <button class="btn btn-primary" onclick={save}>Save</button>
      </div>
    </div>
    <button
      type="button"
      class="modal-backdrop"
      aria-label="Close modal"
      onclick={() => showModal = false}
    ></button>
  </dialog>
{/if}