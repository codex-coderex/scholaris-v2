<script lang="ts">
  import StudentsTable from '$lib/components/StudentsTable.svelte'
  import ProgramsTable from '$lib/components/ProgramsTable.svelte'
  import CollegesTable from '$lib/components/CollegesTable.svelte'

  let activeTab = $state<'students' | 'programs' | 'colleges'>('students')
</script>

<div class="min-h-screen bg-base-200">

  <!-- navbar -->
  <div class="navbar bg-base-100 shadow px-6">
    <span class="text-xl font-bold">Scholaris V2</span>
  </div>

  <!-- tabs -->
  <div class="px-6 pt-4">
    <div role="tablist" class="tabs tabs-bordered">
      <button
        role="tab"
        class="tab {activeTab === 'students' ? 'tab-active' : ''}"
        onclick={() => activeTab = 'students'}
      >
        Students
      </button>
      <button
        role="tab"
        class="tab {activeTab === 'programs' ? 'tab-active' : ''}"
        onclick={() => activeTab = 'programs'}
      >
        Programs
      </button>
      <button
        role="tab"
        class="tab {activeTab === 'colleges' ? 'tab-active' : ''}"
        onclick={() => activeTab = 'colleges'}
      >
        Colleges
      </button>
    </div>
  </div>

  <!-- content -->
  <div class="px-6 pt-4">
    {#if activeTab === 'students'}
      <StudentsTable />
    {:else if activeTab === 'programs'}
      <ProgramsTable />
    {:else}
      <CollegesTable />
    {/if}
  </div>

</div>