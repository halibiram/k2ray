<template>
  <div class="bg-white dark:bg-gray-800 shadow-md rounded-lg p-6">
    <h2 class="text-xl font-semibold text-gray-800 dark:text-gray-100 border-b border-gray-200 dark:border-gray-700 pb-2">System Logs</h2>
    <div class="mt-4 w-full h-96 bg-gray-900 text-white font-mono text-sm rounded-md p-4 overflow-auto">
      <pre>{{ logs || 'Loading logs...' }}</pre>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  logs: string
}>()
</script>
