<template>
  <div class="card shadow-xl" :class="cardColor">
    <div class="card-body">
      <h2 class="card-title">{{ title }}</h2>
      <p class="text-3xl font-bold">{{ statusText }}</p>
    </div>
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
      return 'bg-success text-success-content'
    case 'offline':
    case 'stopped':
      return 'bg-error text-error-content'
    default:
      return 'bg-base-200 text-base-content'
  }
})

const statusText = computed(() => {
  if (!props.status) return 'Unknown'
  return props.status.charAt(0).toUpperCase() + props.status.slice(1)
})
</script>