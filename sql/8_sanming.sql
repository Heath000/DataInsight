-- +migrate Up

-- 插入顶级类（一级类）"综合" 到三明
INSERT INTO categories (category_id, category_name, level, region_id)
VALUES (43, '综合', 1, '8');

-- 获取"综合"类目ID
SET @category_comprehensive = 43;

-- 插入"国民经济主要指标" (二级类，属于"综合") 到三明
INSERT INTO categories (category_id, category_name, parent_id, level, region_id)
VALUES (44, '国民经济主要指标', @category_comprehensive, 2, '8');

-- 插入"国民经济主要指标发展速度" (二级类，属于"综合") 到三明
INSERT INTO categories (category_id, category_name, parent_id, level, region_id)
VALUES (45, '国民经济主要指标发展速度', @category_comprehensive, 2, '8');

-- 插入"平均每天主要社会经济获得" (二级类，属于"综合") 到三明
INSERT INTO categories (category_id, category_name, parent_id, level, region_id)
VALUES (46, '平均每天主要社会经济获得', @category_comprehensive, 2, '8');

-- 插入"国民经济核算" (一级类) 到三明
INSERT INTO categories (category_id, category_name, level, region_id)
VALUES (47, '国民经济核算', 1, '8');

-- 插入"人口与劳动力" (一级类) 到三明
INSERT INTO categories (category_id, category_name, level, region_id)
VALUES (48, '人口与劳动力', 1, '8');

-- 插入"人民生活" (一级类) 到三明
INSERT INTO categories (category_id, category_name, level, region_id)
VALUES (49, '人民生活', 1, '8');

-- +migrate Down

-- 删除"平均每天主要社会经济获得" (二级类，属于"综合") 从三明
DELETE FROM categories WHERE category_id = 46 AND level = 2 AND region_id = '8';

-- 删除"国民经济主要指标发展速度" (二级类，属于"综合") 从三明
DELETE FROM categories WHERE category_id = 45 AND level = 2 AND region_id = '8';

-- 删除"国民经济主要指标" (二级类，属于"综合") 从三明
DELETE FROM categories WHERE category_id = 44 AND level = 2 AND region_id = '8';

-- 删除"综合" (一级类) 从三明
DELETE FROM categories WHERE category_id = 43 AND level = 1 AND region_id = '8';

-- 删除"国民经济核算" (一级类) 从三明
DELETE FROM categories WHERE category_id = 47 AND level = 1 AND region_id = '8';

-- 删除"人口与劳动力" (一级类) 从三明
DELETE FROM categories WHERE category_id = 48 AND level = 1 AND region_id = '8';

-- 删除"人民生活" (一级类) 从三明
DELETE FROM categories WHERE category_id = 49 AND level = 1 AND region_id = '8';
