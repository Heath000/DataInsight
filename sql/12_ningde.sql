-- +migrate Up

-- 插入顶级类（一级类）"综合" 到宁德
INSERT INTO categories (category_id, category_name, level, region_id)
VALUES (50, '综合', 1, '9');

-- 获取"综合"类目ID
SET @category_comprehensive = 50;

-- 插入"国民经济主要指标" (二级类，属于"综合") 到宁德
INSERT INTO categories (category_id, category_name, parent_id, level, region_id)
VALUES (51, '国民经济主要指标', @category_comprehensive, 2, '9');

-- 插入"国民经济主要指标发展速度" (二级类，属于"综合") 到宁德
INSERT INTO categories (category_id, category_name, parent_id, level, region_id)
VALUES (52, '国民经济主要指标发展速度', @category_comprehensive, 2, '9');

-- 插入"平均每天主要社会经济获得" (二级类，属于"综合") 到宁德
INSERT INTO categories (category_id, category_name, parent_id, level, region_id)
VALUES (53, '平均每天主要社会经济获得', @category_comprehensive, 2, '9');

-- 插入"国民经济核算" (一级类) 到宁德
INSERT INTO categories (category_id, category_name, level, region_id)
VALUES (54, '国民经济核算', 1, '9');

-- 插入"人口与劳动力" (一级类) 到宁德
INSERT INTO categories (category_id, category_name, level, region_id)
VALUES (55, '人口与劳动力', 1, '9');

-- 插入"人民生活" (一级类) 到宁德
INSERT INTO categories (category_id, category_name, level, region_id)
VALUES (56, '人民生活', 1, '9');

-- +migrate Down

-- 删除"平均每天主要社会经济获得" (二级类，属于"综合") 从宁德
DELETE FROM categories WHERE category_id = 53 AND level = 2 AND region_id = '9';

-- 删除"国民经济主要指标发展速度" (二级类，属于"综合") 从宁德
DELETE FROM categories WHERE category_id = 52 AND level = 2 AND region_id = '9';

-- 删除"国民经济主要指标" (二级类，属于"综合") 从宁德
DELETE FROM categories WHERE category_id = 51 AND level = 2 AND region_id = '9';

-- 删除"综合" (一级类) 从宁德
DELETE FROM categories WHERE category_id = 50 AND level = 1 AND region_id = '9';

-- 删除"国民经济核算" (一级类) 从宁德
DELETE FROM categories WHERE category_id = 54 AND level = 1 AND region_id = '9';

-- 删除"人口与劳动力" (一级类) 从宁德
DELETE FROM categories WHERE category_id = 55 AND level = 1 AND region_id = '9';

-- 删除"人民生活" (一级类) 从宁德
DELETE FROM categories WHERE category_id = 56 AND level = 1 AND region_id = '9';
