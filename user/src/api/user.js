import api from './client'

export const userRegister = (data) => api.post('/user/auth/register', data)
export const userLogin = (data) => api.post('/user/auth/login', data)
export const userLogout = () => api.post('/user/auth/logout')

export const fetchTasks = (params) => api.get('/user/tasks', { params })
export const fetchTaskDetail = (id) => api.get(`/user/tasks/${id}`)
export const claimTask = (id) => api.post(`/user/tasks/${id}/claim`)

export const fetchCurrentClaim = () => api.get('/user/claims/current')
export const fetchClaimHistory = () => api.get('/user/claims/history')
export const submitClaim = (id, data) => api.post(`/user/claims/${id}/submit`, data)

export const fetchProfile = () => api.get('/user/profile')
export const changePassword = (data) => api.put('/user/profile/password', data)
export const updateWechatQr = (data) => api.put('/user/profile/wechat-qr', data)
export const fetchPointsLogs = (params) => api.get('/user/points/logs', { params })

export const applyWithdrawal = (data) => api.post('/user/withdrawals', data)
export const fetchWithdrawals = () => api.get('/user/withdrawals')

export const uploadImage = (formData) =>
  api.post('/user/upload/image', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
