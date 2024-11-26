package controller

import (
	"net/http"

	"github.com/Heath000/fzuSE2024/model"
	"github.com/gin-gonic/gin"
)

// AdminDataController handles admin-related data operations
type AdminDataController struct{}

// NewAdminDataController initializes a new AdminDataController
func NewAdminDataController() *AdminDataController {
	return &AdminDataController{}
}

// GetProvinceList 获取省份列表
func (ctrl *AdminDataController) GetProvinceList(c *gin.Context) {
	provinces, err := model.Province{}.GetAllProvinces()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch provinces: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success", "data": provinces})
}

// CreateProvince 新增省份
func (ctrl *AdminDataController) CreateProvince(c *gin.Context) {
	var province model.Province
	if err := c.ShouldBindJSON(&province); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	if err := province.Create(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create province: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "success"})
}

// UpdateProvince 修改省份信息
func (ctrl *AdminDataController) UpdateProvince(c *gin.Context) {
	provinceID := c.Param("province_id")
	var province model.Province
	if err := c.ShouldBindJSON(&province); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	if err := province.Update(provinceID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update province: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// DeleteProvince 删除省份
func (ctrl *AdminDataController) DeleteProvince(c *gin.Context) {
	provinceID := c.Param("province_id")

	if err := (&model.Province{}).Delete(provinceID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete province: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// GetRegionList 获取某省份下的区域列表
func (ctrl *AdminDataController) GetRegionList(c *gin.Context) {
	provinceID := c.Param("province_id")

	regions, err := model.Region{}.GetRegionsByProvince(provinceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch regions: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": regions})
}

// CreateRegion 新增区域
func (ctrl *AdminDataController) CreateRegion(c *gin.Context) {
	var region model.Region
	if err := c.ShouldBindJSON(&region); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	if err := region.Create(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create region: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "success"})
}

// UpdateRegion 修改区域信息
func (ctrl *AdminDataController) UpdateRegion(c *gin.Context) {
	regionID := c.Param("region_id")
	var region model.Region
	if err := c.ShouldBindJSON(&region); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	if err := region.Update(regionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update region: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// DeleteRegion 删除区域
func (ctrl *AdminDataController) DeleteRegion(c *gin.Context) {
	regionID := c.Param("region_id")

	if err := (&model.Region{}).Delete(regionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete region: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// GetTopCategories 根据地区查询顶级分类
func (ctrl *AdminDataController) GetTopCategories(c *gin.Context) {
	regionID := c.Param("region_id")

	categories, err := model.Category{}.GetTopLevelCategoriesByRegion(regionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch top categories: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": categories})
}

// GetCategoryList 根据分类查询子分类
func (ctrl *AdminDataController) GetCategoryList(c *gin.Context) {
	categoryID := c.Param("category_id")

	subCategories, err := model.Category{}.GetSubCategories(categoryID)
	if err != nil {
		if err == model.ErrLeafCategory {
			c.JSON(http.StatusNotFound, gin.H{"message": "Leaf category, no subcategories", "data": []model.Category{}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch subcategories: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": subCategories})
}
// CreateCategory 新增分类
func (ctrl *AdminDataController) CreateCategory(c *gin.Context) {
	var category model.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	if err := category.Create(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Category created successfully"})
}
// UpdateCategory 修改分类信息
func (ctrl *AdminDataController) UpdateCategory(c *gin.Context) {
	categoryID := c.Param("category_id")
	var category model.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	if err := category.Update(categoryID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category updated successfully"})
}
// DeleteCategory 删除分类
func (ctrl *AdminDataController) DeleteCategory(c *gin.Context) {
	categoryID := c.Param("category_id")

	if err := (&model.Category{}).Delete(categoryID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}

// GetDataList 获取分类下的数据
func (ctrl *AdminDataController) GetDataList(c *gin.Context) {
	categoryID := c.Param("category_id")

	// 假设从请求中接收一个包含年份的参数数组
	var request struct {
		Years []string `json:"years"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
		return
	}

	data, err := model.BasicData{}.GetBasicDataByCategoryAndYears(categoryID, request.Years)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": data})
}


// CreateData 新增数据
func (ctrl *AdminDataController) CreateData(c *gin.Context) {
	var data model.BasicData
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	if err := data.Create(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create data: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "success"})
}

// UpdateData 修改数据
func (ctrl *AdminDataController) UpdateData(c *gin.Context) {
	regionID := c.Param("region_id")
	categoryID := c.Param("category_id")
	year := c.Param("year")

	var data model.BasicData
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	if err := data.Update(regionID, categoryID, year); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update data: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// DeleteData 删除数据
func (ctrl *AdminDataController) DeleteData(c *gin.Context) {
	regionID := c.Param("region_id")
	categoryID := c.Param("category_id")
	year := c.Param("year")

	if err := (&model.BasicData{}).Delete(regionID, categoryID, year); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete data: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
