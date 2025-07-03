package auto_sender

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"messenging_test/config/pubsub"
	httprequest "messenging_test/http_request"
	"messenging_test/modules/message/biz"
	messageModel "messenging_test/modules/message/model"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type MessageSentEvent struct {
	MessageID string    `json:"messageId"`
	ID        string    `json:"id"`
	SentAt    time.Time `json:"sent_at"`
}

type AutoSender struct {
	db         *gorm.DB
	redis      *redis.Client
	webhookURL string
	webhookKey string
	interval   time.Duration
	ctx        context.Context
	cancel     context.CancelFunc
	mu         sync.Mutex
	running    bool
	pubsub     *pubsub.PubSub
	httpSvc    *httprequest.HttpService
}

func NewAutoSender(db *gorm.DB, redis *redis.Client, webhookURL, webhookKey string, interval time.Duration, pubsub *pubsub.PubSub, httpSvc *httprequest.HttpService) *AutoSender {
	ctx, cancel := context.WithCancel(context.Background())
	return &AutoSender{
		db:         db,
		redis:      redis,
		webhookURL: webhookURL,
		webhookKey: webhookKey,
		interval:   interval,
		ctx:        ctx,
		cancel:     cancel,
		pubsub:     pubsub,
		httpSvc:    httpSvc,
	}
}

func (a *AutoSender) Start() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.running {
		return errors.New("auto sender already running")
	}
	a.ctx, a.cancel = context.WithCancel(context.Background())
	a.running = true
	go a.run()
	return nil
}

func (a *AutoSender) Stop() {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.running {
		a.cancel()
		a.running = false
	}
}

func (a *AutoSender) run() {
	ticker := time.NewTicker(a.interval)
	defer ticker.Stop()
	for {
		select {
		case <-a.ctx.Done():
			return
		case <-ticker.C:
			a.sendMessages()
		}
	}
}

func (a *AutoSender) sendMessages() {
	ctx := context.Background()
	msgs, err := biz.GetUnsentMessages(ctx, a.db, 2)
	if err != nil || len(msgs) == 0 {
		return
	}
	for _, msg := range msgs {
		messageId, sentAt, err := a.sendToWebhook(msg)
		if err == nil {
			_ = biz.MarkMessageSent(ctx, a.db, msg.ID.String())
			if a.pubsub != nil && messageId != "" {
				event := MessageSentEvent{MessageID: messageId, ID: msg.ID.String(), SentAt: sentAt}
				a.pubsub.Publish("message.sent", event)
			}
		}
	}
}

func (a *AutoSender) sendToWebhook(msg messageModel.Message) (string, time.Time, error) {
	payload := fmt.Sprintf(`{"to": "%s", "content": "%s"}`, msg.To, msg.Content)
	headers := map[string]string{"Content-Type": "application/json"}
	if a.webhookKey != "" {
		headers["x-ins-auth-key"] = a.webhookKey
	}
	resp, err := a.httpSvc.Post(a.webhookURL, []byte(payload), headers)
	if err != nil {
		return "", time.Now(), err
	}
	defer resp.Body.Close()
	var respData struct {
		Message   string `json:"message"`
		MessageId string `json:"messageId"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return "", time.Now(), err
	}
	return respData.MessageId, time.Now(), nil
}
