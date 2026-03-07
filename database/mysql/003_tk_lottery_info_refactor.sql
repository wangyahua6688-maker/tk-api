-- tk_lottery_info 结构升级脚本（兼容 MySQL 5.7+）
-- 功能：分类ID、6+1开奖号码分列、回放地址、索引与注释补齐。
-- 执行前请先：USE your_db;

-- 1) category_id
SET @ddl = (
  SELECT IF(COUNT(1) = 0,
    'ALTER TABLE tk_lottery_info ADD COLUMN category_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT ''图库分类ID（关联tk_lottery_category.id）'' AFTER special_lottery_id',
    'SELECT 1')
  FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'tk_lottery_info' AND COLUMN_NAME = 'category_id'
);
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- 2) normal_draw_result
SET @ddl = (
  SELECT IF(COUNT(1) = 0,
    'ALTER TABLE tk_lottery_info ADD COLUMN normal_draw_result VARCHAR(64) NOT NULL DEFAULT '''' COMMENT ''普通号码（6个，逗号分隔）'' AFTER draw_code',
    'SELECT 1')
  FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'tk_lottery_info' AND COLUMN_NAME = 'normal_draw_result'
);
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- 3) special_draw_result
SET @ddl = (
  SELECT IF(COUNT(1) = 0,
    'ALTER TABLE tk_lottery_info ADD COLUMN special_draw_result VARCHAR(16) NOT NULL DEFAULT '''' COMMENT ''特别号码（1个）'' AFTER normal_draw_result',
    'SELECT 1')
  FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'tk_lottery_info' AND COLUMN_NAME = 'special_draw_result'
);
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- 4) playback_url
SET @ddl = (
  SELECT IF(COUNT(1) = 0,
    'ALTER TABLE tk_lottery_info ADD COLUMN playback_url VARCHAR(255) NOT NULL DEFAULT '''' COMMENT ''直播回放地址（直播结束后录入）'' AFTER draw_at',
    'SELECT 1')
  FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'tk_lottery_info' AND COLUMN_NAME = 'playback_url'
);
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- 5) category_id 索引
SET @ddl = (
  SELECT IF(COUNT(1) = 0,
    'ALTER TABLE tk_lottery_info ADD KEY idx_tk_lottery_info_category_id (category_id)',
    'SELECT 1')
  FROM information_schema.STATISTICS
  WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'tk_lottery_info' AND INDEX_NAME = 'idx_tk_lottery_info_category_id'
);
PREPARE stmt FROM @ddl; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- 6) 基于 category_tag 回填 category_id（旧数据兼容）。
UPDATE tk_lottery_info li
LEFT JOIN tk_lottery_category lc
  ON lc.category_key COLLATE utf8mb4_general_ci = li.category_tag COLLATE utf8mb4_general_ci
  OR lc.name COLLATE utf8mb4_general_ci = li.category_tag COLLATE utf8mb4_general_ci
SET li.category_id = CASE
  WHEN li.category_id > 0 THEN li.category_id
  WHEN lc.id IS NOT NULL THEN lc.id
  ELSE 0
END;

-- 7) 基于 draw_result 回填 6+1 字段（旧数据兼容）。
UPDATE tk_lottery_info
SET
  normal_draw_result = CASE
    WHEN TRIM(IFNULL(normal_draw_result, '')) <> '' THEN normal_draw_result
    ELSE TRIM(BOTH ',' FROM SUBSTRING_INDEX(REPLACE(IFNULL(draw_result, ''), ' ', ''), ',', 6))
  END,
  special_draw_result = CASE
    WHEN TRIM(IFNULL(special_draw_result, '')) <> '' THEN special_draw_result
    ELSE TRIM(BOTH ',' FROM SUBSTRING_INDEX(REPLACE(IFNULL(draw_result, ''), ' ', ''), ',', -1))
  END
WHERE TRIM(IFNULL(draw_result, '')) <> '';

-- 8) 统一兼容字段 draw_result（普通6个+特别号）。
UPDATE tk_lottery_info
SET draw_result = CONCAT_WS(',',
  NULLIF(TRIM(IFNULL(normal_draw_result, '')), ''),
  NULLIF(TRIM(IFNULL(special_draw_result, '')), '')
)
WHERE TRIM(IFNULL(normal_draw_result, '')) <> '' OR TRIM(IFNULL(special_draw_result, '')) <> '';

-- 9) 注释更新。
ALTER TABLE tk_lottery_info
  MODIFY COLUMN category_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '图库分类ID（关联tk_lottery_category.id）',
  MODIFY COLUMN category_tag VARCHAR(32) NOT NULL DEFAULT '' COMMENT '分类标识兼容字段（通常等于category_key）',
  MODIFY COLUMN normal_draw_result VARCHAR(64) NOT NULL DEFAULT '' COMMENT '普通号码（6个，逗号分隔）',
  MODIFY COLUMN special_draw_result VARCHAR(16) NOT NULL DEFAULT '' COMMENT '特别号码（1个）',
  MODIFY COLUMN draw_result VARCHAR(120) NOT NULL DEFAULT '' COMMENT '兼容字段：完整开奖号码（普通6个+特别号）',
  MODIFY COLUMN playback_url VARCHAR(255) NOT NULL DEFAULT '' COMMENT '直播回放地址（直播结束后录入）',
  COMMENT = '开奖记录与图纸内容表（含6+1开奖号码、回放地址、详情展示指标）';
