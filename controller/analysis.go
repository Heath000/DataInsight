package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type AnalysisController struct{}

// AnalysisOverview 总览功能接口，计算描述性统计信息
// 输入的数据格式：
//
//	{
//	   "data": [120.5, 132.1, 149.7, 165.3, 178.8, 190.2]
//	}
func (ctrl *AnalysisController) AnalysisOverview(c *gin.Context) {
	// 接收输入数据
	var inputData map[string]interface{}
	if err := c.ShouldBindJSON(&inputData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "error": err.Error()})
		return
	}

	// 将输入数据转换为 JSON 字符串
	jsonData, err := json.Marshal(inputData)
	if err != nil {
		log.Println("Error marshaling JSON:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error processing data"})
		return
	}

	// 调用 Python 脚本进行描述性统计计算
	currentDir, err := os.Getwd()
	if err != nil {
		log.Println("Error getting current directory:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error processing data"})
		return
	}
	scriptPath := filepath.Join(currentDir, "controller", "AnalysisProcessPyScripts", "overall.py")

	cmd := exec.Command("python", scriptPath, string(jsonData))
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error executing Python script: %s, output: %s", err, string(output))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error processing data"})
		return
	}

	// 解析 Python 脚本返回的数据
	var stats map[string]interface{}
	if err := json.Unmarshal(output, &stats); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error parsing Python response"})
		return
	}

	// 返回统计结果
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    stats,
	})
}

// AnalysisLinearRegress 线性回归预测接口
// 输入的数据格式：
//
//	{
//	   "data": [[1.0, 2.0], [2.0, 3.0], [3.0, 4.0], [4.0, 5.0]],  // 自变量数据 (二维数组，每行表示一个样本，每列表示一个自变量)
//	   "labels": [3.0, 5.0, 7.0, 9.0],                           // 目标值 (一维数组，与 data 对应)
//	   "predict_data": [[5.0, 6.0], [6.0, 7.0]]                  // 待预测数据 (二维数组，与 data 特征数一致)
//	}
func (ctrl *AnalysisController) AnalysisLinearRegress(c *gin.Context) {
	// 接收输入数据
	var inputData map[string]interface{}
	if err := c.ShouldBindJSON(&inputData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "error": err.Error()})
		return
	}

	// 检查必要字段
	if _, ok := inputData["data"]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing data in input"})
		return
	}
	if _, ok := inputData["labels"]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing labels in input"})
		return
	}
	if _, ok := inputData["predict_data"]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing predict_data in input"})
		return
	}

	// 将输入数据转换为 JSON 字符串
	jsonData, err := json.Marshal(inputData)
	if err != nil {
		log.Println("Error marshaling JSON:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error processing data"})
		return
	}

	// 调用 Python 脚本进行线性回归预测
	currentDir, err := os.Getwd()
	if err != nil {
		log.Println("Error getting current directory:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error processing data"})
		return
	}
	scriptPath := filepath.Join(currentDir, "controller", "AnalysisProcessPyScripts", "linear_regression.py")

	cmd := exec.Command("python", scriptPath, string(jsonData))
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error executing Python script: %s, output: %s", err, string(output))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error processing data"})
		return
	}

	// 解析 Python 脚本返回的数据
	var response map[string]interface{}
	if err := json.Unmarshal(output, &response); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error parsing Python response"})
		return
	}

	// 返回预测结果或错误信息
	if errMsg, exists := response["error"]; exists {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error in Linear Regression script", "error": errMsg})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    response["predictions"],
	})
}

// AnalysisARIMA ARIMA预测接口
// 输入的数据格式：
//
//	{
//	   "time_series": [2000, 2001, 2002, 2003, 2004, 2005],  // 时间序列信息（年份）
//	   "data": [123.0, 150.5, 180.7, 210.2, 216.5, 241.9],   // 时间序列数据
//	   "p": 1,                                               // 自回归阶数
//	   "d": 1,                                               // 差分次数
//	   "q": 1,                                               // 移动平均阶数
//	   "forecast_steps": 5                                   // 预测步数
//	}
func (ctrl *AnalysisController) AnalysisARIMA(c *gin.Context) {
	// 接收输入数据
	var inputData map[string]interface{}
	if err := c.ShouldBindJSON(&inputData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "error": err.Error()})
		return
	}

	// 检查是否包含时间序列信息
	if _, ok := inputData["time_series"]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing time_series in input"})
		return
	}
	if _, ok := inputData["data"]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing data in input"})
		return
	}

	// 将输入数据转换为 JSON 字符串
	jsonData, err := json.Marshal(inputData)
	if err != nil {
		log.Println("Error marshaling JSON:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error processing data"})
		return
	}

	// 调用 Python 脚本进行 ARIMA 预测
	currentDir, err := os.Getwd()
	if err != nil {
		log.Println("Error getting current directory:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error processing data"})
		return
	}
	scriptPath := filepath.Join(currentDir, "controller", "AnalysisProcessPyScripts", "ARIMA.py")

	cmd := exec.Command("python", scriptPath, string(jsonData))
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error executing Python script: %s, output: %s", err, string(output))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error processing data"})
		return
	}

	// 解析 Python 脚本返回的数据
	var response map[string]interface{}
	if err := json.Unmarshal(output, &response); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error parsing Python response"})
		return
	}

	// 返回预测结果或错误信息
	if errMsg, exists := response["error"]; exists {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error in ARIMA script", "error": errMsg})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":            "success",
		"data":               response["predictions"],
		"future_time_series": response["future_time_series"],
	})
}

