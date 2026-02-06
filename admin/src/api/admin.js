import api from './client'

export const adminLogin = (data) => api.post('/admin/auth/login', data)
export const adminLogout = () => api.post('/admin/auth/logout')

export const fetchDashboard = () => api.get('/admin/tasks')

export const fetchTasks = (params) => api.get('/admin/tasks', { params })
export const createTask = (data) => api.post('/admin/tasks', data)
export const updateTask = (id, data) => api.put(`/admin/tasks/${id}`, data)
export const offTask = (id) => api.put(`/admin/tasks/${id}/off`)
export const onTask = (id) => api.put(`/admin/tasks/${id}/on`)
export const uploadImage = (formData) =>
  api.post('/admin/upload/image', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })

export const fetchClaims = (params) => api.get('/admin/claims', { params })
export const approveClaim = (id, data) => api.post(`/admin/claims/${id}/approve`, data)
export const rejectClaim = (id, data) => api.post(`/admin/claims/${id}/reject`, data)

export const fetchPointsRule = () => api.get('/admin/points/rules')
export const updatePointsRule = (data) => api.put('/admin/points/rules', data)
export const fetchPointsLogs = (params) => api.get('/admin/points/logs', { params })

export const fetchWithdrawals = (params) => api.get('/admin/withdrawals', { params })
export const payWithdrawal = (id) => api.post(`/admin/withdrawals/${id}/pay`)
export const rejectWithdrawal = (id, data) => api.post(`/admin/withdrawals/${id}/reject`, data)

export const fetchUsers = () => api.get('/admin/users')
export const updateUserStatus = (id, data) => api.put(`/admin/users/${id}/status`, data)
export const adjustUserPoints = (id, data) => api.post(`/admin/users/${id}/points`, data)
