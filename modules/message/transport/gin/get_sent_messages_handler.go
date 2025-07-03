package gin

import (
	"messenging_test/modules/message/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetSentMessagesHandler godoc
// @Summary Get list of sent messages
// @Description Returns the list of successfully sent messages
// @Tags messages
// @Produce json
// @Success 200 {array} model.Message
// @Failure 500 {object} map[string]string
// @Router /messages/sent [get]
func GetSentMessagesHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var messages []model.Message
		if err := db.Where("sent = ?", true).Order("sent_at desc").Find(&messages).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, messages)
	}
}
