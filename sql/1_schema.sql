
-- +migrate Up

CREATE TABLE provinces (
    province_id CHAR(10) PRIMARY KEY,
    province_name CHAR(20) NOT NULL UNIQUE
) DEFAULT CHARACTER SET utf8mb4;

CREATE TABLE regions (
    region_id CHAR(10) PRIMARY KEY,
    region_name CHAR(20) NOT NULL,
    province_id CHAR(10),
    FOREIGN KEY (province_id) REFERENCES provinces(province_id) ON DELETE CASCADE
) DEFAULT CHARACTER SET utf8mb4;

CREATE TABLE categories (
    category_id CHAR(10) PRIMARY KEY,
    parent_id CHAR(10),
    category_name CHAR(50) NOT NULL,
    level INT,
    region_id CHAR(10),
    FOREIGN KEY (parent_id) REFERENCES categories(category_id) ON DELETE SET NULL,
    FOREIGN KEY (region_id) REFERENCES regions(region_id) ON DELETE CASCADE
) DEFAULT CHARACTER SET utf8mb4;

CREATE TABLE basic_data (
    region_id CHAR(10),
    category_id CHAR(10),
    data_name CHAR(50) NOT NULL,
    data INT DEFAULT 0,
    year CHAR(4) NOT NULL,
    PRIMARY KEY (region_id, category_id, year),
    FOREIGN KEY (region_id) REFERENCES regions(region_id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(category_id) ON DELETE CASCADE
) DEFAULT CHARACTER SET utf8mb4;


-- +migrate Down
DROP TABLE IF EXISTS basic_data;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS regions;
DROP TABLE IF EXISTS provinces;