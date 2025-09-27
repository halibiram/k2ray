<template>
  <div class="bg-white shadow-md rounded-lg p-6">
    <h2 class="text-xl font-semibold text-gray-800 border-b pb-2">System Information</h2>
    <div v-if="info" class="mt-4 grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
      <div><strong class="font-medium text-gray-600">Hostname:</strong> <span class="text-gray-800">{{ info.hostname }}</span></div>
      <div><strong class="font-medium text-gray-600">Keenetic Model:</strong> <span class="text-gray-800">{{ info.keenetic_model }}</span></div>
      <div><strong class="font-medium text-gray-600">OS:</strong> <span class="text-gray-800">{{ info.os }}</span></div>
      <div><strong class="font-medium text-gray-600">Firmware Version:</strong> <span class="text-gray-800">{{ info.firmware_version }}</span></div>
      <div><strong class="font-medium text-gray-600">Kernel:</strong> <span class="text-gray-800">{{ info.kernel }}</span></div>
      <div><strong class="font-medium text-gray-600">Uptime:</strong> <span class="text-gray-800">{{ info.uptime }}</span></div>
      <div><strong class="font-medium text-gray-600">CPU:</strong> <span class="text-gray-800">{{ info.cpu }} ({{ info.cpu_cores }} cores)</span></div>
      <div><strong class="font-medium text-gray-600">CPU Usage:</strong> <span class="text-gray-800">{{ info.cpu_usage.toFixed(2) }}%</span></div>
      <div><strong class="font-medium text-gray-600">Memory Usage:</strong> <span class="text-gray-800">{{ info.memory_used_mb }} MB / {{ info.memory_total_mb }} MB ({{ info.memory_usage.toFixed(2) }}%)</span></div>
    </div>
    <div v-else>
      <p class="text-gray-500">Loading system information...</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { SystemInfo } from '../../stores/system'

defineProps<{
  info: SystemInfo | null
}>()
</script>
