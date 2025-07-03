package biz

import (
	"context"
	"errors"
	messageModel "messenging_test/modules/message/model"

	"gorm.io/gorm"
)

func GetUnsentMessages(ctx context.Context, db *gorm.DB, limit int) ([]messageModel.Message, error) {
	var messages []messageModel.Message
	if err := db.WithContext(ctx).Where("sent = ?", false).Order("created_at").Limit(limit).Find(&messages).Error; err != nil {
		return nil, err
	}
	for _, m := range messages {
		if len(m.Content) > 160 {
			return nil, errors.New("message content exceeds 160 characters")
		}
	}
	return messages, nil
}
