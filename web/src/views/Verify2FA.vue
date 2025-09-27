<template>
  <div class="flex items-center justify-center min-h-screen bg-gray-100 dark:bg-gray-900">
    <div class="w-full max-w-md p-8 space-y-6 bg-white rounded-lg shadow-md dark:bg-gray-800">
      <h2 class="text-2xl font-bold text-center text-gray-900 dark:text-white">
        Two-Factor Authentication
      </h2>
      <p class="text-center text-gray-600 dark:text-gray-400">
        Enter the code from your authenticator app.
      </p>
      <form @submit.prevent="verifyCode">
        <div class="space-y-4">
          <div>
            <label for="2fa-code" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
              6-Digit Code
            </label>
            <input
              id="2fa-code"
              v-model="code"
              type="text"
              required
              class="w-full px-3 py-2 mt-1 border border-gray-300 rounded-md shadow-sm dark:bg-gray-700 dark:border-gray-600 dark:text-white focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              placeholder="123456"
              autocomplete="one-time-code"
            />
          </div>
          <button
            type="submit"
            class="w-full px-4 py-2 font-medium text-white bg-indigo-600 rounded-md hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
            :disabled="isLoading"
          >
            <span v-if="isLoading">Verifying...</span>
            <span v-else>Verify</span>
          </button>
        </div>
      </form>
      <p v-if="error" class="mt-4 text-sm text-center text-red-600">
        {{ error }}
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useRouter } from 'vue-router'

const code = ref('')
const isLoading = ref(false)
const error = ref<string | null>(null)

const authStore = useAuthStore()
const router = useRouter()

async function verifyCode() {
  if (!authStore.twoFactorToken) {
    error.value = '2FA process was not initiated correctly. Please log in again.'
    return
  }

  isLoading.value = true
  error.value = null

  try {
    await authStore.verify2FA(code.value)
    router.push({ name: 'Dashboard' })
  } catch (err) {
    error.value = (err as Error).message || 'An unknown error occurred.'
  } finally {
    isLoading.value = false
  }
}
</script>