import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

import LoginView from '../views/Login.vue'
import Verify2FAView from '../views/Verify2FA.vue'
import Setup2FAView from '../views/Setup2FA.vue'
import DashboardView from '../views/Dashboard.vue'
import ConfigManagerView from '../views/ConfigManager.vue'
import SystemStatusView from '../views/SystemStatus.vue'
import MonitoringView from '../views/Monitoring.vue'
import Layout from '../components/common/Layout.vue'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: LoginView,
  },
  {
    path: '/verify-2fa',
    name: 'Verify2FA',
    component: Verify2FAView,
  },
  {
    path: '/',
    component: Layout,
    meta: { requiresAuth: true },
    children: [
      { path: '', name: 'Dashboard', component: DashboardView },
      { path: '/configurations', name: 'Configurations', component: ConfigManagerView },
      { path: '/system-status', name: 'SystemStatus', component: SystemStatusView },
      { path: '/monitoring', name: 'Monitoring', component: MonitoringView },
      { path: '/settings/2fa', name: 'Setup2FA', component: Setup2FAView },
    ],
  },
  // Catch-all to redirect to the main page
  {
    path: '/:pathMatch(.*)*',
    redirect: '/',
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

// Navigation Guard
router.beforeEach((to, _from, next) => {
  // We need to initialize the store here to use it outside of a component setup
  const authStore = useAuthStore()

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    // If route requires auth and user is not authenticated, redirect to login
    next({ name: 'Login' })
  } else {
    // Otherwise, allow navigation
    next()
  }
})

export default router
