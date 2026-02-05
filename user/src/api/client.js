import axios from 'axios'
import { notify } from '../store/notify'

// 根据环境变量获取 API 基础地址
const getBaseURL = () => {
  // 优先使用环境变量配置
  if (import.meta.env.VITE_API_BASE_URL) {
    return import.meta.env.VITE_API_BASE_URL
  }
  
  // 根据开发/生产环境返回默认值
  return import.meta.env.MODE === 'production' 
    ? 'https://api.yourdomain.com'  // 生产环境默认地址
    : 'http://127.0.0.1:8000'        // 开发环境默认地址
}

const api = axios.create({
  baseURL: getBaseURL(),
  timeout: 10000,
})

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('user_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

api.interceptors.response.use(
  (response) => response.data,
  (error) => {
    const status = error?.response?.status
    if (status === 401) {
      localStorage.removeItem('user_token')
      if (window.location.pathname !== '/login') {
        window.location.href = '/login'
      }
    }
    const message = error?.response?.data?.message || error.message || '请求失败'
    notify(message, 'error')
    return Promise.reject(new Error(message))
  }
)

export default api
