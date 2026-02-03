# 初始化与种子数据

## 1. 导入数据库结构
- 使用 `database/schema.sql`

## 2. 默认规则
- 使用 `database/seed.sql` 会插入默认积分规则

## 3. 创建管理员账号
- 生成密码哈希：
  - `php -r "echo password_hash('admin123', PASSWORD_BCRYPT), PHP_EOL;"`
- 执行 SQL：
  - `INSERT INTO admin_users (username, password_hash, status) VALUES ('admin', '<HASH>', 1);`
