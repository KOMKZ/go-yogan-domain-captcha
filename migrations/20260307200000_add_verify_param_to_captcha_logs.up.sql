SET @col_exists := (
  SELECT COUNT(*)
  FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE()
    AND TABLE_NAME = 'captcha_logs'
    AND COLUMN_NAME = 'verify_param'
);

SET @ddl := IF(
  @col_exists = 0,
  'ALTER TABLE `captcha_logs` ADD COLUMN `verify_param` TEXT COMMENT ''验证参数原始JSON'' AFTER `request_id`',
  'SELECT 1'
);

PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
