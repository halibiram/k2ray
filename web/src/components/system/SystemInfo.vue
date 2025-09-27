<template>
  <div class="bg-white dark:bg-gray-800 shadow-md rounded-lg p-6">
    <h2 class="text-xl font-semibold text-gray-800 dark:text-gray-100 border-b border-gray-200 dark:border-gray-700 pb-2">System Information</h2>
    <div v-if="info" class="mt-4 grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
      <div><strong class="font-medium text-gray-600 dark:text-gray-400">Hostname:</strong> <span class="text-gray-800 dark:text-gray-200">{{ info.hostname }}</span></div>
      <div><strong class="font-medium text-gray-600 dark:text-gray-400">Keenetic Model:</strong> <span class="text-gray-800 dark:text-gray-200">{{ info.keenetic_model }}</span></div>
      <div><strong class="font-medium text-gray-600 dark:text-gray-400">OS:</strong> <span class="text-gray-800 dark:text-gray-200">{{ info.os }}</span></div>
      <div><strong class="font-medium text-gray-600 dark:text-gray-400">Firmware Version:</strong> <span class="text-gray-800 dark:text-gray-200">{{ info.firmware_version }}</span></div>
      <div><strong class="font-medium text-gray-600 dark:text-gray-400">Kernel:</strong> <span class="text-gray-800 dark:text-gray-200">{{ info.kernel }}</span></div>
      <div><strong class="font-medium text-gray-600 dark:text-gray-400">Uptime:</strong> <span class="text-gray-800 dark:text-gray-200">{{ info.uptime }}</span></div>
      <div><strong class="font-medium text-gray-600 dark:text-gray-400">CPU:</strong> <span class="text-gray-800 dark:text-gray-200">{{ info.cpu }} ({{ info.cpu_cores }} cores)</span></div>
      <div><strong class="font-medium text-gray-600 dark:text-gray-400">CPU Usage:</strong> <span class="text-gray-800 dark:text-gray-200">{{ info.cpu_usage.toFixed(2) }}%</span></div>
      <div><strong class="font-medium text-gray-600 dark:text-gray-400">Memory Usage:</strong> <span class="text-gray-800 dark:text-gray-200">{{ info.memory_used_mb }} MB / {{ info.memory_total_mb }} MB ({{ info.memory_usage.toFixed(2) }}%)</span></div>
    </div>
    <div v-else>
      <p class="text-gray-500 dark:text-gray-400">Loading system information...</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { SystemInfo } from '../../stores/system'

defineProps<{
  info: SystemInfo | null
}>()
</script>
