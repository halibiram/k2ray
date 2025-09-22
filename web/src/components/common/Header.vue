<template>
  <header class="flex items-center justify-between px-6 py-4 bg-white border-b-4 border-indigo-600">
    <div class="flex items-center">
      <h2 class="text-xl font-semibold text-gray-800">K2Ray</h2>
    </div>
    <div class="flex items-center">
      <button @click="logout" class="text-sm font-medium text-gray-600 hover:text-indigo-500 focus:outline-none">
        Logout
      </button>
    </div>
  </header>
</template>

<script setup lang="ts">
import { useAuthStore } from '../../stores/auth'
import { useRouter } from 'vue-router'

const authStore = useAuthStore()
const router = useRouter()

const logout = () => {
  authStore.logout()
  router.push('/login')
}
</script>
