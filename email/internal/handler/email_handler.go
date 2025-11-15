package handler

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/MamangRust/simple_microservice_ecommerce/email/internal/mailer"
)

type emailHandler struct {
	ctx    context.Context
	mailer *mailer.Mailer
}

func NewEmailHandler(ctx context.Context, mailer *mailer.Mailer) *emailHandler {
	return &emailHandler{
		ctx:    ctx,
		mailer: mailer,
	}
}

func (h *emailHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *emailHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *emailHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	start := time.Now()
	total := 0
	failed := 0

	log.Println("[EmailHandler] Started consuming messages...")

	for msg := range claim.Messages() {
		total++
		var payload map[string]interface{}

		if err := json.Unmarshal(msg.Value, &payload); err != nil {
			failed++
			log.Printf("[EmailHandler] ❌ Failed to unmarshal message: %v\n", err)
			continue
		}

		email, _ := payload["email"].(string)
		subject, _ := payload["subject"].(string)
		body, _ := payload["body"].(string)

		if err := h.mailer.Send(email, subject, body); err != nil {
			failed++
			log.Printf("[EmailHandler] ❌ Failed to send email to %s: %v\n", email, err)
		} else {
			log.Printf("[EmailHandler] ✅ Email sent to %s", email)
		}

		sess.MarkMessage(msg, "")
	}

	duration := time.Since(start)
	log.Printf("[EmailHandler] Finished consuming. Total: %d, Failed: %d, Duration: %.2fs\n",
		total, failed, duration.Seconds())

	return nil
}
