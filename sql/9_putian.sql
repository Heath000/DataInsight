-- +migrate Up

-- 插入顶级类（一级类）"综合" 到莆田
INSERT INTO categories (category_id, category_name, level, region_id)
VALUES (29, '综合', 1, '7');

-- 获取"综合"类目ID
SET @category_comprehensive = 29;

-- 插入"国民经济主要指标" (二级类，属于"综合") 到莆田
INSERT INTO categories (category_id, category_name, parent_id, level, region_id)
VALUES (30, '国民经济主要指标', @category_comprehensive, 2, '7');

-- 插入"国民经济主要指标发展速度" (二级类，属于"综合") 到莆田
INSERT INTO categories (category_id, category_name, parent_id, level, region_id)
VALUES (31, '国民经济主要指标发展速度', @category_comprehensive, 2, '7');

-- 插入"平均每天主要社会经济获得" (二级类，属于"综合") 到莆田
INSERT INTO categories (category_id, category_name, parent_id, level, region_id)
VALUES (32, '平均每天主要社会经济获得', @category_comprehensive, 2, '7');

-- 插入"国民经济核算" (一级类) 到莆田
INSERT INTO categories (category_id, category_name, level, region_id)
VALUES (33, '国民经济核算', 1, '7');

-- 插入"人口与劳动力" (一级类) 到莆田
INSERT INTO categories (category_id, category_name, level, region_id)
VALUES (34, '人口与劳动力', 1, '7');

-- 插入"人民生活" (一级类) 到莆田
INSERT INTO categories (category_id, category_name, level, region_id)
VALUES (35, '人民生活', 1, '7');

-- +migrate Down

-- 删除"平均每天主要社会经济获得" (二级类，属于"综合") 从莆田
DELETE FROM categories WHERE category_id = 32 AND level = 2 AND region_id = '7';

-- 删除"国民经济主要指标发展速度" (二级类，属于"综合") 从莆田
DELETE FROM categories WHERE category_id = 31 AND level = 2 AND region_id = '7';

-- 删除"国民经济主要指标" (二级类，属于"综合") 从莆田
DELETE FROM categories WHERE category_id = 30 AND level = 2 AND region_id = '7';

-- 删除"综合" (一级类) 从莆田
DELETE FROM categories WHERE category_id = 29 AND level = 1 AND region_id = '7';

-- 删除"国民经济核算" (一级类) 从莆田
DELETE FROM categories WHERE category_id = 33 AND level = 1 AND region_id = '7';

-- 删除"人口与劳动力" (一级类) 从莆田
DELETE FROM categories WHERE category_id = 34 AND level = 1 AND region_id = '7';

-- 删除"人民生活" (一级类) 从莆田
DELETE FROM categories WHERE category_id = 35 AND level = 1 AND region_id = '7';
