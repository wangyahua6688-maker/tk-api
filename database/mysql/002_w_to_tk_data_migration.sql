-- w_* -> tk_* 数据迁移（外链、分类、图纸）
-- 说明：
-- 1) 推荐使用 go 工具执行（tk-business/tools/migrate_w_to_tk），该工具已内置字段/索引存在性判断；
-- 2) 本 SQL 用于人工审计迁移逻辑，不负责“判断是否存在”的分支控制。

-- ---------- 外链 ----------
UPDATE tk_external_link t
JOIN w_external_link w
  ON BINARY t.name = BINARY w.name
 AND BINARY t.url = BINARY w.url
 AND BINARY t.position = BINARY w.position
SET
  t.icon_url = IFNULL(w.icon_url, ''),
  t.group_key = IFNULL(w.group_key, ''),
  t.status = w.status,
  t.sort = w.sort,
  t.updated_at = IFNULL(w.updated_at, t.updated_at);

INSERT INTO tk_external_link (
  name, url, position, icon_url, group_key, status, sort, created_at, updated_at
)
SELECT
  w.name, w.url, w.position, IFNULL(w.icon_url, ''), IFNULL(w.group_key, ''), w.status, w.sort,
  IFNULL(w.created_at, NOW(3)), IFNULL(w.updated_at, NOW(3))
FROM w_external_link w
LEFT JOIN tk_external_link t
  ON BINARY t.name = BINARY w.name
 AND BINARY t.url = BINARY w.url
 AND BINARY t.position = BINARY w.position
WHERE t.id IS NULL;

-- ---------- 分类 ----------
UPDATE tk_lottery_category t
JOIN w_lottery_category w
  ON BINARY t.category_key = BINARY w.category_key
SET
  t.name = w.name,
  t.search_keywords = IFNULL(w.search_keywords, ''),
  t.show_on_home = w.show_on_home,
  t.status = w.status,
  t.sort = w.sort,
  t.updated_at = IFNULL(w.updated_at, t.updated_at);

INSERT INTO tk_lottery_category (
  category_key, name, search_keywords, show_on_home, status, sort, created_at, updated_at
)
SELECT
  w.category_key, w.name, IFNULL(w.search_keywords, ''), w.show_on_home, w.status, w.sort,
  IFNULL(w.created_at, NOW(3)), IFNULL(w.updated_at, NOW(3))
FROM w_lottery_category w
LEFT JOIN tk_lottery_category t
  ON BINARY t.category_key = BINARY w.category_key
WHERE t.id IS NULL;

-- ---------- 图纸 ----------
UPDATE tk_lottery_info t
JOIN w_lottery_info w
  ON t.special_lottery_id = w.special_lottery_id
 AND BINARY t.issue = BINARY w.issue
 AND BINARY t.title = BINARY w.title
SET
  t.category_id = COALESCE((
    SELECT c.id FROM tk_lottery_category c
    WHERE c.category_key COLLATE utf8mb4_general_ci = IFNULL(w.category_tag, '') COLLATE utf8mb4_general_ci
       OR c.name COLLATE utf8mb4_general_ci = IFNULL(w.category_tag, '') COLLATE utf8mb4_general_ci
    ORDER BY c.id ASC LIMIT 1
  ), t.category_id),
  t.category_tag = IFNULL(w.category_tag, ''),
  t.year = w.year,
  t.cover_image_url = w.cover_image_url,
  t.detail_image_url = w.detail_image_url,
  t.draw_code = w.draw_code,
  t.normal_draw_result = TRIM(BOTH ',' FROM SUBSTRING_INDEX(REPLACE(IFNULL(w.draw_result, ''), ' ', ''), ',', 6)),
  t.special_draw_result = TRIM(BOTH ',' FROM SUBSTRING_INDEX(REPLACE(IFNULL(w.draw_result, ''), ' ', ''), ',', -1)),
  t.draw_result = w.draw_result,
  t.draw_at = w.draw_at,
  t.playback_url = IFNULL(t.playback_url, ''),
  t.is_current = w.is_current,
  t.status = w.status,
  t.sort = w.sort,
  t.likes_count = w.likes_count,
  t.comment_count = w.comment_count,
  t.favorite_count = w.favorite_count,
  t.read_count = w.read_count,
  t.poll_enabled = w.poll_enabled,
  t.poll_default_expand = w.poll_default_expand,
  t.recommend_info_ids = IFNULL(w.recommend_info_ids, ''),
  t.updated_at = IFNULL(w.updated_at, t.updated_at);

INSERT INTO tk_lottery_info (
  special_lottery_id, category_id, category_tag, issue, year, title,
  cover_image_url, detail_image_url, draw_code, normal_draw_result, special_draw_result, draw_result, draw_at, playback_url,
  is_current, status, sort, likes_count, comment_count, favorite_count, read_count,
  poll_enabled, poll_default_expand, recommend_info_ids, created_at, updated_at
)
SELECT
  w.special_lottery_id,
  COALESCE((
    SELECT c.id FROM tk_lottery_category c
    WHERE c.category_key COLLATE utf8mb4_general_ci = IFNULL(w.category_tag, '') COLLATE utf8mb4_general_ci
       OR c.name COLLATE utf8mb4_general_ci = IFNULL(w.category_tag, '') COLLATE utf8mb4_general_ci
    ORDER BY c.id ASC LIMIT 1
  ), 0),
  IFNULL(w.category_tag, ''), w.issue, w.year, w.title,
  w.cover_image_url, w.detail_image_url, w.draw_code,
  TRIM(BOTH ',' FROM SUBSTRING_INDEX(REPLACE(IFNULL(w.draw_result, ''), ' ', ''), ',', 6)),
  TRIM(BOTH ',' FROM SUBSTRING_INDEX(REPLACE(IFNULL(w.draw_result, ''), ' ', ''), ',', -1)),
  w.draw_result, w.draw_at, '',
  w.is_current, w.status, w.sort, w.likes_count, w.comment_count, w.favorite_count, w.read_count,
  w.poll_enabled, w.poll_default_expand, IFNULL(w.recommend_info_ids, ''),
  IFNULL(w.created_at, NOW(3)), IFNULL(w.updated_at, NOW(3))
FROM w_lottery_info w
LEFT JOIN tk_lottery_info t
  ON t.special_lottery_id = w.special_lottery_id
 AND BINARY t.issue = BINARY w.issue
 AND BINARY t.title = BINARY w.title
WHERE t.id IS NULL;
