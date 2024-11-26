package controller

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Heath000/fzuSE2024/model"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// FileController 处理文件相关的操作
type FileController struct{}

// GetFileList 获取用户的文件列表
func (f *FileController) GetFileList(c *gin.Context) {
	// 从 JWT 中提取 claims 并获取用户 ID
	claims := jwt.ExtractClaims(c)
	idValue, ok := claims["ID"].(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "User ID not found in JWT",
		})
		return
	}

	// 将提取的 ID 转换为 uint 类型
	userID := uint(idValue)

	// 使用提取的 userID 获取用户的文件列表
	var fileModel model.File
	files, err := fileModel.GetFileListByUserId(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to retrieve files",
		})
		return
	}

	// 返回文件列表
	c.JSON(http.StatusOK, gin.H{
		"files": files,
	})
}

// GetFile 获取单个文件信息
func (f *FileController) GetFile(c *gin.Context) {
	// 从请求中获取文件ID
	fileIDStr := c.Param("file_id")
	//filename := c.Param("filename")
	fileID, err := strconv.ParseUint(fileIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid file ID format",
		})
		return
	}

	// 从 JWT 中提取 claims 并获取用户 ID
	claims := jwt.ExtractClaims(c)
	idValue, ok := claims["ID"].(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "User ID not found in JWT",
		})
		return
	}
	userID := uint(idValue)

	// 根据 file_id 和 userID 获取该文件的信息
	var fileModel model.File
	file, err := fileModel.GetFileByIDAndUserID(uint(fileID), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "File not found",
		})
		return
	}

	// 真正地获取服务器里的文件，然后发回****************
	// 获取当前工作目录并拼接脚本路径
	currentDir, err := os.Getwd()
	if err != nil {
		log.Println("Error getting current directory:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error processing data"})
		return
	}
	filePath := filepath.Join(currentDir, "file", fileID)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "File not found on server",
		})
		return
	}

	// 返回文件内容
	c.File(filePath)
	// code ends here
}

// DeleteFile 删除文件
func (f *FileController) DeleteFile(c *gin.Context) {
	// 从请求中获取文件ID
	fileIDStr := c.Param("file_id")
	//filename := c.Param("filename")
	fileID, err := strconv.ParseUint(fileIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid file ID format",
		})
		return
	}

	// 从 JWT 中提取 claims 并获取用户 ID
	claims := jwt.ExtractClaims(c)
	idValue, ok := claims["ID"].(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "User ID not found in JWT",
		})
		return
	}
	userID := uint(idValue)

	// 删除文件，使用 fileID 和 userID 双重验证
	var fileModel model.File
	err = fileModel.DeleteFileByIDAndUserID(uint(fileID), userID)
	if err != nil {
		if err == model.ErrDataNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "File not found or access denied",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "Failed to delete file",
			})
		}
		return
	}
	currentDir, err := os.Getwd()
	if err != nil {
		log.Println("Error getting current directory:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error processing data"})
		return
	}
	filePath := filepath.Join(currentDir, "file", fileID)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "File does not exist on server",
		})
		return
	}

	// 删除服务器中的文件
	if err := os.Remove(filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to delete file from server",
		})
		log.Println("File deletion error:", err)
		return
	}

	// 成功删除
	c.JSON(http.StatusOK, gin.H{
		"message": "File deleted successfully",
	})
}

// UploadFile 上传文件
func (f *FileController) UploadFile(c *gin.Context) {
	// 从 JWT 中提取 claims 并获取用户 ID
	claims := jwt.ExtractClaims(c)
	idValue, ok := claims["ID"].(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "User ID not found in JWT",
		})
		return
	}
	userID := uint(idValue)

	// 获取上传的文件
	file, _ := c.FormFile("file")
	if file == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "No file uploaded",
		})
		return
	}

	// 真正地将文件保存到服务器的指定路径********************** 代码写在这里
	// 在数据库中保存文件记录
	fileInfo := model.File{
		UserID:     userID,
		Filename:   file.Filename,
		UploadTime: time.Now(), // 假设 UploadTime 自动设置
	}
	// 定义文件存储路径
	// 获取当前工作目录并拼接脚本路径
	currentDir, err := os.Getwd()
	if err != nil {
		log.Println("Error getting current directory:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error processing data"})
		return
	}
	uploadPath := filepath.Join(currentDir, "file")
	fullPath := uploadPath + file.FileID

	// 保存文件到指定路径
	if err := c.SaveUploadedFile(file, fullPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to save file to server",
		})
		return
	}
	//code ends here

	if err := fileInfo.PostFileInfo(userID, file.Filename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to save file information in database",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully",
	})
}
