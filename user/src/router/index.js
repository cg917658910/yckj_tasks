import { createRouter, createWebHistory } from 'vue-router'

import Login from '../pages/Login.vue'
import MyTasks from '../pages/MyTasks.vue'
import PointsHistory from '../pages/PointsHistory.vue'
import Profile from '../pages/Profile.vue'
import TaskDetail from '../pages/TaskDetail.vue'
import TaskList from '../pages/TaskList.vue'
import Withdraw from '../pages/Withdraw.vue'

const routes = [
  { path: '/', redirect: '/tasks' },
  { path: '/login', name: 'login', component: Login, meta: { public: true, hideLayout: true } },
  { path: '/tasks', name: 'tasks', component: TaskList },
  { path: '/tasks/:id', name: 'task-detail', component: TaskDetail },
  { path: '/my-tasks', name: 'my-tasks', component: MyTasks },
  { path: '/profile', name: 'profile', component: Profile },
  { path: '/withdraw', name: 'withdraw', component: Withdraw },
  { path: '/points-history', name: 'points-history', component: PointsHistory },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

router.beforeEach((to) => {
  if (to.meta.public) return true
  const token = localStorage.getItem('user_token')
  if (token && to.path === '/login') {
    return { path: '/tasks' }
  }
  if (!token && !to.meta.public) {
    return { path: '/login' }
  }
  return true
})

export default router
