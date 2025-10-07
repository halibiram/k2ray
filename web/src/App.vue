<template>
  <div class="keenetic-layout">
    <!-- Header -->
    <header class="keenetic-header">
      <div class="header-brand">
        <span class="brand-text">KEENETIC V2RAY</span>
      </div>
      <div class="header-actions">
        <button class="action-btn" @click="toggleSettings">
          <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
            <path d="M12 8c1.1 0 2-.9 2-2s-.9-2-2-2-2 .9-2 2 .9 2 2 2zm0 2c-1.1 0-2 .9-2 2s.9 2 2 2 2-.9 2-2-.9-2-2-2zm0 6c-1.1 0-2 .9-2 2s.9 2 2 2 2-.9 2-2-.9-2-2-2z"/>
          </svg>
        </button>
        <span class="status-text">{{ currentUser }}</span>
        <button class="logout-btn" @click="logout">{{ logoutText }}</button>
      </div>
    </header>

    <div class="keenetic-main">
      <!-- Sidebar -->
      <aside class="keenetic-sidebar">
        <nav class="sidebar-nav">
          <router-link 
            v-for="item in navigationItems" 
            :key="item.path"
            :to="item.path" 
            class="nav-item"
            :class="{ 'active': $route.path === item.path }"
          >
            <div class="nav-icon">
              <component :is="item.icon" />
            </div>
          </router-link>
        </nav>
      </aside>

      <!-- Main Content -->
      <main class="keenetic-content">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

// Simple icon components
const HomeIcon = () => 'H'
const ServerIcon = () => 'S'
const CogIcon = () => 'C'
const ChartBarIcon = () => 'M'
const ShieldCheckIcon = () => 'G'
const NetworkIcon = () => 'N'
const DocumentTextIcon = () => 'L'
const UserIcon = () => 'U'

// Navigation items matching Keenetic style
const navigationItems = ref([
  { path: '/', icon: HomeIcon, title: 'Ana Sayfa' },
  { path: '/v2ray-servers', icon: ServerIcon, title: 'V2Ray Sunucuları' },
  { path: '/connections', icon: NetworkIcon, title: 'Bağlantılar' },
  { path: '/monitoring', icon: ChartBarIcon, title: 'İzleme' },
  { path: '/security', icon: ShieldCheckIcon, title: 'Güvenlik' },
  { path: '/logs', icon: DocumentTextIcon, title: 'Günlükler' },
  { path: '/settings', icon: CogIcon, title: 'Ayarlar' },
  { path: '/profile', icon: UserIcon, title: 'Profil' }
])

const currentUser = ref('admin')
const logoutText = ref('Çıkış')

const toggleSettings = () => {
  router.push('/settings')
}

const logout = () => {
  // Logout logic
  console.log('Çıkış yapılıyor...')
}
</script>

<style scoped>
/* Keenetic Dark Theme Styles */
.keenetic-layout {
  @apply flex flex-col h-screen bg-slate-900;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

.keenetic-header {
  @apply flex items-center justify-between px-4 py-3;
  background: linear-gradient(135deg, #1e293b 0%, #334155 100%);
  border-bottom: 1px solid #475569;
}

.header-brand {
  @apply flex items-center;
}

.brand-text {
  @apply text-xl font-bold text-cyan-400;
  text-shadow: 0 0 10px rgba(34, 211, 238, 0.3);
}

.header-actions {
  @apply flex items-center space-x-4;
}

.action-btn {
  @apply p-2 text-slate-300 hover:text-cyan-400 transition-colors duration-200;
}

.status-text {
  @apply text-sm text-slate-300;
}

.logout-btn {
  @apply px-3 py-1 text-sm text-cyan-400 border border-cyan-400 rounded hover:bg-cyan-400 hover:text-slate-900 transition-all duration-200;
}

.keenetic-main {
  @apply flex flex-1 overflow-hidden;
}

.keenetic-sidebar {
  @apply w-16 bg-slate-800;
  background: linear-gradient(180deg, #1e293b 0%, #0f172a 100%);
  border-right: 1px solid #334155;
}

.sidebar-nav {
  @apply flex flex-col py-4;
}

.nav-item {
  @apply flex items-center justify-center h-12 mx-2 mb-2 text-slate-400 hover:text-cyan-400 hover:bg-slate-700/50 rounded-lg transition-all duration-200 relative;
}

.nav-item.active {
  @apply text-cyan-400 bg-slate-700/70;
}

.nav-item.active::after {
  content: '';
  @apply absolute left-0 top-1/2 w-1 h-6 bg-cyan-400 rounded-r;
  transform: translateY(-50%);
}

.nav-icon {
  @apply w-6 h-6 flex items-center justify-center font-semibold;
}

.keenetic-content {
  @apply flex-1 overflow-auto bg-slate-900;
  background: radial-gradient(ellipse at top, #1e293b 0%, #0f172a 70%);
}

/* Mobile Responsive */
@media (max-width: 768px) {
  .keenetic-sidebar {
    @apply w-12;
  }
  
  .nav-item {
    @apply h-10 text-xs;
  }
  
  .header-brand .brand-text {
    @apply text-lg;
  }
}
</style>