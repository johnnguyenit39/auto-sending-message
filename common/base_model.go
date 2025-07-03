package common

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"-"`
}

// BeforeCreate is a GORM hook that runs before a new record is created
func (baseModel *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	if baseModel.ID == uuid.Nil {
		baseModel.ID = uuid.New()
	}
	if baseModel.CreatedAt.IsZero() {
		baseModel.CreatedAt = time.Now().UTC()
	}
	baseModel.UpdatedAt = time.Now().UTC()
	return nil
}

func (baseModel *BaseModel) BeforeUpdate(tx *gorm.DB) (err error) {
	baseModel.UpdatedAt = time.Now().UTC()
	return
}
