<template>
  <header class="flex items-center justify-between px-6 py-4 bg-white border-b-4 border-indigo-600">
    <div class="flex items-center">
      <h2 class="text-xl font-semibold text-gray-800">{{ $t('header.title') }}</h2>
    </div>
    <div class="flex items-center space-x-4">
      <button @click="themeStore.toggleTheme" class="px-3 py-1 text-sm font-medium rounded-md" :class="themeStore.theme === 'dark' ? 'text-white bg-gray-700' : 'text-gray-600 bg-gray-200'">
        {{ themeStore.theme === 'dark' ? 'Light' : 'Dark' }}
      </button>
      <select v-model="locale" class="block w-full px-2 py-1 text-sm border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500">
        <option value="en">English</option>
        <option value="tr">Türkçe</option>
      </select>
      <button @click="logout" class="text-sm font-medium text-gray-600 hover:text-indigo-500 focus:outline-none">
        {{ $t('common.logout') }}
      </button>
    </div>
  </header>
</template>

<script setup lang="ts">
import { useAuthStore } from '../../stores/auth'
import { useThemeStore } from '../../stores/theme'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'

const authStore = useAuthStore()
const themeStore = useThemeStore()
const router = useRouter()
const { locale } = useI18n()

const logout = () => {
  authStore.logout()
  router.push('/login')
}
</script>
