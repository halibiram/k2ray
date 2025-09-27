<template>
  <div class="p-4 sm:p-6 lg:p-8">
    <h1 class="text-3xl font-bold text-base-content">Dashboard</h1>
    <p class="mt-2 text-base-content/70">System status overview</p>

    <div class="mt-8 grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
      <StatusCard title="K2Ray API Status" :status="systemStore.apiStatus" />
      <StatusCard title="V2Ray Service Status" :status="systemStore.v2rayStatus?.status" />
    </div>

    <div v-if="systemStore.error" class="mt-4 alert alert-error">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        class="stroke-current shrink-0 h-6 w-6"
        fill="none"
        viewBox="0 0 24 24"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
        />
      </svg>
      <span>An error occurred: {{ systemStore.error }}</span>
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