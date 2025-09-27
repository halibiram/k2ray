<template>
  <div class="bg-white shadow-md rounded-lg p-6">
    <h2 class="text-xl font-semibold text-gray-800 border-b pb-2">System Logs</h2>
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
