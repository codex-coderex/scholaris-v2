<script lang="ts">
  import { Building2, LibraryBig, Plus, Users } from 'lucide-svelte'

  import StudentsSection from '$lib/features/students/StudentsSection.svelte'
  import ProgramsSection from '$lib/features/programs/ProgramsSection.svelte'
  import CollegesSection from '$lib/features/colleges/CollegesSection.svelte'
  import ThemeToggle from '$lib/components/shared/ThemeToggle.svelte'
  import ToastHost from '../lib/components/overlays/ToastHost.svelte'
  import type { AppMode } from '$lib/types/schema'

  type Tab = 'students' | 'programs' | 'colleges'

  let activeTab = $state<Tab>('students')
  let mode = $state<AppMode>('light')
  let studentsSection = $state<any>(null)
  let programsSection = $state<any>(null)
  let collegesSection = $state<any>(null)

  const modules: { key: Tab; label: string; title: string; subtitle: string; icon: typeof Users }[] = [
    {
      key: 'students',
      label: 'Students',
      title: 'Student Records',
      subtitle: 'Manage enrolled student records',
      icon: Users
    },
    {
      key: 'programs',
      label: 'Programs',
      title: 'Program Catalog',
      subtitle: 'Maintain academic programs and assignments',
      icon: LibraryBig
    },
    {
      key: 'colleges',
      label: 'Colleges',
      title: 'College Registry',
      subtitle: 'Manage university colleges',
      icon: Building2
    }
  ]

  const current = $derived(modules.find(m => m.key === activeTab) ?? modules[0])
  const addLabel = $derived(activeTab === 'students' ? 'Add Student' : activeTab === 'programs' ? 'Add Program' : 'Add College')

  function toggleMode() {
    mode = mode === 'light' ? 'dark' : 'light'
  }

  function openActiveAdd() {
    if (activeTab === 'students') {
      studentsSection?.openAdd()
    } else if (activeTab === 'programs') {
      programsSection?.openAdd()
    } else {
      collegesSection?.openAdd()
    }
  }
</script>

<svelte:head>
  <title>Scholaris — {current.title}</title>
</svelte:head>

<div class="ssis-app" data-mode={mode}>
  <div class="ssis-shell">
    <aside class="ssis-sidebar">
      <div class="ssis-brand">
        <div class="font-lora text-[2.1rem] font-semibold leading-none text-(--accent)">Scholaris</div>
        <div class="mt-2 text-[0.82rem] uppercase tracking-[0.12em] text-(--text-muted)">Student Information System</div>
      </div>

      <nav class="ssis-nav mt-1">
        {#each modules as m}
          <button
            type="button"
            onclick={() => activeTab = m.key}
            aria-current={activeTab === m.key ? 'page' : undefined}
            class="ssis-nav-item flex w-full items-center gap-3 text-left {activeTab === m.key ? 'is-active' : ''}"
          >
            <m.icon class="h-4 w-4 shrink-0" />
            <span>{m.label}</span>
          </button>
        {/each}
      </nav>

      <div class="mt-auto border-t border-(--border-subtle) pt-6">
        <ThemeToggle {mode} onToggle={toggleMode} />
      </div>
    </aside>

    <section class="ssis-main">
      <header class="ssis-header">
        <div>
          <h1 class="ssis-title">{current.label}</h1>
          <p class="ssis-subtitle">{current.subtitle}</p>
        </div>

        <div class="header-actions">
          <button
            type="button"
            onclick={openActiveAdd}
            class="btn ssis-btn-primary ssis-add-btn px-4 py-2 text-sm font-semibold normal-case"
          >
            <span class="inline-flex items-center gap-1.5"><Plus class="h-3.5 w-3.5" /> {addLabel}</span>
          </button>
        </div>
      </header>

      <main class="mt-6">
        {#if activeTab === 'students'}
          <StudentsSection bind:this={studentsSection} mode={mode} />
        {:else if activeTab === 'programs'}
          <ProgramsSection bind:this={programsSection} mode={mode} />
        {:else}
          <CollegesSection bind:this={collegesSection} mode={mode} />
        {/if}
      </main>
    </section>
  </div>

  <ToastHost />
</div>