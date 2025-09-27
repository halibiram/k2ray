import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue'),
  },
  {
    path: '/verify-2fa',
    name: 'Verify2FA',
    component: () => import('../views/Verify2FA.vue'),
  },
  {
    path: '/',
    component: () => import('../components/common/Layout.vue'),
    meta: { requiresAuth: true },
    children: [
      { path: '', name: 'Dashboard', component: () => import('../views/Dashboard.vue') },
      {
        path: '/configurations',
        name: 'Configurations',
        component: () => import('../views/ConfigManager.vue'),
      },
      {
        path: '/system-status',
        name: 'SystemStatus',
        component: () => import('../views/SystemStatus.vue'),
      },
      { path: '/monitoring', name: 'Monitoring', component: () => import('../views/Monitoring.vue') },
      { path: '/settings/2fa', name: 'Setup2FA', component: () => import('../views/Setup2FA.vue') },
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