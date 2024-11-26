
-- +migrate Up
CREATE TABLE `users` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) DEFAULT '',
  `email` varchar(100) DEFAULT '',
  `password` varchar(100) DEFAULT '',
  `updated_at` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 插入用户数据,密码是"password123"的哈希加密形式
INSERT INTO users (id, name, email, password, created_at, updated_at)
VALUES (1, 'Test User', 'testuser@example.com', '$2b$12$OxPgUXTZgsawpNQpYZjIyOjkMHqNBIn0dxonVGEDvoQ3/SC8kVL7q
', NOW(), NOW());
-- 插入管理员数据,密码是"adminPSW"的哈希加密形式
INSERT INTO users (id, name, email, password, created_at, updated_at)
VALUES (2, 'admin', 'admin', '$2b$12$UzuwIYad47XQToQ8CRaRGeF8YVeipYGaEcxaSoKFq4s55l4Wa2RfK', NOW(), NOW());
-- +migrate Down
DELETE FROM users WHERE id = 2;
DELETE FROM users WHERE id = 1;
DROP TABLE `users`;
