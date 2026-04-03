<script lang="ts">
  import {
    GetColleges,
    CreateCollege,
    UpdateCollege,
    DeleteCollege
  } from '../../../bindings/scholaris-v2/appservice'

  // state

  let colleges   = $state<any[]>([])
  let total      = $state(0)
  let search     = $state('')
  let sortBy     = $state('code')
  let order      = $state('ASC')
  let page       = $state(1)
  let pageSize   = $state(10)
  let loading    = $state(false)
  let error      = $state('')

  // modal state

  let showModal  = $state(false)
  let editMode   = $state(false)
  let form       = $state({ code: '', name: '' })

  // derived

  let totalPages = $derived(Math.ceil(total / pageSize))

  // load

  async function load() {
    loading = true
    error   = ''
    try {
      const [data, count] = await GetColleges(search, sortBy, order, page, pageSize)
      colleges = data ?? []
      total    = count
    } catch (e: any) {
      error = e.message ?? 'Failed to load colleges'
    } finally {
      loading = false
    }
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
    form     = { code: '', name: '' }
    editMode = false
    showModal = true
  }

  function openEdit(c: any) {
    form     = { code: c.code, name: c.name }
    editMode = true
    showModal = true
  }

  async function save() {
    try {
      if (editMode) {
        await UpdateCollege(form)
      } else {
        await CreateCollege(form)
      }
      showModal = false
      load()
    } catch (e: any) {
      error = e.message ?? 'Failed to save'
    }
  }

  async function remove(code: string) {
    if (!confirm(`Delete college ${code}?`)) return
    try {
      await DeleteCollege(code)
      load()
    } catch (e: any) {
      error = e.message ?? 'Failed to delete'
    }
  }

  // init

  $effect(() => { load() })
</script>

<div class="card bg-base-100 shadow">
  <div class="card-body">

    <!-- toolbar -->
    <div class="flex justify-between items-center mb-4">
      <input
        class="input input-bordered w-72"
        placeholder="Search colleges..."
        bind:value={search}
        oninput={() => { page = 1; load() }}
      />
      <button class="btn btn-primary" onclick={openCreate}>
        + Add College
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
              <button onclick={() => setSort('code')}>
                Code {sortBy === 'code' ? (order === 'ASC' ? '↑' : '↓') : ''}
              </button>
            </th>
            <th>
              <button onclick={() => setSort('name')}>
                Name {sortBy === 'name' ? (order === 'ASC' ? '↑' : '↓') : ''}
              </button>
            </th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {#if loading}
            <tr><td colspan="3" class="text-center py-8">Loading...</td></tr>
          {:else if colleges.length === 0}
            <tr><td colspan="3" class="text-center py-8">No colleges found</td></tr>
          {:else}
            {#each colleges as c}
              <tr>
                <td>{c.code}</td>
                <td>{c.name}</td>
                <td class="flex gap-2">
                  <button class="btn btn-sm btn-outline" onclick={() => openEdit(c)}>Edit</button>
                  <button class="btn btn-sm btn-error btn-outline" onclick={() => remove(c.code)}>Delete</button>
                </td>
              </tr>
            {/each}
          {/if}
        </tbody>
      </table>
    </div>

    <!-- pagination -->
    <div class="flex justify-between items-center mt-4">
      <span class="text-sm text-base-content/60">
        {total} total
      </span>
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
        {editMode ? 'Edit College' : 'Add College'}
      </h3>

      <fieldset class="fieldset">
        <label class="fieldset-label" for="college-code">Code</label>
        <input
          id="college-code"
          class="input input-bordered w-full"
          bind:value={form.code}
          disabled={editMode}
          placeholder="e.g. CCS"
        />
      </fieldset>

      <fieldset class="fieldset mt-2">
        <label class="fieldset-label" for="college-name">Name</label>
        <input
          id="college-name"
          class="input input-bordered w-full"
          bind:value={form.name}
          placeholder="e.g. College of Computer Studies"
        />
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