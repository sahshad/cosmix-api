package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`
}

func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.Must(uuid.NewV7())
	}
	return nil
}