// AnalysisGreyPredict 灰色预测接口
// 输入的数据格式：
//
//	{
//	   "data": [120.5, 132.1, 149.7, 165.3, 178.8, 190.2],   // 时间序列数据
//	   "time_series": [2000, 2001, 2002, 2003, 2004, 2005], // 时间序列信息（年份）
//	   "forecast_steps": 5                                  // 预测步数
//	}
func (ctrl *AnalysisController) AnalysisGreyPredict(c *gin.Context) {
	// 接收输入数据
	var inputData map[string]interface{}
	if err := c.ShouldBindJSON(&inputData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "error": err.Error()})
		return
	}

	// 检查是否包含必要的字段
	if _, ok := inputData["time_series"]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing time_series in input"})
		return
	}
	if _, ok := inputData["data"]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing data in input"})
		return
	}

	// 将输入数据转换为 JSON 字符串
	jsonData, err := json.Marshal(inputData)
	if err != nil {
		log.Println("Error marshaling JSON:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error processing data"})
		return
	}

	// 调用 Python 脚本进行灰色预测
	currentDir, err := os.Getwd()
	if err != nil {
		log.Println("Error getting current directory:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error processing data"})
		return
	}
	scriptPath := filepath.Join(currentDir, "controller", "AnalysisProcessPyScripts", "grey_predict.py")

	cmd := exec.Command("python", scriptPath, string(jsonData))
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error executing Python script: %s, output: %s", err, string(output))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error processing data"})
		return
	}

	// 解析 Python 脚本返回的数据
	var response map[string]interface{}
	if err := json.Unmarshal(output, &response); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error parsing Python response"})
		return
	}

	// 返回预测结果或错误信息
	if errMsg, exists := response["error"]; exists {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error in Grey Model script", "error": errMsg})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":            "success",
		"data":               response["predictions"],
		"future_time_series": response["future_time_series"],
	})
}

// AnalysisBP BP神经网络回归预测接口
// 输入的数据格式：
//
//	{
//	   "data": [[1.0, 2.0], [2.0, 3.0], [3.0, 4.0], [4.0, 5.0]],  // 自变量数据 (二维数组，每行表示一个样本，每列表示一个自变量)
//	   "labels": [3.0, 5.0, 7.0, 9.0],                           // 目标值 (一维数组，与 data 对应)
//	   "predict_data": [[5.0, 6.0], [6.0, 7.0]],                 // 待预测数据 (二维数组，与 data 特征数一致)
//	   "hidden_layers": 100,                                     // 隐藏层节点数
//	   "max_iter": 1000,                                         // 最大迭代次数
//	   "learning_rate_init": 0.01                                // 学习率
//	}
func (ctrl *AnalysisController) AnalysisBP(c *gin.Context) {
	// 接收输入数据
	var inputData map[string]interface{}
	if err := c.ShouldBindJSON(&inputData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "error": err.Error()})
		return
	}

	// 检查必要字段
	if _, ok := inputData["data"]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing data in input"})
		return
	}
	if _, ok := inputData["labels"]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing labels in input"})
		return
	}
	if _, ok := inputData["predict_data"]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing predict_data in input"})
		return
	}

	// 将输入数据转换为 JSON 字符串
	jsonData, err := json.Marshal(inputData)
	if err != nil {
		log.Println("Error marshaling JSON:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error processing data"})
		return
	}

	// 调用 Python 脚本进行 BP 神经网络预测
	currentDir, err := os.Getwd()
	if err != nil {
		log.Println("Error getting current directory:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error processing data"})
		return
	}
	scriptPath := filepath.Join(currentDir, "controller", "AnalysisProcessPyScripts", "BP.py")

	cmd := exec.Command("python", scriptPath, string(jsonData))
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error executing Python script: %s, output: %s", err, string(output))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error processing data"})
		return
	}

	// 解析 Python 脚本返回的数据
	var response map[string]interface{}
	if err := json.Unmarshal(output, &response); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error parsing Python response"})
		return
	}

	// 返回预测结果或错误信息
	if errMsg, exists := response["error"]; exists {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error in BP script", "error": errMsg})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    response["predictions"],
	})
}

