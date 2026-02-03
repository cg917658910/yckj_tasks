<template>
  <div class="login">
    <div class="login-card">
      <div class="login-brand">
        <div class="logo">任</div>
        <div>
          <div class="brand-title">任务管理系统</div>
          <div class="brand-sub">用户登录/注册</div>
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
        <button class="primary-btn" @click="submit('login')" :disabled="loading">
          {{ loading ? '处理中...' : '登录' }}
        </button>
        <button class="ghost-btn" @click="submit('register')" :disabled="loading">
          注册账号
        </button>
      </div>

      <div v-if="error" class="login-error">{{ error }}</div>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { userLogin, userRegister } from '../api/user'
import { notify } from '../store/notify'

const router = useRouter()
const loading = ref(false)
const error = ref('')
const form = reactive({
  username: '',
  password: '',
})

const submit = async (type) => {
  error.value = ''
  if (!form.username || !form.password) {
    error.value = '请输入账号和密码'
    return
  }

  try {
    loading.value = true
    const res = type === 'register' ? await userRegister(form) : await userLogin(form)
    const token = res.data?.token
    if (!token) {
      throw new Error('未获取到登录凭证')
    }
    localStorage.setItem('user_token', token)
    notify(type === 'register' ? '注册成功' : '登录成功', 'success')
    router.replace('/tasks')
  } catch (err) {
    error.value = err.message || '操作失败'
  } finally {
    loading.value = false
  }
}
</script>
