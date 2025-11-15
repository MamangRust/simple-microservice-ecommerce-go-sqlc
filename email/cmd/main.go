package main

import (
	"context"
	"log"

	"github.com/MamangRust/simple_microservice_ecommerce/email/internal/config"
	"github.com/MamangRust/simple_microservice_ecommerce/email/internal/handler"
	"github.com/MamangRust/simple_microservice_ecommerce/email/internal/mailer"
	"github.com/MamangRust/simple_microservice_ecommerce/email/pkg/dotenv"
	"github.com/MamangRust/simple_microservice_ecommerce/email/pkg/kafka"
	"github.com/spf13/viper"
)

func main() {
	if err := dotenv.Viper(); err != nil {
		log.Fatalf("[MAIN] ❌ Failed to load .env file: %v", err)
	}

	ctx := context.Background()

	cfg := config.Config{
		KafkaBrokers: []string{viper.GetString("KAFKA_BROKERS")},
		SMTPServer:   viper.GetString("SMTP_SERVER"),
		SMTPPort:     viper.GetInt("SMTP_PORT"),
		SMTPUser:     viper.GetString("SMTP_USER"),
		SMTPPass:     viper.GetString("SMTP_PASS"),
	}

	m := mailer.NewMailer(
		ctx,
		cfg.SMTPServer,
		cfg.SMTPPort,
		cfg.SMTPUser,
		cfg.SMTPPass,
	)

	h := handler.NewEmailHandler(ctx, m)

	myKafka := kafka.NewKafka(cfg.KafkaBrokers)

	log.Println("[MAIN] 🚀 Starting Kafka consumers...")

	err := myKafka.StartConsumers([]string{
		"email-service-topic-auth-register",
		"email-service-topic-auth-forgot-password",
		"email-service-topic-auth-verify-code-success",
		"email-service-topic-merchant-create",
		"email-service-topic-merchant-update-status",
		"email-service-topic-merchant-document-create",
		"email-service-topic-merchant-document-update-status",
		"email-service-topic-transaction-create",
	}, "email-service-group", h)

	if err != nil {
		log.Fatalf("[MAIN] ❌ Error starting consumer: %v", err)
	}

	log.Println("[MAIN] ✅ Email service is running and listening for Kafka messages...")
	select {}
}
