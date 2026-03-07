-- tk_draw_record 开奖区结构升级脚本（MySQL 5.7+）
-- 目标：
-- 1) 新增开奖区独立主表 tk_draw_record；
-- 2) 将 tk_lottery_info 表语义修正为“图库内容与竞猜配置”。
-- 执行前请先：USE your_db;

CREATE TABLE IF NOT EXISTS tk_draw_record (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  special_lottery_id BIGINT UNSIGNED NOT NULL COMMENT '所属彩种ID（关联tk_special_lottery.id）',
  issue VARCHAR(32) NOT NULL COMMENT '开奖期号（如2026-063）',
  year INT NOT NULL COMMENT '年份（如2026）',
  draw_at DATETIME(3) NOT NULL COMMENT '开奖时间',
  normal_draw_result VARCHAR(64) NOT NULL DEFAULT '' COMMENT '普通号码（6个，逗号分隔）',
  special_draw_result VARCHAR(16) NOT NULL DEFAULT '' COMMENT '特别号码（1个）',
  draw_result VARCHAR(120) NOT NULL DEFAULT '' COMMENT '兼容字段：完整开奖号码（普通6个+特别号）',
  draw_labels VARCHAR(255) NOT NULL DEFAULT '' COMMENT '号码标签（与号码一一对应，格式示例：羊/土,蛇/金）',
  zodiac_labels VARCHAR(255) NOT NULL DEFAULT '' COMMENT '号码对应属相标签（与号码一一对应，示例：羊,蛇,马）',
  wuxing_labels VARCHAR(255) NOT NULL DEFAULT '' COMMENT '号码对应五行标签（与号码一一对应，示例：土,金,火）',
  playback_url VARCHAR(255) NOT NULL DEFAULT '' COMMENT '开奖回放地址（直播结束后录入）',
  special_single_double VARCHAR(16) NOT NULL DEFAULT '' COMMENT '特码单双（如：双）',
  special_big_small VARCHAR(16) NOT NULL DEFAULT '' COMMENT '特码大小（如：大）',
  sum_single_double VARCHAR(16) NOT NULL DEFAULT '' COMMENT '总和单双（如：双）',
  sum_big_small VARCHAR(16) NOT NULL DEFAULT '' COMMENT '总和大小（如：大）',
  recommend_six VARCHAR(120) NOT NULL DEFAULT '' COMMENT '六肖推荐（空格分隔）',
  recommend_four VARCHAR(120) NOT NULL DEFAULT '' COMMENT '四肖推荐（空格分隔）',
  recommend_one VARCHAR(32) NOT NULL DEFAULT '' COMMENT '一肖推荐',
  recommend_ten VARCHAR(255) NOT NULL DEFAULT '' COMMENT '十码推荐（空格分隔）',
  special_code VARCHAR(16) NOT NULL DEFAULT '' COMMENT '特码数字',
  normal_code VARCHAR(120) NOT NULL DEFAULT '' COMMENT '正码（逗号分隔）',
  zheng1 VARCHAR(120) NOT NULL DEFAULT '' COMMENT '正1特描述（如：大双,合双,尾大,蓝波）',
  zheng2 VARCHAR(120) NOT NULL DEFAULT '' COMMENT '正2特描述',
  zheng3 VARCHAR(120) NOT NULL DEFAULT '' COMMENT '正3特描述',
  zheng4 VARCHAR(120) NOT NULL DEFAULT '' COMMENT '正4特描述',
  zheng5 VARCHAR(120) NOT NULL DEFAULT '' COMMENT '正5特描述',
  zheng6 VARCHAR(120) NOT NULL DEFAULT '' COMMENT '正6特描述',
  is_current TINYINT NOT NULL DEFAULT 0 COMMENT '是否当前期：1是；0否',
  status TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1启用；0停用',
  sort INT NOT NULL DEFAULT 0 COMMENT '排序值，越小越靠前',
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (id),
  UNIQUE KEY uk_tk_draw_record_issue (special_lottery_id, issue),
  KEY idx_tk_draw_record_special (special_lottery_id),
  KEY idx_tk_draw_record_year (year),
  KEY idx_tk_draw_record_draw_at (draw_at),
  KEY idx_tk_draw_record_status_sort (status, sort)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='开奖区开奖记录表（首页开奖区/历史开奖/开奖详情）';

ALTER TABLE tk_lottery_info
  COMMENT = '图库图纸内容与竞猜配置表（不承载开奖区历史主数据）';
