package gin

import (
	"net/http"
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

var autoSenderRunning int32

// StartAutoSenderHandler godoc
// @Summary Start automatic message sending
// @Description Start the background process to send messages every 2 minutes
// @Tags auto-sender
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /auto-sender/start [post]
func StartAutoSenderHandler(startFunc func()) gin.HandlerFunc {
	return func(c *gin.Context) {
		if atomic.LoadInt32(&autoSenderRunning) == 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Auto sender already running"})
			return
		}
		atomic.StoreInt32(&autoSenderRunning, 1)
		go startFunc()
		c.JSON(http.StatusOK, gin.H{"message": "Auto sender started"})
	}
}

// StopAutoSenderHandler godoc
// @Summary Stop automatic message sending
// @Description Stop the background process for automatic message sending
// @Tags auto-sender
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /auto-sender/stop [post]
func StopAutoSenderHandler(stopFunc func()) gin.HandlerFunc {
	return func(c *gin.Context) {
		if atomic.LoadInt32(&autoSenderRunning) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Auto sender not running"})
			return
		}
		atomic.StoreInt32(&autoSenderRunning, 0)
		stopFunc()
		c.JSON(http.StatusOK, gin.H{"message": "Auto sender stopped"})
	}
}
