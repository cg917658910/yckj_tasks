# API 设计（草案）

> 说明：接口分为管理端（/admin）与用户端（/user）。认证建议使用 JWT 或 Session（ThinkPHP8 可选中间件实现）。

## 认证
- POST `/admin/auth/login`
- POST `/admin/auth/logout`
- POST `/user/auth/register`
- POST `/user/auth/login`
- POST `/user/auth/logout`

## 管理端（任务发布平台）
### 积分管理
- GET `/admin/points/rules` 获取当前规则
- PUT `/admin/points/rules` 更新规则（记录变更日志）

### 任务管理
- GET `/admin/tasks` 任务列表（含状态过滤）
- POST `/admin/tasks` 发布任务
- PUT `/admin/tasks/{id}` 编辑任务
- PUT `/admin/tasks/{id}/off` 下架任务

### 任务审核
- GET `/admin/claims` 任务领取/提交列表
- POST `/admin/claims/{id}/approve` 审核通过（可自定义奖励积分）
- POST `/admin/claims/{id}/reject` 审核驳回（需填写原因）

### 提现管理
- GET `/admin/withdrawals` 提现申请列表
- POST `/admin/withdrawals/{id}/pay` 审核通过（已打款）
- POST `/admin/withdrawals/{id}/reject` 审核驳回

### 用户管理
- GET `/admin/users` 用户列表
- PUT `/admin/users/{id}/status` 禁用/启用
- POST `/admin/users/{id}/points` 手动调整积分
- GET `/admin/users/{id}/tasks` 用户任务记录
- GET `/admin/users/{id}/withdrawals` 用户提现记录

## 用户端（任务领取平台）
### 任务大厅
- GET `/user/tasks` 可领取任务列表
- GET `/user/tasks/{id}` 任务详情
- POST `/user/tasks/{id}/claim` 领取任务（限制：用户同时仅可领取 1 个）

### 我的任务
- GET `/user/claims/current` 当前任务
- POST `/user/claims/{id}/submit` 提交成果（截图+备注）
- GET `/user/claims/history` 历史任务

### 个人中心
- GET `/user/profile` 个人信息与积分
- PUT `/user/profile/password` 修改密码
- PUT `/user/profile/wechat-qr` 上传微信收款二维码

### 提现
- POST `/user/withdrawals` 申请提现（校验最低提现金额、积分充足、已绑定二维码）
- GET `/user/withdrawals` 提现记录
