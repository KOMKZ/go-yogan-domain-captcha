CREATE TABLE IF NOT EXISTS `captcha_logs` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `scene` VARCHAR(100) NOT NULL DEFAULT '' COMMENT '验证场景',
    `ip` VARCHAR(50) NOT NULL DEFAULT '' COMMENT '请求IP',
    `user_agent` VARCHAR(500) NOT NULL DEFAULT '' COMMENT 'User-Agent',
    `success` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '验证是否通过',
    `request_id` VARCHAR(100) NOT NULL DEFAULT '' COMMENT '三方请求ID',
    `verify_param` TEXT COMMENT '验证参数原始JSON',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_scene` (`scene`),
    KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