// AnalysisSVM SVM分类接口
// 输入的数据格式：
//
//	{
//	   "data": [[1.0, 2.0], [2.0, 3.0], [3.0, 4.0], [4.0, 5.0]],  // 自变量数据 (二维数组，每行表示一个样本，每列表示一个自变量)
//	   "labels": ["A", "B", "A", "B"],                          // 类别名称 (与 data 对应的目标值)
//	   "predict_data": [[5.0, 6.0], [6.0, 7.0]],                 // 待分类数据 (二维数组，与 data 特征数一致)
//	   "C": 1.0,                                                // 惩罚系数
//	   "tol": 0.001,                                            // 误差收敛条件
//	   "max_iter": 1000                                         // 最大迭代次数
//	}
func (ctrl *AnalysisController) AnalysisSVM(c *gin.Context) {
	// 接收输入数据
	var inputData map[string]interface{}
	if err := c.ShouldBindJSON(&inputData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "error": err.Error()})
		return
	}

	// 检查必要字段
	if _, ok := inputData["data"]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing data in input"})
		return
	}
	if _, ok := inputData["labels"]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing labels in input"})
		return
	}
	if _, ok := inputData["predict_data"]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing predict_data in input"})
		return
	}

	// 将输入数据转换为 JSON 字符串
	jsonData, err := json.Marshal(inputData)
	if err != nil {
		log.Println("Error marshaling JSON:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error processing data"})
		return
	}

	// 调用 Python 脚本进行 SVM 分类
	currentDir, err := os.Getwd()
	if err != nil {
		log.Println("Error getting current directory:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error processing data"})
		return
	}
	scriptPath := filepath.Join(currentDir, "controller", "AnalysisProcessPyScripts", "SVM.py")

	cmd := exec.Command("python", scriptPath, string(jsonData))
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error executing Python script: %s, output: %s", err, string(output))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error processing data"})
		return
	}

	// 解析 Python 脚本返回的数据
	var response map[string]interface{}
	if err := json.Unmarshal(output, &response); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error parsing Python response"})
		return
	}

	// 返回预测结果或错误信息
	if errMsg, exists := response["error"]; exists {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error in SVM script", "error": errMsg})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "success",
		"data":          response["predictions"],
		"probabilities": response["probabilities"],
	})
}

// AnalysisRandomForest 随机森林分类接口
// 输入的数据格式：
//
//	{
//	   "data": [[1.0, 2.0], [2.0, 3.0], [3.0, 4.0], [4.0, 5.0]],  // 自变量数据 (二维数组，每行表示一个样本，每列表示一个自变量)
//	   "labels": ["A", "B", "A", "B"],                          // 类别名称 (与 data 对应的目标值)
//	   "predict_data": [[5.0, 6.0], [6.0, 7.0]],                 // 待分类数据 (二维数组，与 data 特征数一致)
//	   "n_estimators": 100,                                     // 决策树数量
//	   "max_depth": 10,                                         // 树的最大深度
//	   "max_leaf_nodes": 20,                                    // 叶节点的最大数量
//	   "min_samples_split": 2,                                  // 内部节点分裂的最小样本数
//	   "min_samples_leaf": 1                                    // 叶节点的最小样本数
//	}
func (ctrl *AnalysisController) AnalysisRandomForest(c *gin.Context) {
	// 接收输入数据
	var inputData map[string]interface{}
	if err := c.ShouldBindJSON(&inputData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "error": err.Error()})
		return
	}

	// 检查必要字段
	if _, ok := inputData["data"]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing data in input"})
		return
	}
	if _, ok := inputData["labels"]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing labels in input"})
		return
	}
	if _, ok := inputData["predict_data"]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing predict_data in input"})
		return
	}

	// 将输入数据转换为 JSON 字符串
	jsonData, err := json.Marshal(inputData)
	if err != nil {
		log.Println("Error marshaling JSON:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error processing data"})
		return
	}

	// 调用 Python 脚本进行随机森林分类
	currentDir, err := os.Getwd()
	if err != nil {
		log.Println("Error getting current directory:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error processing data"})
		return
	}
	scriptPath := filepath.Join(currentDir, "controller", "AnalysisProcessPyScripts", "RandomForest.py")

	cmd := exec.Command("python", scriptPath, string(jsonData))
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error executing Python script: %s, output: %s", err, string(output))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error processing data"})
		return
	}

	// 解析 Python 脚本返回的数据
	var response map[string]interface{}
	if err := json.Unmarshal(output, &response); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error parsing Python response"})
		return
	}

	// 返回预测结果或错误信息
	if errMsg, exists := response["error"]; exists {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error in Random Forest script", "error": errMsg})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "success",
		"data":          response["predictions"],
		"probabilities": response["probabilities"],
	})
}
