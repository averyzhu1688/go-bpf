# ************************************************************
# Sequel Ace SQL dump
# 版本号： 20077
#
# https://sequel-ace.com/
# https://github.com/Sequel-Ace/Sequel-Ace
#
# 主机: 10.0.0.106 (MySQL 8.4.0)
# 数据库: go-gin
# 生成时间: 2025-05-13 03:17:31 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
SET NAMES utf8mb4;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE='NO_AUTO_VALUE_ON_ZERO', SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# 转储表 t_sys_roles
# ------------------------------------------------------------

DROP TABLE IF EXISTS `t_sys_roles`;

CREATE TABLE `t_sys_roles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(50) NOT NULL,
  `code` varchar(50) NOT NULL,
  `description` varchar(200) DEFAULT NULL,
  `permissions` json DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_t_sys_roles_name` (`name`),
  UNIQUE KEY `idx_t_sys_roles_code` (`code`),
  KEY `idx_t_sys_roles_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

LOCK TABLES `t_sys_roles` WRITE;
/*!40000 ALTER TABLE `t_sys_roles` DISABLE KEYS */;

INSERT INTO `t_sys_roles` (`id`, `created_at`, `updated_at`, `deleted_at`, `name`, `code`, `description`, `permissions`)
VALUES
	(1,'2025-05-12 20:50:24.985','2025-05-12 20:50:24.985',NULL,'超级管理员','superuser','系统最高管理员，拥有所有权限','[\"*\"]'),
	(2,'2025-05-12 20:50:24.985','2025-05-12 20:50:24.985',NULL,'管理员','admin','系统管理员，拥有所有权限','[\"*\"]'),
	(3,'2025-05-12 20:50:24.985','2025-05-12 20:50:24.985',NULL,'普通用户','user','普通用户，拥有基本权限','[\"user:view\", \"user:edit\", \"content:view\", \"content:create\", \"content:edit\"]'),
	(4,'2025-05-12 20:50:24.985','2025-05-12 20:50:24.985',NULL,'访客','guest','访客，仅拥有查看权限','[\"user:view\", \"content:view\"]');

/*!40000 ALTER TABLE `t_sys_roles` ENABLE KEYS */;
UNLOCK TABLES;


# 转储表 t_sys_users
# ------------------------------------------------------------

DROP TABLE IF EXISTS `t_sys_users`;

CREATE TABLE `t_sys_users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `username` varchar(50) NOT NULL,
  `password` varchar(100) NOT NULL,
  `email` varchar(100) NOT NULL,
  `phone` varchar(20) DEFAULT NULL,
  `nickname` varchar(50) DEFAULT NULL,
  `role_id` bigint unsigned DEFAULT '3',
  `status` bigint DEFAULT '1',
  `last_login` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_t_sys_users_username` (`username`),
  UNIQUE KEY `idx_t_sys_users_email` (`email`),
  KEY `idx_t_sys_users_deleted_at` (`deleted_at`),
  KEY `fk_t_sys_users_role` (`role_id`),
  CONSTRAINT `fk_t_sys_users_role` FOREIGN KEY (`role_id`) REFERENCES `t_sys_roles` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

LOCK TABLES `t_sys_users` WRITE;
/*!40000 ALTER TABLE `t_sys_users` DISABLE KEYS */;

INSERT INTO `t_sys_users` (`id`, `created_at`, `updated_at`, `deleted_at`, `username`, `password`, `email`, `phone`, `nickname`, `role_id`, `status`, `last_login`)
VALUES (1,'2025-05-12 20:50:25.070','2025-05-12 20:50:25.070',NULL,'admin','$2a$10$Hzk5L1zvBhIlYv6ouuwe4uN4/bDQ.BGSkuSBjulzjtT5J0MiMUkn6','admin@example.com','','系统管理员',1,1,NULL)

/*!40000 ALTER TABLE `t_sys_users` ENABLE KEYS */;
UNLOCK TABLES;



/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
