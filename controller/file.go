package controller

import (
	"fmt"
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

// GetFile 获取单个文件
func (f *FileController) GetFile(c *gin.Context) {
	// 从请求中获取文件ID
	fileIDStr := c.Param("file_id")
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

	// 获取当前工作目录并拼接脚本路径
	currentDir, err := os.Getwd()
	if err != nil {
		log.Println("Error getting current directory:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error processing data"})
		return
	}

	// 使用 fileID + filename + extension 构建文件路径
	filePath := filepath.Join(currentDir, "file", strconv.FormatUint(uint64(fileID), 10)+"-"+file.Filename)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "File not found on server",
		})
		return
	}

	// 返回文件内容，返回时只提供 filename + extension
	c.File(filePath)
}

// DeleteFile 删除文件
func (f *FileController) DeleteFile(c *gin.Context) {
	// 从请求中获取文件ID
	fileIDStr := c.Param("file_id")
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
	file, err := fileModel.GetFileByIDAndUserID(uint(fileID), userID)
	if err != nil {
		if err == model.ErrDataNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "File not found or access denied",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "Failed to retrieve file information from database",
			})
		}
		return
	}

	// 获取当前工作目录并拼接脚本路径
	currentDir, err := os.Getwd()
	if err != nil {
		log.Println("Error getting current directory:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error processing data"})
		return
	}

	// 使用 fileID + filename + extension 构建文件路径
	filePath := filepath.Join(currentDir, "file", strconv.FormatUint(uint64(fileID), 10)+"-"+file.Filename)

	// 检查文件是否存在
	_, fileErr := os.Stat(filePath)

	// 如果文件存在，则删除文件
	if fileErr == nil {
		// 删除服务器中的文件
		if err := os.Remove(filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "Failed to delete file from server",
			})
			log.Println("File deletion error:", err)
			return
		}
	} else if os.IsNotExist(fileErr) {
		// 如果文件不存在，检查数据库记录并删除
		log.Println("File not found on server, checking database record...")

		// 删除数据库中的记录
		err = fileModel.DeleteFileByIDAndUserID(uint(fileID), userID)
		if err != nil {
			if err == model.ErrDataNotFound {
				c.JSON(http.StatusNotFound, gin.H{
					"code":    404,
					"message": "File not found in database",
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "Failed to delete file record from database",
				})
			}
			log.Println("Database record deletion error:", err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "File record deleted from database, but no physical file found",
		})
		return
	}

	// 如果文件存在且成功删除，删除数据库中的记录
	err = fileModel.DeleteFileByIDAndUserID(uint(fileID), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to delete file record from database",
		})
		log.Println("Database record deletion error:", err)
		return
	}

	// 成功删除文件和数据库记录
	c.JSON(http.StatusOK, gin.H{
		"message": "File deleted successfully, record removed from database",
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

	// 在数据库中保存文件记录
	fileInfo := model.File{
		UserID:     userID,
		Filename:   file.Filename,
		UploadTime: time.Now(),
	}

	// 调用 PostFileInfo 插入数据库
	if err := fileInfo.PostFileInfo(userID, file.Filename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to save file information in database",
		})
		return
	}

	// 使用 GetLastFileID 获取最新的 fileID
	lastFileID, err := model.GetLastFileID()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to retrieve last file ID",
		})
		return
	}

	// 获取当前工作目录并拼接脚本路径
	currentDir, err := os.Getwd()
	if err != nil {
		log.Println("Error getting current directory:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error processing data"})
		return
	}

	// 定义文件存储路径
	uploadPath := filepath.Join(currentDir, "file")
	// 确保文件夹存在，不存在时创建它
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		log.Println("Error creating upload directory:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating directory"})
		return
	}

	// 获取文件扩展名
	fileExtension := filepath.Ext(file.Filename) // 获取文件的扩展名（如 .jpg、.png）
	if fileExtension == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "File does not have an extension",
		})
		return
	}

	// 使用 fileID + filename + extension 生成文件路径
	fullPath := filepath.Join(uploadPath, strconv.FormatUint(uint64(lastFileID), 10)+"-"+file.Filename)

	// 打印生成的文件路径
	fmt.Println("Saving file to path:", fullPath)

	// 保存文件到指定路径
	if err := c.SaveUploadedFile(file, fullPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to save file to server",
		})
		return
	}

	// 成功响应
	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully",
	})
}
