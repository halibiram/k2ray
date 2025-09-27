<template>
  <div class="hidden md:flex flex-col w-64 bg-white dark:bg-gray-800 border-r border-gray-200 dark:border-gray-700">
    <div class="flex items-center justify-center h-16 bg-gray-100 dark:bg-gray-900">
      <span class="text-gray-800 dark:text-white font-bold uppercase">K2Ray</span>
    </div>
    <div class="flex flex-col flex-1 overflow-y-auto">
      <nav class="flex-1 px-2 py-4 bg-white dark:bg-gray-800">
        <router-link
          to="/"
          class="flex items-center px-4 py-2 text-gray-600 dark:text-gray-100 hover:bg-gray-200 dark:hover:bg-gray-700 rounded-md"
          active-class="bg-gray-200 dark:bg-gray-900 text-gray-900 dark:text-white"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
          </svg>
          Dashboard
        </router-link>
        <router-link
          to="/configurations"
          class="flex items-center px-4 py-2 mt-2 text-gray-500 dark:text-gray-400 hover:bg-gray-200 dark:hover:bg-gray-700 hover:text-gray-700 dark:hover:text-gray-100 rounded-md"
          active-class="bg-gray-200 dark:bg-gray-900 text-gray-900 dark:text-white"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
          </svg>
          Configurations
        </router-link>
        <router-link
          to="/system-status"
          class="flex items-center px-4 py-2 mt-2 text-gray-500 dark:text-gray-400 hover:bg-gray-200 dark:hover:bg-gray-700 hover:text-gray-700 dark:hover:text-gray-100 rounded-md"
          active-class="bg-gray-200 dark:bg-gray-900 text-gray-900 dark:text-white"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
          </svg>
          System
        </router-link>
        <router-link
          to="/monitoring"
          class="flex items-center px-4 py-2 mt-2 text-gray-500 dark:text-gray-400 hover:bg-gray-200 dark:hover:bg-gray-700 hover:text-gray-700 dark:hover:text-gray-100 rounded-md"
          active-class="bg-gray-200 dark:bg-gray-900 text-gray-900 dark:text-white"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 12l3-3 3 3 4-4M8 21l4-4 4 4M3 4h18" />
          </svg>
          Monitoring
        </router-link>
      </nav>
    </div>
  </div>
</template>

<script setup lang="ts">
// Placeholder component
</script>
