-- Seed data for task_system
USE task_system;

-- Ensure default points rules exist
INSERT INTO points_rules (exchange_rate, min_withdraw_amount, register_bonus_points)
SELECT 10, 10.00, 10
WHERE NOT EXISTS (SELECT 1 FROM points_rules);

-- Create an admin account manually
-- 1) Generate hash: php -r "echo password_hash('admin123', PASSWORD_BCRYPT), PHP_EOL;"
-- 2) Insert with the generated hash
-- INSERT INTO admin_users (username, password_hash, status) VALUES ('admin', '<HASH>', 1);
