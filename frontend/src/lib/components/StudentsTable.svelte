<script lang="ts">
  import {
    GetStudents,
    CreateStudent,
    UpdateStudent,
    DeleteStudent
  } from '../../../bindings/scholaris-v2/appservice'
  import { GetPrograms } from '../../../bindings/scholaris-v2/appservice'
  import { Student as StudentModel } from '../../../bindings/scholaris-v2/internal/models/models.js'

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

  // modal state

  let showModal  = $state(false)
  let editMode   = $state(false)

  function createEmptyForm() {
    return new StudentModel({
      id: '',
      first_name: '',
      last_name: '',
      year: 1,
      gender: 'Male',
      program_code: ''
    })
  }

  let form = $state(createEmptyForm())

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

  // modal

  function openCreate() {
    form = createEmptyForm()
    editMode  = false
    showModal = true
  }

  function openEdit(s: any) {
    form = new StudentModel(s)
    editMode  = true
    showModal = true
  }

  async function save() {
    try {
      if (editMode) {
        await UpdateStudent(form)
      } else {
        await CreateStudent(form)
      }
      showModal = false
      load()
    } catch (e: any) {
      error = e.message ?? 'Failed to save'
    }
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

<div class="card bg-base-100 shadow">
  <div class="card-body">

    <!-- toolbar -->
    <div class="flex justify-between items-center mb-4">
      <input
        class="input input-bordered w-72"
        placeholder="Search students..."
        bind:value={search}
        oninput={() => { page = 1; load() }}
      />
      <button class="btn btn-primary" onclick={openCreate}>
        + Add Student
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
              <button onclick={() => setSort('s.id')}>
                ID {sortBy === 's.id' ? (order === 'ASC' ? '↑' : '↓') : ''}
              </button>
            </th>
            <th>
              <button onclick={() => setSort('s.first_name')}>
                First Name {sortBy === 's.first_name' ? (order === 'ASC' ? '↑' : '↓') : ''}
              </button>
            </th>
            <th>
              <button onclick={() => setSort('s.last_name')}>
                Last Name {sortBy === 's.last_name' ? (order === 'ASC' ? '↑' : '↓') : ''}
              </button>
            </th>
            <th>Year</th>
            <th>Gender</th>
            <th>
              <button onclick={() => setSort('p.code')}>
                Program {sortBy === 'p.code' ? (order === 'ASC' ? '↑' : '↓') : ''}
              </button>
            </th>
            <th>College</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {#if loading}
            <tr><td colspan="8" class="text-center py-8">Loading...</td></tr>
          {:else if students.length === 0}
            <tr><td colspan="8" class="text-center py-8">No students found</td></tr>
          {:else}
            {#each students as s}
              <tr>
                <td class="font-mono">{s.id}</td>
                <td>{s.first_name}</td>
                <td>{s.last_name}</td>
                <td>{s.year}</td>
                <td>{s.gender}</td>
                <td>
                  <span class="badge badge-ghost">{s.program?.code}</span>
                </td>
                <td>
                  <span class="badge badge-outline">{s.program?.college?.code}</span>
                </td>
                <td class="flex gap-2">
                  <button class="btn btn-sm btn-outline" onclick={() => openEdit(s)}>Edit</button>
                  <button class="btn btn-sm btn-error btn-outline" onclick={() => remove(s.id)}>Delete</button>
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
        {editMode ? 'Edit Student' : 'Add Student'}
      </h3>

      <fieldset class="fieldset">
        <label class="fieldset-label" for="student-id">ID</label>
        <input
          id="student-id"
          class="input input-bordered w-full"
          bind:value={form.id}
          disabled={editMode}
          placeholder="e.g. 2024-0001"
        />
      </fieldset>

      <div class="grid grid-cols-2 gap-2 mt-2">
        <fieldset class="fieldset">
          <label class="fieldset-label" for="student-first-name">First Name</label>
          <input
            id="student-first-name"
            class="input input-bordered w-full"
            bind:value={form.first_name}
            placeholder="First name"
          />
        </fieldset>

        <fieldset class="fieldset">
          <label class="fieldset-label" for="student-last-name">Last Name</label>
          <input
            id="student-last-name"
            class="input input-bordered w-full"
            bind:value={form.last_name}
            placeholder="Last name"
          />
        </fieldset>
      </div>

      <div class="grid grid-cols-2 gap-2 mt-2">
        <fieldset class="fieldset">
          <label class="fieldset-label" for="student-year">Year</label>
          <select id="student-year" class="select select-bordered w-full" bind:value={form.year}>
            <option value={1}>1st Year</option>
            <option value={2}>2nd Year</option>
            <option value={3}>3rd Year</option>
            <option value={4}>4th Year</option>
          </select>
        </fieldset>

        <fieldset class="fieldset">
          <label class="fieldset-label" for="student-gender">Gender</label>
          <select id="student-gender" class="select select-bordered w-full" bind:value={form.gender}>
            <option value="Male">Male</option>
            <option value="Female">Female</option>
          </select>
        </fieldset>
      </div>

      <fieldset class="fieldset mt-2">
        <label class="fieldset-label" for="student-program">Program</label>
        <select id="student-program" class="select select-bordered w-full" bind:value={form.program_code}>
          <option value="">Select program</option>
          {#each programs as p}
            <option value={p.code}>{p.code} — {p.name}</option>
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