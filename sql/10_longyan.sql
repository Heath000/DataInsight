-- +migrate Up

-- 插入顶级类（一级类）"综合" 到龙岩
INSERT INTO categories (category_id, category_name, level, region_id)
VALUES (36, '综合', 1, '8');

-- 获取"综合"类目ID
SET @category_comprehensive = 36;

-- 插入"国民经济主要指标" (二级类，属于"综合") 到龙岩
INSERT INTO categories (category_id, category_name, parent_id, level, region_id)
VALUES (37, '国民经济主要指标', @category_comprehensive, 2, '8');

-- 插入"国民经济主要指标发展速度" (二级类，属于"综合") 到龙岩
INSERT INTO categories (category_id, category_name, parent_id, level, region_id)
VALUES (38, '国民经济主要指标发展速度', @category_comprehensive, 2, '8');

-- 插入"平均每天主要社会经济获得" (二级类，属于"综合") 到龙岩
INSERT INTO categories (category_id, category_name, parent_id, level, region_id)
VALUES (39, '平均每天主要社会经济获得', @category_comprehensive, 2, '8');

-- 插入"国民经济核算" (一级类) 到龙岩
INSERT INTO categories (category_id, category_name, level, region_id)
VALUES (40, '国民经济核算', 1, '8');

-- 插入"人口与劳动力" (一级类) 到龙岩
INSERT INTO categories (category_id, category_name, level, region_id)
VALUES (41, '人口与劳动力', 1, '8');

-- 插入"人民生活" (一级类) 到龙岩
INSERT INTO categories (category_id, category_name, level, region_id)
VALUES (42, '人民生活', 1, '8');

-- +migrate Down

-- 删除"平均每天主要社会经济获得" (二级类，属于"综合") 从龙岩
DELETE FROM categories WHERE category_id = 39 AND level = 2 AND region_id = '8';

-- 删除"国民经济主要指标发展速度" (二级类，属于"综合") 从龙岩
DELETE FROM categories WHERE category_id = 38 AND level = 2 AND region_id = '8';

-- 删除"国民经济主要指标" (二级类，属于"综合") 从龙岩
DELETE FROM categories WHERE category_id = 37 AND level = 2 AND region_id = '8';

-- 删除"综合" (一级类) 从龙岩
DELETE FROM categories WHERE category_id = 36 AND level = 1 AND region_id = '8';

-- 删除"国民经济核算" (一级类) 从龙岩
DELETE FROM categories WHERE category_id = 40 AND level = 1 AND region_id = '8';

-- 删除"人口与劳动力" (一级类) 从龙岩
DELETE FROM categories WHERE category_id = 41 AND level = 1 AND region_id = '8';

-- 删除"人民生活" (一级类) 从龙岩
DELETE FROM categories WHERE category_id = 42 AND level = 1 AND region_id = '8';
