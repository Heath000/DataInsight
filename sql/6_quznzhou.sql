-- +migrate Up

-- 插入顶级类（一级类）"综合" 到泉州
INSERT INTO categories (category_id, category_name, level, region_id)
VALUES (22, '综合', 1, '5');

-- 插入"国民经济主要指标" (二级类，属于"综合") 到泉州
INSERT INTO categories (category_id, category_name, parent_id, level, region_id)
SELECT 23, '国民经济主要指标', c.category_id, 2, '5'
FROM categories c
WHERE c.category_name = '综合' AND c.level = 1 AND c.region_id = '5';

-- 插入"国民经济主要指标发展速度" (二级类，属于"综合") 到泉州
INSERT INTO categories (category_id, category_name, parent_id, level, region_id)
SELECT 24, '国民经济主要指标发展速度', c.category_id, 2, '5'
FROM categories c
WHERE c.category_name = '综合' AND c.level = 1 AND c.region_id = '5';

-- 插入"平均每天主要社会经济获得" (二级类，属于"综合") 到泉州
INSERT INTO categories (category_id, category_name, parent_id, level, region_id)
SELECT 25, '平均每天主要社会经济获得', c.category_id, 2, '5'
FROM categories c
WHERE c.category_name = '综合' AND c.level = 1 AND c.region_id = '5';

-- 插入"国民经济核算" (一级类) 到泉州
INSERT INTO categories (category_id, category_name, level, region_id)
VALUES (26, '国民经济核算', 1, '5');

-- 插入"人口与劳动力" (一级类) 到泉州
INSERT INTO categories (category_id, category_name, level, region_id)
VALUES (27, '人口与劳动力', 1, '5');

-- 插入"人民生活" (一级类) 到泉州
INSERT INTO categories (category_id, category_name, level, region_id)
VALUES (28, '人民生活', 1, '5');

-- +migrate Down

-- 删除"平均每天主要社会经济获得" (二级类，属于"综合") 从泉州
DELETE FROM categories WHERE category_id = 25 AND level = 2 AND region_id = '5';

-- 删除"国民经济主要指标发展速度" (二级类，属于"综合") 从泉州
DELETE FROM categories WHERE category_id = 24 AND level = 2 AND region_id = '5';

-- 删除"国民经济主要指标" (二级类，属于"综合") 从泉州
DELETE FROM categories WHERE category_id = 23 AND level = 2 AND region_id = '5';

-- 删除"综合" (一级类) 从泉州
DELETE FROM categories WHERE category_id = 22 AND level = 1 AND region_id = '5';

-- 删除"国民经济核算" (一级类) 从泉州
DELETE FROM categories WHERE category_id = 26 AND level = 1 AND region_id = '5';

-- 删除"人口与劳动力" (一级类) 从泉州
DELETE FROM categories WHERE category_id = 27 AND level = 1 AND region_id = '5';

-- 删除"人民生活" (一级类) 从泉州
DELETE FROM categories WHERE category_id = 28 AND level = 1 AND region_id = '5';
