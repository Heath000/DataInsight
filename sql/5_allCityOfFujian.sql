-- +migrate Up
INSERT INTO regions (region_id, region_name, province_id) VALUES
('2', '厦门', '1'),
('3', '莆田', '1'),
('4', '三明', '1'),
('5', '泉州', '1'),
('6', '漳州', '1'),
('7', '南平', '1'),
('8', '龙岩', '1'),
('9', '宁德', '1');


-- +migrate Down

DELETE FROM regions WHERE region_id IN ('2', '3', '4', '5', '6', '7', '8', '9');