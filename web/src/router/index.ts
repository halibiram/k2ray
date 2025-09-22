import { createRouter, createWebHistory } from 'vue-router'

import LoginView from '../views/Login.vue'
import DashboardView from '../views/Dashboard.vue'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: LoginView,
  },
  {
    path: '/',
    name: 'Dashboard',
    component: DashboardView,
    meta: { requiresAuth: true },
  },
  // Redirect to dashboard if root is accessed, will be handled by nav guard later
  {
    path: '/:pathMatch(.*)*',
    redirect: '/',
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

// Add navigation guard later
// router.beforeEach((to, from, next) => { ... })

export default router
