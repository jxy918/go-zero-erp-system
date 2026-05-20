package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type InventoryAdjustRequestModel struct {
	DB *gorm.DB
}

func NewInventoryAdjustRequestModel(db *gorm.DB) *InventoryAdjustRequestModel {
	return &InventoryAdjustRequestModel{DB: db}
}

func (m *InventoryAdjustRequestModel) TableName() string {
	return "inventory_adjust_requests"
}

func (m *InventoryAdjustRequestModel) Create(req *InventoryAdjustRequest) error {
	req.RequestNo = m.generateRequestNo()
	return m.DB.Create(req).Error
}

func (m *InventoryAdjustRequestModel) Update(req *InventoryAdjustRequest) error {
	return m.DB.Save(req).Error
}

func (m *InventoryAdjustRequestModel) Delete(id uint) error {
	return m.DB.Unscoped().Delete(&InventoryAdjustRequest{}, id).Error
}

func (m *InventoryAdjustRequestModel) GetByID(id uint) (*InventoryAdjustRequest, error) {
	var req InventoryAdjustRequest
	err := m.DB.Preload("Product").Preload("Warehouse").Preload("Applicant").Preload("Approver").
		First(&req, id).Error
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func (m *InventoryAdjustRequestModel) List(page, pageSize int, status int) ([]InventoryAdjustRequest, int64, error) {
	var requests []InventoryAdjustRequest
	var total int64

	query := m.DB.Model(&InventoryAdjustRequest{}).Preload("Product").Preload("Warehouse").Preload("Applicant").Preload("Approver")

	if status > 0 {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&requests).Error
	if err != nil {
		return nil, 0, err
	}

	return requests, total, nil
}

func (m *InventoryAdjustRequestModel) Approve(id uint, approverID uint, note string) error {
	now := time.Now()
	return m.DB.Model(&InventoryAdjustRequest{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":       2,
		"approver_id":  approverID,
		"approve_time": &now,
		"approve_note": note,
	}).Error
}

func (m *InventoryAdjustRequestModel) Reject(id uint, approverID uint, note string) error {
	return m.DB.Model(&InventoryAdjustRequest{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":       3,
		"approver_id":  approverID,
		"approve_note": note,
	}).Error
}

func (m *InventoryAdjustRequestModel) generateRequestNo() string {
	now := time.Now()
	var count int64
	m.DB.Model(&InventoryAdjustRequest{}).Where("DATE(created_at) = ?", now.Format("2006-01-02")).Count(&count)
	return fmt.Sprintf("ADJ-%s-%03d", now.Format("20060102"), count+1)
}
