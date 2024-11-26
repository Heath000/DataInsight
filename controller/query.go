package controller

import (
	"net/http"

	"github.com/Heath000/fzuSE2024/model"
	"github.com/gin-gonic/gin"
)

// QueryController is the query controller
type QueryController struct{}

// 规定请求体 JSON 结构
// ProvinceId 请求体，用于传递省份ID
type ProvinceId struct {
	ProvinceId string `form:"province_id" json:"province_id" binding:"required"`
}

// RegionId 请求体，用于传递地区ID
type RegionId struct {
	RegionId string `form:"region_id" json:"region_id" binding:"required"`
}

// CategoryId 请求体，用于传递分类ID
type CategoryId struct {
	CategoryId string `form:"category_id" json:"category_id" binding:"required"`
}

// CategoryIdAndYear 请求体，用于传递分类ID和年份列表
type CategoryIdAndYear struct {
	CategoryId string   `form:"category_id" json:"category_id" binding:"required"` // 分类ID
	Years      []string `form:"years" json:"years" binding:"required"`             // 年份列表
}

// NewQueryController initializes a new QueryController
func NewQueryController() *QueryController {
	return &QueryController{}
}
func (ctrl *QueryController) GetProvinceList(c *gin.Context) {
	provinces, err := model.Province{}.GetAllProvinces()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch provinces: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success", "data": provinces})
}
// QueryRegions 获取指定省份下的所有地区
func (ctrl *QueryController) QueryRegions(c *gin.Context) {
	// 绑定请求体的参数
	var form ProvinceId
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查 ProvinceId 是否有效
	if form.ProvinceId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ProvinceId is required"})
		return
	}

	// 调用模型方法获取地区列表
	regions, err := model.Region{}.GetRegionsByProvince(form.ProvinceId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch regions: " + err.Error()})
		return
	}

	// 返回 JSON 格式的地区列表
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    regions,
	})
}

// QueryTopCategories 获取指定地区的顶级分类
func (ctrl *QueryController) QueryTopCategories(c *gin.Context) {
	// 绑定请求体，获取 region_id
	var regionID RegionId
	if err := c.ShouldBindJSON(&regionID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// 调用模型方法查询顶级分类
	categories, err := model.Category{}.GetTopLevelCategoriesByRegion(regionID.RegionId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch top categories: " + err.Error()})
		return
	}

	// 返回 JSON 格式的顶级分类列表
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    categories,
	})
}

// QuerySubCategories 获取指定分类的子分类
func (ctrl *QueryController) QuerySubCategories(c *gin.Context) {
	// 绑定请求体，获取 category_id
	var categoryID CategoryId
	if err := c.ShouldBindJSON(&categoryID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// 调用模型方法查询子分类
	subCategories, err := model.Category{}.GetSubCategories(categoryID.CategoryId)
	if err != nil {
		if err == model.ErrLeafCategory {
			c.JSON(http.StatusOK, gin.H{"message": "Leaf category, no subcategories", "data": []model.Category{}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch subcategories: " + err.Error()})
		return
	}

	// 返回 JSON 格式的子分类列表
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    subCategories,
	})
}
//获取叶子category里面可选的年份
func (ctrl *QueryController) QueryAvailableYears(c *gin.Context) {
    var request CategoryId
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
        return
    }

    // 调用模型方法获取年份列表
    years, err := model.BasicData{}.GetAvailableYearsByLeafCategory(request.CategoryId)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch available years: " + err.Error()})
        return
    }

    // 返回 JSON 格式的年份列表
    c.JSON(http.StatusOK, gin.H{
        "message": "success",
        "years":   years,
    })
}

// QueryData 获取指定分类在指定年份范围内的时间序列数据
func (ctrl *QueryController) QueryData(c *gin.Context) {
	// 绑定请求体，获取 category_id 和 years
	var requestData CategoryIdAndYear
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// 调用模型方法查询时间序列数据
	basicData, err := model.BasicData{}.GetBasicDataByCategoryAndYears(requestData.CategoryId, requestData.Years)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch basic data: " + err.Error()})
		return
	}

	// 返回 JSON 格式的时间序列数据
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    basicData,
	})
}
