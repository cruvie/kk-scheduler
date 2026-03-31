<script setup lang="ts">
import type {NavigationMenuItem} from '@nuxt/ui'

const items: NavigationMenuItem[] = [{
  label: 'Home',
  icon: 'i-heroicons-home',
  to: '/'
}, {
  label: 'Job List',
  icon: 'i-heroicons-cog-6-tooth',
  to: '/job'
}, {
  label: 'Service List',
  icon: 'i-heroicons-server-stack',
  to: '/service'
}, {
  label: 'Task Execution',
  icon: 'i-heroicons-clipboard-document-list',
  to: '/task_execution'
}]
</script>

<template>
  <UDashboardGroup>
    <UDashboardSidebar>
      <template #header>
        <h1 class="font-bold text-lg">KK Scheduler</h1>
      </template>

      <UNavigationMenu
          :items="items"
          orientation="vertical"
      />
    </UDashboardSidebar>


    <UDashboardPanel
        resizable
        :ui="{ root: 'p-4'}"
    >
      <div class="lg:hidden mb-4">
        <UDashboardSidebarToggle/>
      </div>
      <slot/>
    </UDashboardPanel>
  </UDashboardGroup>

</template>