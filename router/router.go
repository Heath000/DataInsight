package router

import (
	"fmt"

	"github.com/Heath000/fzuSE2024/controller"
	"github.com/Heath000/fzuSE2024/middleware"
	"github.com/Heath000/fzuSE2024/model"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// Route makes the routing
func Route(app *gin.Engine) {
	fmt.Println("Routes are registered")
	indexController := new(controller.IndexController)
	app.GET(
		"/", indexController.GetIndex,
	)

	auth := app.Group("/auth")
	authMiddleware := middleware.Auth()
	auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/hello", func(c *gin.Context) {
			claims := jwt.ExtractClaims(c)
			user, _ := c.Get("email")
			c.JSON(200, gin.H{
				"email": claims["email"],
				"name":  user.(*model.User).Name,
				"text":  "Hello World.",
			})
		})
	}
	// 管理员路由组,能够修改用户，内置区域经济两大数据库
	adminUserController := new(controller.AdminUserController)
	adminUser := app.Group("/admin/user")
	adminUser.Use(middleware.Auth().MiddlewareFunc()) // 挂载鉴权中间件,含"/admin"路由的只有email="admin"的账号才被通过
	{
		adminUser.GET("/userlist", adminUserController.GetUserList)          // 获取用户列表
		adminUser.GET("/get_user/:id", adminUserController.GetUser)          //获取指定用户信息byId
		adminUser.POST("/create_user", adminUserController.CreateUser)       // 管理员创建用户
		adminUser.DELETE("/delete_user/:id", adminUserController.DeleteUser) // 管理员删除用户
		adminUser.PUT("/update_user", adminUserController.UpdateUser)        //管理员修改用户信息
	}
	userController := new(controller.UserController)
	app.GET(
		"/user/:id", userController.GetUser,
	).GET(
		"/signup", userController.SignupForm,
	).POST(
		"/signup", userController.Signup,
	).GET(
		"/login", userController.LoginForm,
	).POST(
		"/login", authMiddleware.LoginHandler,
	)

	api := app.Group("/api")
	{
		api.GET("/version", indexController.GetVersion)
	}

	chartController := new(controller.ChartController)
	chart := app.Group("/chart")
	{
		chart.POST("/pie", chartController.ChartPie)
		chart.POST("/bar", chartController.ChartBar)
		chart.POST("/line", chartController.ChartLine)
		chart.POST("/linebarmixed", chartController.ChartLineBarMixed)
	}

	analysisController := new(controller.AnalysisController)
	analysis := app.Group("/analysis")
	{
		analysis.POST("/overall", analysisController.AnalysisOverview)
		analysis.POST("/linear_regression", analysisController.AnalysisLinearRegress)
		analysis.POST("/grey_predict", analysisController.AnalysisGreyPredict)
		analysis.POST("/ARIMA", analysisController.AnalysisARIMA)
		analysis.POST("/BP", analysisController.AnalysisBP)
		analysis.POST("/SVM", analysisController.AnalysisSVM)
		analysis.POST("/RandomForest", analysisController.AnalysisRandomForest)
	}

	dataProcessingController := new(controller.DataProcessingController)
	data := app.Group("/dataProcessing")
	{
		data.POST("/standardize", dataProcessingController.DataStandalize)
		data.POST("/outliers", dataProcessingController.DataOutliersHandle)
		data.POST("/missing", dataProcessingController.DataMissingValuesHandle)
		data.POST("/feature/variance", dataProcessingController.DataFeatureVariance)
		data.POST("/feature/correlation", dataProcessingController.DataFeatureCorrelation)
		data.POST("/feature/chi_square", dataProcessingController.DataFeatureChiSquare)
	}

	queryController := new(controller.QueryController)
	query := app.Group("/query")
	{
		query.GET("/provinces", queryController.GetProvinceList)
		query.POST("/region", queryController.QueryRegions)
		query.POST("/top_category", queryController.QueryTopCategories)
		query.POST("/sub_category", queryController.QuerySubCategories)
		query.POST("/available_year", queryController.QueryAvailableYears)
		query.POST("/data", queryController.QueryData)

	}

	llmController := new(controller.LlmController)
	llm := app.Group("/llm")
	{
		llm.POST("/report", llmController.GetReport)
	}

	fileController := new(controller.FileController)
	file := app.Group("/file")
	file.Use(authMiddleware.MiddlewareFunc()) // 所有文件相关的路由都需要鉴权
	{
		file.GET("/get_file_list", fileController.GetFileList)          // 获取文件列表
		file.GET("/get_file/:file_id", fileController.GetFile)          // 获取单个文件
		file.DELETE("/delete_file/:file_id", fileController.DeleteFile) // 删除文件
		file.POST("/upload_file", fileController.UploadFile)            // 上传文件
	}

	adminDataController := new(controller.AdminDataController)
	adminData := app.Group("/admin/data")
	adminData.Use(middleware.Auth().MiddlewareFunc()) // 挂载鉴权中间件,含"/admin"路由的只有email="admin"的账号才被通过
	{
		// 省份相关接口
		//adminData.GET("/provinces", adminDataController.GetProvinceList)  // 获取省份列表
		adminData.POST("/province", adminDataController.CreateProvince)                // 新增省份
		adminData.PUT("/province/:province_id", adminDataController.UpdateProvince)    // 修改省份信息
		adminData.DELETE("/province/:province_id", adminDataController.DeleteProvince) // 删除省份

		// 区域相关接口
		//adminData.GET("/regions/:province_id", adminDataController.GetRegionList)  // 获取某省份下的区域列表
		adminData.POST("/region", adminDataController.CreateRegion)              // 新增区域
		adminData.PUT("/region/:region_id", adminDataController.UpdateRegion)    // 修改区域信息
		adminData.DELETE("/region/:region_id", adminDataController.DeleteRegion) // 删除区域

		// 分类相关接口
		//adminData.GET("/top_categories/:region_id", adminDataController.GetTopCategories)  // 获取区域下的顶级分类
		//adminData.GET("/sub_categories/:region_id", adminDataController.GetCategoryList)  // 获取这个分类的子分类
		adminData.POST("/category", adminDataController.CreateCategory)                // 新增分类
		adminData.PUT("/category/:category_id", adminDataController.UpdateCategory)    // 修改分类信息
		adminData.DELETE("/category/:category_id", adminDataController.DeleteCategory) // 删除分类

		// 数据相关接口
		//adminData.GET("/data/:region_id/:category_id", adminDataController.GetDataList)  // 获取分类下的数据
		adminData.POST("/data", adminDataController.CreateData)                                 // 新增数据
		adminData.PUT("/data/:region_id/:category_id/:year", adminDataController.UpdateData)    // 修改数据
		adminData.DELETE("/data/:region_id/:category_id/:year", adminDataController.DeleteData) // 删除数据
	}

}
