<script setup lang="ts">
import type { NavigationMenuItem } from '@nuxt/ui'

const open = ref(true)

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
  <div class="flex flex-1 h-screen">
    <USidebar
      v-model:open="open"
      collapsible="icon"
      :ui="{
        container: 'h-full',
        inner: 'bg-elevated/25',
        body: 'py-2'
      }"
    >
      <template #header>
        <h1 class="font-bold text-lg truncate px-2">KK Scheduler</h1>
      </template>

      <template #default="{ state }">
        <UNavigationMenu
          :items="items"
          orientation="vertical"
          :ui="{ link: 'p-1.5 overflow-hidden' }"
        />
      </template>

      <template #footer>
        <UButton
          :icon="open ? 'i-heroicons-chevron-left' : 'i-heroicons-chevron-right'"
          color="neutral"
          variant="ghost"
          square
          class="w-full"
          @click="open = !open"
        />
      </template>
    </USidebar>

    <div class="flex-1 flex flex-col min-w-0">
      <div class="h-(--ui-header-height) shrink-0 flex items-center px-4 border-b border-default">
        <UButton
          v-if="!open"
          icon="i-heroicons-bars-3"
          color="neutral"
          variant="ghost"
          aria-label="Expand sidebar"
          @click="open = true"
        />
      </div>

      <div class="flex-1 p-4 overflow-auto">
        <slot />
      </div>
    </div>
  </div>
</template>
