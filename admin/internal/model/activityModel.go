package model

import (
	"time"

	"gorm.io/gorm"
)

// Activity 活动日志模型
type Activity struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	Username  string    `gorm:"size:50;index" json:"username"`
	Action    string    `gorm:"size:255" json:"action"`
	URL       string    `gorm:"size:255" json:"url"`
	IP        string    `gorm:"size:50" json:"ip"`
	CreatedAt time.Time `json:"created_at"`
}

// ActivityModel 活动日志模型接口
type ActivityModel interface {
	Create(activity *Activity) error
	GetRecent(limit int) ([]Activity, int64, error)
	ListByUsername(username string, page, pageSize int) ([]Activity, int64, error)
	List(page, pageSize int) ([]Activity, int64, error)
	ListByUserID(userID uint, page, pageSize int) ([]Activity, int64, error)
}

// activityModel 活动日志模型实现
type activityModel struct {
	db *gorm.DB
}

// NewActivityModel 创建活动日志模型实例
func NewActivityModel(db *gorm.DB) ActivityModel {
	return &activityModel{db: db}
}

// Create 创建活动日志
func (m *activityModel) Create(activity *Activity) error {
	return m.db.Create(activity).Error
}

// GetRecent 获取最近的活动日志
func (m *activityModel) GetRecent(limit int) ([]Activity, int64, error) {
	var activities []Activity
	err := m.db.Order("created_at DESC").Limit(limit).Find(&activities).Error
	return activities, int64(len(activities)), err
}

// ListByUsername 根据用户名搜索活动日志
func (m *activityModel) ListByUsername(username string, page, pageSize int) ([]Activity, int64, error) {
	var activities []Activity
	var total int64

	offset := (page - 1) * pageSize

	// 构建查询
	query := m.db
	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}

	// 获取总数
	if err := query.Model(&Activity{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&activities).Error; err != nil {
		return nil, 0, err
	}

	return activities, total, nil
}

// List 获取活动日志列表（分页）
func (m *activityModel) List(page, pageSize int) ([]Activity, int64, error) {
	var activities []Activity
	var total int64

	offset := (page - 1) * pageSize

	// 获取总数
	if err := m.db.Model(&Activity{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	if err := m.db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&activities).Error; err != nil {
		return nil, 0, err
	}

	return activities, total, nil
}

// ListByUserID 根据用户ID获取活动日志列表
func (m *activityModel) ListByUserID(userID uint, page, pageSize int) ([]Activity, int64, error) {
	var activities []Activity
	var total int64

	offset := (page - 1) * pageSize

	// 获取总数
	if err := m.db.Model(&Activity{}).
		Where("user_id = ?", userID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	if err := m.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&activities).Error; err != nil {
		return nil, 0, err
	}

	return activities, total, nil
}
