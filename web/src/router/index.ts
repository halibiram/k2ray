import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

import LoginView from '../views/Login.vue'
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
    path: '/',
    component: Layout,
    meta: { requiresAuth: true },
    children: [
      { path: '', name: 'Dashboard', component: DashboardView },
      { path: '/configurations', name: 'Configurations', component: ConfigManagerView },
      { path: '/system-status', name: 'SystemStatus', component: SystemStatusView },
      { path: '/monitoring', name: 'Monitoring', component: MonitoringView },
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
router.beforeEach((to, from, next) => {
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
