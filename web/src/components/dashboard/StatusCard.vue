<template>
  <div class="p-4 rounded-lg shadow-md" :class="cardColor">
    <h3 class="text-lg font-semibold text-white">{{ title }}</h3>
    <p class="mt-2 text-3xl font-bold text-white">{{ statusText }}</p>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  title: string
  status: 'online' | 'offline' | 'loading' | 'running' | 'stopped' | null | undefined
}>()

const cardColor = computed(() => {
  switch (props.status) {
    case 'online':
    case 'running':
      return 'bg-green-500'
    case 'offline':
    case 'stopped':
      return 'bg-red-500'
    default:
      return 'bg-gray-400'
  }
})

const statusText = computed(() => {
  if (!props.status) return 'Unknown'
  return props.status.charAt(0).toUpperCase() + props.status.slice(1)
})
</script>
