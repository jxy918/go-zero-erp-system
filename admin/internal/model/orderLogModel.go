package model

import (
	"time"

	"gorm.io/gorm"
)

type OrderLog struct {
	ID           uint      `gorm:"primarykey"`
	OrderID      uint      `gorm:"column:order_id;not null"`
	OrderType    int       `gorm:"column:order_type;not null;default:1"`
	BeforeStatus int       `gorm:"column:before_status;not null"`
	AfterStatus  int       `gorm:"column:after_status;not null"`
	OperatorID   uint      `gorm:"column:operator_id;not null"`
	OperatorName string    `gorm:"column:operator_name;not null;size:100"`
	Remark       string    `gorm:"column:remark;size:500"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (OrderLog) TableName() string {
	return "order_logs"
}

type OrderLogModel struct {
	db *gorm.DB
}

func NewOrderLogModel(db *gorm.DB) *OrderLogModel {
	return &OrderLogModel{db: db}
}

func (m *OrderLogModel) Insert(data *OrderLog) error {
	return m.db.Create(data).Error
}

func (m *OrderLogModel) List(orderID, orderType, operatorID uint, startTime, endTime string) ([]*OrderLog, error) {
	var logs []*OrderLog
	query := m.db.Model(&OrderLog{})

	if orderID > 0 {
		query = query.Where("order_id = ?", orderID)
	}
	if orderType > 0 {
		query = query.Where("order_type = ?", orderType)
	}
	if operatorID > 0 {
		query = query.Where("operator_id = ?", operatorID)
	}
	if startTime != "" {
		query = query.Where("created_at >= ?", startTime)
	}
	if endTime != "" {
		query = query.Where("created_at <= ?", endTime)
	}

	err := query.Order("created_at DESC").Find(&logs).Error
	return logs, err
}

func (m *OrderLogModel) GetByID(id uint) (*OrderLog, error) {
	var log OrderLog
	err := m.db.Where("id = ?", id).First(&log).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}
