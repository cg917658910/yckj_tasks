<template>
  <div class="login">
    <div class="login-card">
      <div class="login-brand">
        <div class="logo">T</div>
        <div>
          <div class="brand-title">任务管理系统</div>
          <div class="brand-sub">管理员登录</div>
        </div>
      </div>

      <div class="login-form">
        <label>
          账号
          <input v-model="form.username" placeholder="请输入用户名" />
        </label>
        <label>
          密码
          <input v-model="form.password" type="password" placeholder="请输入密码" />
        </label>
      </div>

      <div class="login-actions">
        <button class="primary-btn" @click="submit" :disabled="loading">
          {{ loading ? '登录中...' : '登录' }}
        </button>
        <div class="login-hint">默认账号需在数据库中创建</div>
      </div>

      <div v-if="error" class="login-error">{{ error }}</div>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { adminLogin } from '../api/admin'

const router = useRouter()
const loading = ref(false)
const error = ref('')
const form = reactive({
  username: '',
  password: '',
})

const submit = async () => {
  error.value = ''
  if (!form.username || !form.password) {
    error.value = '请输入账号和密码'
    return
  }

  try {
    loading.value = true
    const res = await adminLogin(form)
    // code != 0 代表登录失败
    if (res.code !== 0) {
      throw new Error(res.message || '登录失败')
    }
    const token = res.data?.token
    if (!token) {
      throw new Error('未获取到登录凭证')
    }
    localStorage.setItem('admin_token', token)
    router.replace('/dashboard')
  } catch (err) {
    error.value = err.message || '登录失败'
  } finally {
    loading.value = false
  }
}
</script>
