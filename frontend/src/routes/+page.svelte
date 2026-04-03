<script lang="ts">
  import StudentsTable from '$lib/components/StudentsTable.svelte'
  import ProgramsTable from '$lib/components/ProgramsTable.svelte'
  import CollegesTable from '$lib/components/CollegesTable.svelte'

  let activeTab = $state<'students' | 'programs' | 'colleges'>('students')

  const modules = [
    {
      key: 'students' as const,
      label: 'Students',
      title: 'Student records',
      description: 'Browse enrolled students, sort by fields, and manage profiles.'
    },
    {
      key: 'programs' as const,
      label: 'Programs',
      title: 'Program catalog',
      description: 'Maintain degree programs and their college assignments.'
    },
    {
      key: 'colleges' as const,
      label: 'Colleges',
      title: 'College registry',
      description: 'Organize college metadata for the academic structure.'
    }
  ]

  function selectModule(key: 'students' | 'programs' | 'colleges') {
    activeTab = key
  }

  const currentModule = $derived(modules.find((module) => module.key === activeTab) ?? modules[0])
</script>

<svelte:head>
  <title>Scholaris V2</title>
</svelte:head>

<div class="min-h-screen bg-slate-100 text-slate-900">
  <main class="mx-auto flex min-h-screen w-full max-w-6xl flex-row gap-4 p-4 sm:p-6">
    <aside class="w-64 shrink-0 rounded border border-slate-300 bg-white p-4">
      <div class="rounded border border-slate-200 bg-slate-50 p-3">
        <div class="text-sm font-semibold text-slate-900">SSIS</div>
        <div class="text-sm text-slate-600">MSU-IIT</div>
      </div>

      <nav class="mt-4 space-y-2 text-sm">
        {#each modules as module}
          <button
            type="button"
            class="flex w-full items-start gap-2 rounded border px-3 py-2 text-left transition {activeTab === module.key ? 'border-slate-900 bg-slate-900 text-white' : 'border-slate-300 bg-white text-slate-700 hover:bg-slate-50'}"
            aria-current={activeTab === module.key ? 'page' : undefined}
            onclick={() => selectModule(module.key)}
          >
            <span>
              <span class="block font-medium">{module.label}</span>
              <span class="block text-xs leading-5 {activeTab === module.key ? 'text-slate-200' : 'text-slate-500'}">{module.description}</span>
            </span>
          </button>
        {/each}
      </nav>
    </aside>

    <section class="min-w-0 flex-1 rounded border border-slate-300 bg-white">
      <header class="border-b border-slate-200 px-5 py-4 sm:px-6">
        <div class="flex flex-col gap-3 lg:flex-row lg:items-end lg:justify-between">
          <div>
            <div class="text-xs uppercase tracking-[0.3em] text-slate-500">Scholaris V2</div>
            <h1 class="mt-1 text-2xl font-semibold text-slate-900 sm:text-3xl">{currentModule.title}</h1>
            <p class="mt-2 max-w-2xl text-sm leading-6 text-slate-600">{currentModule.description}</p>
          </div>
        </div>
      </header>

      <div class="p-4 sm:p-5 lg:p-6">
        <div class="mt-4">
          {#if activeTab === 'students'}
            <StudentsTable />
          {:else if activeTab === 'programs'}
            <ProgramsTable />
          {:else}
            <CollegesTable />
          {/if}
        </div>
      </div>
    </section>
  </main>
</div>