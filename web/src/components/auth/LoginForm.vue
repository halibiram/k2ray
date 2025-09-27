<template>
  <form class="space-y-6" @submit.prevent="handleSubmit">
    <div>
      <label for="username" class="block text-sm font-medium text-gray-700">Username</label>
      <div class="mt-1">
        <input
          id="username"
          v-model="username"
          name="username"
          type="text"
          required
          class="block w-full px-3 py-2 placeholder-gray-400 border border-gray-300 rounded-md shadow-sm appearance-none focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
        />
      </div>
    </div>

    <div>
      <label for="password" class="block text-sm font-medium text-gray-700">Password</label>
      <div class="mt-1">
        <input
          id="password"
          v-model="password"
          name="password"
          type="password"
          required
          class="block w-full px-3 py-2 placeholder-gray-400 border border-gray-300 rounded-md shadow-sm appearance-none focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
        />
      </div>
    </div>

    <div>
      <button
        type="submit"
        class="flex justify-center w-full px-4 py-2 text-sm font-medium text-white bg-indigo-600 border border-transparent rounded-md shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
      >
        Sign in
      </button>
    </div>

    <div v-if="authStore.error" class="p-4 text-sm text-red-700 bg-red-100 rounded-md" role="alert">
      <span class="font-medium">Error:</span> {{ authStore.error }}
    </div>
  </form>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../../stores/auth'

const username = ref('')
const password = ref('')
const router = useRouter()
const authStore = useAuthStore()

const handleSubmit = async () => {
  const success = await authStore.login({
    username: username.value,
    password: password.value,
  })

  if (success) {
    router.push('/') // Redirect to Dashboard
  }
}
</script>
