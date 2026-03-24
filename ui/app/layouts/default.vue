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
}]
</script>

<template>
  <UDashboardGroup>
    <UDashboardSidebar>
      <template #header>
        <h1 class="font-bold text-lg">KK Schedule</h1>
      </template>

      <template #default>
        <UNavigationMenu
            :items="items"
            orientation="vertical"
        />
      </template>
    </UDashboardSidebar>


    <UDashboardPanel
        resizable
        :ui="{ root: 'p-4'}"
    >
      <slot/>
    </UDashboardPanel>
  </UDashboardGroup>

</template>