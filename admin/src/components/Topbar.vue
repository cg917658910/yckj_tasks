<template>
  <header class="topbar">
    <div class="search">
      <input v-model="keyword" placeholder="搜索任务/用户/提现" />
    </div>
    <div class="top-actions">
      <!-- <button class="ghost-btn">导出数据</button> -->
      <button class="primary-btn" @click="createTask">发布任务</button>
      <button class="ghost-btn" @click="logout">退出登录</button>
      <div class="top-avatar">A</div>
    </div>
  </header>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { adminLogout } from '../api/admin'

const keyword = ref('')
const router = useRouter()

const logout = async () => {
  try {
    await adminLogout()
  } catch (err) {
    // ignore logout errors
  }
  localStorage.removeItem('admin_token')
  router.replace('/login')
}
// go to task creation page
const createTask = () => {
  router.push('/tasks?create=true')
}
</script>
