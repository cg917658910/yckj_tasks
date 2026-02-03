<template>
  <aside class="sidebar">
    <div class="brand">
      <div class="logo">任</div>
      <div class="brand-text">任务管理系统</div>
    </div>

    <div class="profile">
      <div class="avatar">{{ avatarText }}</div>
      <div>
        <div class="name">{{ profile.username || '用户' }}</div>
        <div class="role">普通用户</div>
      </div>
    </div>

    <div class="score-card">
      <div class="score-title">我的积分</div>
      <div class="score-value">{{ profile.available_points || 0 }}</div>
      <div class="score-sub">可提现约 {{ withdrawable }} 元</div>
      <div class="score-actions">
        <button class="link" @click="goPoints">积分历史</button>
        <button class="link" @click="goWithdraw">收款码</button>
      </div>
    </div>

    <nav class="menu">
      <RouterLink class="menu-item" to="/tasks">
        <span class="menu-icon icon-tasks"></span>
        任务列表
      </RouterLink>
      <RouterLink class="menu-item" to="/my-tasks">
        <span class="menu-icon icon-mytasks"></span>
        我的任务
      </RouterLink>
      <RouterLink class="menu-item" to="/profile">
        <span class="menu-icon icon-profile"></span>
        个人中心
      </RouterLink>
      <RouterLink class="menu-item" to="/withdraw">
        <span class="menu-icon icon-withdraw"></span>
        提现中心
      </RouterLink>
      <RouterLink class="menu-item" to="/points-history">
        <span class="menu-icon icon-points"></span>
        积分历史
      </RouterLink>
    </nav>

    <div class="sidebar-footer">
      <button class="ghost-btn" @click="logout">退出登录</button>
    </div>
  </aside>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { userLogout } from '../api/user'
import { profileState, loadProfile } from '../store/profile'

const router = useRouter()
const profile = profileState

const avatarText = computed(() => {
  const name = profile.value.username || 'U'
  return name.slice(0, 1).toUpperCase()
})

const withdrawable = computed(() => {
  const points = profile.value.available_points || 0
  return (points / 10).toFixed(2)
})

const load = async () => {
  await loadProfile()
}

const logout = async () => {
  try {
    await userLogout()
  } catch (err) {
    // ignore
  }
  localStorage.removeItem('user_token')
  router.replace('/login')
}

const goWithdraw = () => {
  router.push('/withdraw')
}

const goPoints = () => {
  router.push('/points-history')
}

onMounted(load)
</script>
