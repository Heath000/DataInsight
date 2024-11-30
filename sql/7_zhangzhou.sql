-- +migrate Up

-- 插入顶级类（一级类）"综合" 到漳州
INSERT INTO categories (category_id, category_name, level, region_id)
VALUES (57, '综合', 1, '6');

-- 获取"综合"类目ID
SET @category_comprehensive = 57;

-- 插入"国民经济主要指标" (二级类，属于"综合") 到漳州
INSERT INTO categories (category_id, category_name, parent_id, level, region_id)
VALUES (58, '国民经济主要指标', @category_comprehensive, 2, '6');

-- 插入"国民经济主要指标发展速度" (二级类，属于"综合") 到漳州
INSERT INTO categories (category_id, category_name, parent_id, level, region_id)
VALUES (59, '国民经济主要指标发展速度', @category_comprehensive, 2, '6');

-- 插入"平均每天主要社会经济获得" (二级类，属于"综合") 到漳州
INSERT INTO categories (category_id, category_name, parent_id, level, region_id)
VALUES (60, '平均每天主要社会经济获得', @category_comprehensive, 2, '6');

-- 插入"国民经济核算" (一级类) 到漳州
INSERT INTO categories (category_id, category_name, level, region_id)
VALUES (61, '国民经济核算', 1, '6');

-- 插入"人口与劳动力" (一级类) 到漳州
INSERT INTO categories (category_id, category_name, level, region_id)
VALUES (62, '人口与劳动力', 1, '6');

-- 插入"人民生活" (一级类) 到漳州
INSERT INTO categories (category_id, category_name, level, region_id)
VALUES (63, '人民生活', 1, '6');

-- +migrate Down

-- 删除"平均每天主要社会经济获得" (二级类，属于"综合") 从漳州
DELETE FROM categories WHERE category_id = 60 AND level = 2 AND region_id = '6';

-- 删除"国民经济主要指标发展速度" (二级类，属于"综合") 从漳州
DELETE FROM categories WHERE category_id = 59 AND level = 2 AND region_id = '6';

-- 删除"国民经济主要指标" (二级类，属于"综合") 从漳州
DELETE FROM categories WHERE category_id = 58 AND level = 2 AND region_id = '6';

-- 删除"综合" (一级类) 从漳州
DELETE FROM categories WHERE category_id = 57 AND level = 1 AND region_id = '6';

-- 删除"国民经济核算" (一级类) 从漳州
DELETE FROM categories WHERE category_id = 61 AND level = 1 AND region_id = '6';

-- 删除"人口与劳动力" (一级类) 从漳州
DELETE FROM categories WHERE category_id = 62 AND level = 1 AND region_id = '6';

-- 删除"人民生活" (一级类) 从漳州
DELETE FROM categories WHERE category_id = 63 AND level = 1 AND region_id = '6';
