ALTER TABLE `captcha_logs` ADD COLUMN `verify_param` TEXT COMMENT '验证参数原始JSON' AFTER `request_id`;
