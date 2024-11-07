package models

import "time"

type BaseModel struct {
	CreatedBy  string    `json:"created_by"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedBy string    `json:"modified_by"`
	UpdatedAt  time.Time `json:"updated_at"`
	IsDeleted  bool      `gorm:"default:false" json:"is_deleted"`
}
