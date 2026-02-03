import { createRouter, createWebHistory } from 'vue-router'

import Dashboard from '../pages/Dashboard.vue'
import Tasks from '../pages/Tasks.vue'
import Claims from '../pages/Claims.vue'
import Points from '../pages/Points.vue'
import PointsHistory from '../pages/PointsHistory.vue'
import Withdrawals from '../pages/Withdrawals.vue'
import Users from '../pages/Users.vue'
import Settings from '../pages/Settings.vue'
import Login from '../pages/Login.vue'

const routes = [
  { path: '/', redirect: '/dashboard' },
  { path: '/login', name: 'login', component: Login, meta: { public: true, hideLayout: true } },
  { path: '/dashboard', name: 'dashboard', component: Dashboard, meta: { requiresAuth: true } },
  { path: '/tasks', name: 'tasks', component: Tasks, meta: { requiresAuth: true } },
  { path: '/claims', name: 'claims', component: Claims, meta: { requiresAuth: true } },
  { path: '/points', name: 'points', component: Points, meta: { requiresAuth: true } },
  { path: '/points/history', name: 'points-history', component: PointsHistory, meta: { requiresAuth: true } },
  { path: '/withdrawals', name: 'withdrawals', component: Withdrawals, meta: { requiresAuth: true } },
  { path: '/users', name: 'users', component: Users, meta: { requiresAuth: true } },
  { path: '/settings', name: 'settings', component: Settings, meta: { requiresAuth: true } },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to) => {
  if (to.meta.public) return true
  const token = localStorage.getItem('admin_token')
  if (token && to.path === '/login') {
    return { path: '/dashboard' }
  }
  if (!token && to.meta.requiresAuth) {
    return { path: '/login' }
  }
  return true
})

export default router
