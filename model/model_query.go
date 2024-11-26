package model

import (
	"errors"
)

// Province 表结构
type Province struct {
	ProvinceID   string   `gorm:"type:char(10);primaryKey" json:"province_id"`        // 省份ID
	ProvinceName string   `gorm:"type:char(20);not null;unique" json:"province_name"` // 省份名称，非空且唯一
	Regions      []Region `gorm:"foreignKey:ProvinceID" json:"-"`                     // 与地区的关系，一对多
}

// Region 表结构
type Region struct {
	RegionID   string     `gorm:"type:char(10);primaryKey" json:"region_id"`                  // 地区ID
	RegionName string     `gorm:"type:char(20);not null" json:"region_name"`                  // 地区名称，非空
	ProvinceID string     `gorm:"type:char(10);not null" json:"province_id"`                  // 省份ID，外键关联到省份表
	Province   Province   `gorm:"foreignKey:ProvinceID;constraint:OnDelete:CASCADE" json:"-"` // 外键设置级联删除
	Categories []Category `gorm:"foreignKey:RegionID" json:"-"`                               // 与类别的关系，一对多
}

// Category 表结构
type Category struct {
	CategoryID   string    `gorm:"type:char(10);primaryKey" json:"category_id"`               // 类别ID
	ParentID     *string   `gorm:"type:char(10)" json:"parent_id"`                            // 父类别ID，允许为空
	CategoryName string    `gorm:"type:char(20);not null" json:"category_name"`               // 类别名称，非空
	Level        int       `gorm:"type:int" json:"level"`                                     // 类别层级
	RegionID     string    `gorm:"type:char(10);not null" json:"region_id"`                   // 地区ID，外键关联到地区表
	Region       Region    `gorm:"foreignKey:RegionID;constraint:OnDelete:CASCADE" json:"-"`  // 外键设置级联删除
	Parent       *Category `gorm:"foreignKey:ParentID;constraint:OnDelete:SET NULL" json:"-"` // 自关联，设置父类，删除置空
}

// BasicData 表结构
type BasicData struct {
	CategoryID string   `gorm:"type:char(10);not null" json:"category_id"`                  // 分类ID，外键关联到类别表
	DataName   string   `gorm:"type:char(20);not null" json:"data_name"`                    // 数据名称，非空
	Data       int      `gorm:"type:int;default:0;check:data >= 0" json:"data"`             // 数据，默认为0且不允许负数
	Year       string   `gorm:"type:char(4);not null" json:"year"`                          // 年份，非空，且应是4位字符
	Category   Category `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE" json:"-"` // 外键设置级联删除
}

// 查询接口1：根据省份 ID 获取其下的地区
func (Region) GetRegionsByProvince(provinceID string) ([]Region, error) {
	var regions []Region
	err := DB().Where("province_id = ?", provinceID).Find(&regions).Error
	if err != nil {
		return nil, err
	}
	return regions, nil
}

// 查询接口2：根据地区 ID 获取顶级类别（parent_id 为空的类别）
func (Category) GetTopLevelCategoriesByRegion(regionID string) ([]Category, error) {
	var categories []Category
	err := DB().Where("region_id = ? AND parent_id IS NULL", regionID).Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// 查询接口3：根据类别 ID 获取其下的子类别，若无子类别则返回叶子类别错误
var ErrLeafCategory = errors.New("no subcategories available; this is a leaf category")

func (Category) GetSubCategories(categoryID string) ([]Category, error) {
	var subCategories []Category

	// 查询该分类的下级分类
	err := DB().Where("parent_id = ?", categoryID).Find(&subCategories).Error
	if err != nil {
		return nil, err // 返回数据库查询错误
	}

	// 如果没有下级分类，返回 ErrLeafCategory 错误
	if len(subCategories) == 0 {
		return nil, ErrLeafCategory
	}

	// 如果有下级分类，直接返回下级分类
	return subCategories, nil
}

// 查询接口4：根据类别 ID 和年份范围查询 BasicData 中的相关数据
func (BasicData) GetBasicDataByCategoryAndYears(categoryID string, years []string) ([]BasicData, error) {
	var basicData []BasicData
	err := DB().Where("category_id = ? AND year IN ?", categoryID, years).Find(&basicData).Error
	if err != nil {
		return nil, err
	}
	return basicData, nil
}

// 查询接口5:查询这个叶子category下有多少年份可以查
// GetAvailableYearsByLeafCategory 根据叶子类别 ID 查询该类别下可查询的年份列表

// GetAllProvinces 获取所有省份信息
func (Province) GetAllProvinces() ([]Province, error) {
	var provinces []Province

	// 查询数据库中所有省份信息
	err := DB().Find(&provinces).Error
	if err != nil {
		return nil, err
	}

	return provinces, nil
}
func (BasicData) GetAvailableYearsByLeafCategory(categoryID string) ([]string, error) {
	var years []string

	// 查询 BasicData 表中指定 category_id 的所有不同年份
	err := DB().Model(&BasicData{}).
		Select("DISTINCT year").
		Where("category_id = ?", categoryID).
		Order("year ASC").
		Pluck("year", &years).Error

	if err != nil {
		return nil, err
	}

	return years, nil
}


//以下是开放给管理员，可以对数据进行修改的接口
// Province 表相关方法

// Create 新增省份
func (p *Province) Create() error {
	return DB().Create(p).Error
}

// Update 修改省份
func (p *Province) Update(provinceID string) error {
	return DB().Model(&Province{}).Where("province_id = ?", provinceID).Updates(p).Error
}

// Delete 删除省份
func (Province) Delete(provinceID string) error {
	return DB().Where("province_id = ?", provinceID).Delete(&Province{}).Error
}

// Region 表相关方法

// Create 新增区域
func (r *Region) Create() error {
	return DB().Create(r).Error
}

// Update 修改区域
func (r *Region) Update(regionID string) error {
	return DB().Model(&Region{}).Where("region_id = ?", regionID).Updates(r).Error
}

// Delete 删除区域
func (Region) Delete(regionID string) error {
	return DB().Where("region_id = ?", regionID).Delete(&Region{}).Error
}

// Category 表相关方法

// Create 新增分类
func (c *Category) Create() error {
	return DB().Create(c).Error
}

// Update 修改分类
func (c *Category) Update(categoryID string) error {
	return DB().Model(&Category{}).Where("category_id = ?", categoryID).Updates(c).Error
}

// Delete 删除分类
func (Category) Delete(categoryID string) error {
	return DB().Where("category_id = ?", categoryID).Delete(&Category{}).Error
}

// BasicData 表相关方法

// Create 新增数据
func (d *BasicData) Create() error {
	return DB().Create(d).Error
}

// Update 修改数据
func (d *BasicData) Update(regionID, categoryID, year string) error {
	return DB().
		Model(&BasicData{}).
		Where("region_id = ? AND category_id = ? AND year = ?", regionID, categoryID, year).
		Updates(d).Error
}

// Delete 删除数据
func (BasicData) Delete(regionID, categoryID, year string) error {
	return DB().
		Where("region_id = ? AND category_id = ? AND year = ?", regionID, categoryID, year).
		Delete(&BasicData{}).Error
}