# 任务管理系统 - 后端 (ThinkPHP8 + MySQL)

## 目录结构
- `backend/database/schema.sql` 数据库建表与初始化规则
- `backend/docs/api.md` API 说明（后端接口清单）
- `backend/docs/status.md` 状态定义
- `backend/.env.example` 环境变量示例
- `backend/app` 业务代码（控制器、服务、模型、中间件）
- `backend/route/app.php` 路由配置
- `backend/config/jwt.php` JWT 配置

## 初始化方式
由于 ThinkPHP8 框架本体需要通过 Composer 初始化，推荐流程如下：

1. 在任意临时目录执行：
   - `composer create-project topthink/think=8.* task-system`
2. 将本仓库的以下目录拷贝到 `task-system` 根目录合并：
   - `app/` `route/` `config/jwt.php` `app/common.php`
3. 配置 `.env` 并导入 `database/schema.sql`

> 如果你希望我直接在当前仓库内完成 ThinkPHP8 初始化，请告诉我，我可以在你本地直接执行命令并处理合并。

## 说明
- 认证方式：JWT（Header: `Authorization: Bearer <token>`）
- 上传接口保存到 `public/uploads`，返回 URL 形如 `/uploads/xxx.jpg`
- 任务领取规则：单任务单人领取；单用户同一时间仅可有 1 个未完成任务
- 依赖库：`firebase/php-jwt`（用于 JWT 编解码）
