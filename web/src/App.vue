<template>
  <router-view />
</template>

<script setup lang="ts">
import { useThemeStore } from './stores/theme'

// Initialize the theme store to apply the theme on app load
useThemeStore()
</script>
