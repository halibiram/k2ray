<template>
  <div>
    <h1 class="text-3xl font-bold text-gray-800">Dashboard</h1>
    <p class="mt-2 text-gray-600">System status overview</p>

    <div class="mt-8 grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
      <StatusCard title="K2Ray API Status" :status="systemStore.apiStatus" />
      <StatusCard title="V2Ray Service Status" :status="systemStore.v2rayStatus?.status" />
    </div>

    <div v-if="systemStore.error" class="mt-4 p-4 text-red-700 bg-red-100 border border-red-400 rounded-md">
      <p>An error occurred: {{ systemStore.error }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useSystemStore } from '../stores/system'
import { useWebSocket } from '../services/useWebSocket'
import StatusCard from '../components/dashboard/StatusCard.vue'

const systemStore = useSystemStore()

// Initialize WebSocket connection
useWebSocket()

onMounted(() => {
  systemStore.fetchApiStatus()
  systemStore.fetchV2rayStatus()
})
</script>
