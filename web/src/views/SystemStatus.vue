<template>
  <div>
    <h1 class="text-3xl font-bold text-gray-800 dark:text-gray-100">System Status</h1>
    <p class="mt-2 text-gray-600 dark:text-gray-300">Detailed system information and logs.</p>

    <div class="mt-8">
      <SystemInfo :info="systemStore.systemInfo" />
    </div>

    <div class="mt-8">
      <LogViewer :logs="systemStore.systemLogs" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useSystemStore } from '../stores/system'
import SystemInfo from '../components/system/SystemInfo.vue'
import LogViewer from '../components/system/LogViewer.vue'

const systemStore = useSystemStore()

onMounted(() => {
  systemStore.fetchSystemInfo()
  systemStore.fetchSystemLogs()
})
</script>
