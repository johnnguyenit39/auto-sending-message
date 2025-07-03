package biz

import (
	"context"
	messageModel "messenging_test/modules/message/model"
	"time"

	"gorm.io/gorm"
)

func MarkMessageSent(ctx context.Context, db *gorm.DB, id string) error {
	return db.WithContext(ctx).Model(&messageModel.Message{}).Where("id = ?", id).Updates(map[string]interface{}{
		"sent":    true,
		"sent_at": time.Now(),
	}).Error
}
