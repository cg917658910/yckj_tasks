# Go Gin Backend

并行的 Go 版后端，接口路径与 PHP 版一致（`/admin` / `/user`）。

## 配置
通过环境变量：
- `APP_ADDR` 监听地址（默认 :8080）
- `DB_DSN` MySQL DSN
- `JWT_SECRET` JWT 密钥
- `JWT_ISSUER` JWT 发行者
- `JWT_TTL` JWT 过期秒数
- `UPLOAD_DIR` 上传目录（默认 ./public/uploads）
- `BASE_URL` 访问域名（用于生成上传访问地址）
- `ALLOW_ORIGIN` CORS 允许来源

## 运行
```bash
cd go-backend
# go mod download
# go run ./cmd/server
```

## 已实现接口
- 管理端：任务、审核、积分规则/历史、提现、用户、上传
- 用户端：任务大厅、我的任务、个人中心、提现、积分历史、上传
- 认证：管理员/用户登录注册

> Go 版与 PHP 版共用 MySQL 表结构。
