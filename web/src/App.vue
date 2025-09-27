<template>
  <div class="drawer">
    <div class="drawer-content flex flex-col min-h-screen">
      <!-- Navbar -->
      <header class="w-full navbar bg-base-300">
        <input id="my-drawer-3" type="checkbox" class="drawer-toggle" />
        <div class="flex-none lg:hidden">
          <label for="my-drawer-3" aria-label="open sidebar" class="btn btn-square btn-ghost">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              class="inline-block w-6 h-6 stroke-current"
            >
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"></path>
            </svg>
          </label>
        </div>
        <div class="flex-1 px-2 mx-2">K2Ray</div>
        <nav class="flex-none hidden lg:block">
          <ul class="menu menu-horizontal">
            <!-- Navbar menu content here -->
            <li><a>Navbar Item 1</a></li>
            <li><a>Navbar Item 2</a></li>
          </ul>
        </nav>
        <ThemeToggler />
      </header>
      <!-- Page content here -->
      <main class="flex-1">
        <router-view />
      </main>
    </div>
    <aside class="drawer-side">
      <label for="my-drawer-3" aria-label="close sidebar" class="drawer-overlay"></label>
      <ul class="menu p-4 w-80 min-h-full bg-base-200">
        <!-- Sidebar content here -->
        <li><a>Sidebar Item 1</a></li>
        <li><a>Sidebar Item 2</a></li>
      </ul>
    </aside>
  </div>
</template>

<script setup lang="ts">
import { useThemeStore } from './stores/theme'
import ThemeToggler from './components/ThemeToggler.vue'

// Initialize the theme store to apply the theme on app load
useThemeStore()
</script>