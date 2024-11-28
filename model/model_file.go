package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// File represents the file model
type File struct {
	FileID     uint      `gorm:"primaryKey;autoIncrement;column:file_id" json:"fileId"`          // 文件 ID，主键且自动递增
	UserID     uint      `gorm:"not null;column:user_id" json:"userId"`                          // 用户 ID，外键，关联 users 表
	Filename   string    `gorm:"not null;column:filename" json:"filename"`                       // 文件名，最大长度 255
	UploadTime time.Time `gorm:"default:current_timestamp;column:upload_time" json:"uploadTime"` // 上传时间，默认为当前时间戳

	// Relationships
	User User `gorm:"foreignKey:UserID;references:ID" json:"user"` // 定义外键关联，指向 User 表的 ID 字段
}

// 虽然gorm显式映射到小写复数，但是这里显式声明一下
func (File) TableName() string {
	return "files"
}

func GetUserIDByEmail(email string) (uint, error) {
	var user User
	if err := DB().Where("email = ?", email).First(&user).Error; err != nil {
		return 0, errors.New("user not found")
	}
	return user.ID, nil
}

// GetFileListByUserId 获取指定用户 ID 的所有文件
func (f *File) GetFileListByUserId(userID uint) ([]File, error) {
	var files []File
	err := DB().Where("user_id = ?", userID).Find(&files).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrDataNotFound
	}
	return files, err
}

// GetFileByIDAndUserID 根据文件 ID 和用户 ID 获取文件
func (f *File) GetFileByIDAndUserID(fileID, userID uint) (*File, error) {
	var file File
	err := DB().Where("file_id = ? AND user_id = ?", fileID, userID).First(&file).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrDataNotFound
	}
	return &file, err
}

// DeleteFileByIDAndUserID 根据文件 ID 和用户 ID 删除文件
func (f *File) DeleteFileByIDAndUserID(fileID, userID uint) error {
	err := DB().Where("file_id = ? AND user_id = ?", fileID, userID).Delete(&File{}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrDataNotFound
	}
	return err
}

// PutFileInfo 根据文件 ID 和用户 ID 更新文件信息
func (f *File) PutFileInfo(fileID, userID uint, filename string) error {
	err := DB().Model(&File{}).Where("file_id = ? AND user_id = ?", fileID, userID).Updates(File{Filename: filename}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrDataNotFound
	}
	return err
}

// PostFileInfo 创建新的文件记录
func (f *File) PostFileInfo(userID uint, filename string) error {
	newFile := File{
		UserID:   userID,
		Filename: filename,
	}
	err := DB().Create(&newFile).Error
	if err != nil {
		return err
	}
	return nil
}
// GetLastFileID 获取最新的 fileID，返回最后插入的文件的 fileID
func GetLastFileID() (uint, error) {
	var lastFile File
	// 使用 Order 降序排列，Limit 1 只取一条最新的记录
	err := DB().Order("file_id desc").Limit(1).First(&lastFile).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果没有找到记录，说明文件表为空，返回 0 和错误
			return 0, errors.New("no files found")
		}
		// 如果查询出错，返回错误
		return 0, err
	}
	// 返回最新记录的 fileID
	return lastFile.FileID, nil
}