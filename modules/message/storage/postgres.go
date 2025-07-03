package storage

import (
	"gorm.io/gorm"
)

type MessageStore struct {
	DB *gorm.DB
}

func NewMessageStore(db *gorm.DB) *MessageStore {
	return &MessageStore{DB: db}
}
