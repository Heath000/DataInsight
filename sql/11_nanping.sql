-- +migrate Up

-- 插入顶级类（一级类）"综合" 到南平
INSERT INTO categories (category_id, category_name, level, region_id)
VALUES (100, '综合', 1, '9');

-- 获取"综合"类目ID
SET @category_comprehensive = 100;

-- 插入"国民经济主要指标" (二级类，属于"综合") 到南平
INSERT INTO categories (category_id, category_name, parent_id, level, region_id)
VALUES (101, '国民经济主要指标', @category_comprehensive, 2, '9');

-- 插入"国民经济主要指标发展速度" (二级类，属于"综合") 到南平
INSERT INTO categories (category_id, category_name, parent_id, level, region_id)
VALUES (102, '国民经济主要指标发展速度', @category_comprehensive, 2, '9');

-- 插入"平均每天主要社会经济获得" (二级类，属于"综合") 到南平
INSERT INTO categories (category_id, category_name, parent_id, level, region_id)
VALUES (103, '平均每天主要社会经济获得', @category_comprehensive, 2, '9');

-- 插入"国民经济核算" (一级类) 到南平
INSERT INTO categories (category_id, category_name, level, region_id)
VALUES (104, '国民经济核算', 1, '9');

-- 插入"人口与劳动力" (一级类) 到南平
INSERT INTO categories (category_id, category_name, level, region_id)
VALUES (105, '人口与劳动力', 1, '9');

-- 插入"人民生活" (一级类) 到南平
INSERT INTO categories (category_id, category_name, level, region_id)
VALUES (106, '人民生活', 1, '9');

-- +migrate Down

-- 删除"平均每天主要社会经济获得" (二级类，属于"综合") 从南平
DELETE FROM categories WHERE category_id = 103 AND level = 2 AND region_id = '9';

-- 删除"国民经济主要指标发展速度" (二级类，属于"综合") 从南平
DELETE FROM categories WHERE category_id = 102 AND level = 2 AND region_id = '9';

-- 删除"国民经济主要指标" (二级类，属于"综合") 从南平
DELETE FROM categories WHERE category_id = 101 AND level = 2 AND region_id = '9';

-- 删除"综合" (一级类) 从南平
DELETE FROM categories WHERE category_id = 100 AND level = 1 AND region_id = '9';

-- 删除"国民经济核算" (一级类) 从南平
DELETE FROM categories WHERE category_id = 104 AND level = 1 AND region_id = '9';

-- 删除"人口与劳动力" (一级类) 从南平
DELETE FROM categories WHERE category_id = 105 AND level = 1 AND region_id = '9';

-- 删除"人民生活" (一级类) 从南平
DELETE FROM categories WHERE category_id = 106 AND level = 1 AND region_id = '9';
