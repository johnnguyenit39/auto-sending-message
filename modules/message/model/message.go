package model

import (
	"log"
	"messenging_test/common"
	"time"

	"gorm.io/gorm"
)

type Message struct {
	common.BaseModel
	To      string    `gorm:"type:varchar(20);not null" json:"to"`
	Content string    `gorm:"type:varchar(160);not null" json:"content"`
	Sent    bool      `gorm:"default:false" json:"sent"`
	SentAt  time.Time `gorm:"default:null" json:"sent_at"`
}

func (*Message) TableName() string {
	return "messages"
}

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&Message{})
	if err != nil {
		log.Println(err.Error())
	}
	return err
}
