package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"user/internal/config"
	"user/internal/db/models"
	"user/internal/utils/logger"
	"user/pkg/client/postgres"
)

type userServiceConsumerGroupHandler struct {
	db       *gorm.DB
	producer sarama.SyncProducer
}

func (userServiceConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (userServiceConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h userServiceConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
		checkUserAndRespond(h.producer, h.db, message)
		session.MarkMessage(message, "")
	}
	return nil
}

func main() {

	cfg, err := config.LoadAPIConfig()
	if err != nil {
		log.Fatalf("Failed to load user service API configs: %v", err)
	}

	zlog, _ := logger.InitLogger(cfg.Env)

	zlog.Info("user microservice started", zap.String("version", cfg.Version))

	db, err := postgres.Connect(cfg.Database)
	if err != nil {
		zlog.Fatal("Failed to connect to user DB", zap.String("error", err.Error()))
	}

	if db == nil {
		zlog.Fatal("Received nil for user database object after connection")
		return
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		zlog.Fatal("Failed to migrate user database: %v", zap.String("error", err.Error()))
	}

	producer, err := sarama.NewSyncProducer([]string{"kafka:19092"}, nil)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}

	//TODO: Temp For Testing
	//Kafka Initialization
	var consumer sarama.Consumer
	for retries := 0; retries < 10; retries++ {
		consumer, err = sarama.NewConsumer([]string{"kafka:19092"}, nil)
		if err == nil {
			break
		}
		log.Printf("Retrying to connect to Kafka consumer: attempt %d", retries+1)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}

	partitions, err := consumer.Partitions("user_requests")
	if err != nil {
		log.Fatalf("Failed to get partitions: %v", err)
	}

	for _, partition := range partitions {
		log.Printf("Consuming from partition %d", partition)
		partitionConsumer, err := consumer.ConsumePartition("user_requests", partition, sarama.OffsetOldest)
		if err != nil {
			log.Printf("Failed to consume partition %d: %v", partition, err)
			continue
		}

		go func(pc sarama.PartitionConsumer) {
			for message := range pc.Messages() {
				log.Printf("Received message from partition %d at offset %d: %s", message.Partition, message.Offset, string(message.Value))
				go checkUserAndRespond(producer, db, message)
			}
		}(partitionConsumer)
	}

	conf := sarama.NewConfig()
	conf.Version = sarama.V3_6_0_0 // Update this version to match your Kafka brokers version

	consumerGroup, err := sarama.NewConsumerGroup([]string{"kafka:19092"}, "user_service_group", conf)
	if err != nil {
		log.Fatalf("Failed to create consumer group: %v", err)
	}

	handler := userServiceConsumerGroupHandler{
		db:       db,
		producer: producer,
	}

	go func() {
		for {
			err := consumerGroup.Consume(context.Background(), []string{"user_requests"}, handler)
			if err != nil {
				log.Printf("Error from consumer group: %v", err)
			}
		}
	}()

	//go func() {
	//	defer func(consumer sarama.Consumer) {
	//		err := consumer.Close()
	//		if err != nil {
	//			log.Fatalf("Failed to close Kafka consumer: %v", err)
	//		}
	//	}(consumer)
	//
	//	partitions, err := consumer.Partitions("user_requests")
	//	if err != nil {
	//		log.Fatalf("Failed to get partitions: %v", err)
	//	}
	//
	//	for _, partition := range partitions {
	//		// Consume partition
	//		partitionConsumer, err := consumer.ConsumePartition("user_requests", partition, sarama.OffsetOldest)
	//		if err != nil {
	//			log.Printf("Failed to consume partition %d: %v", partition, err)
	//			continue
	//		}
	//
	//		go func(pc sarama.PartitionConsumer) {
	//			for message := range pc.Messages() {
	//				// Process each message
	//				go checkUserAndRespond(producer, db, message)
	//			}
	//			// Close the partition consumer after we finish processing messages
	//			err := pc.Close()
	//			if err != nil {
	//				log.Printf("Failed to close partition %d: %v", partition, err)
	//				return
	//			}
	//		}(partitionConsumer)
	//	}
	//}()

	//Consuming messages and routing them to handlers

	server := &http.Server{
		Addr: cfg.Server.Host + ":" + cfg.Server.Port,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zlog.Fatal("Failed to listen and serve", zap.String("error", err.Error()))

		}
	}()

	<-quit
	zlog.Info("Shutting down user microservice server...")

	//context with a timeout to allow for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server
	if err := server.Shutdown(ctx); err != nil {
		zlog.Fatal("Server forced to shutdown:", zap.Error(err))
	}

}

// //////////////////////TEMP///////TESTING/////////////
// TODO: Move to brokers service
type Message struct {
	RequestID string `json:"request_id"`
	Service   string `json:"service"`
	Action    string `json:"action"`
	Payload   string `json:"payload"`
}

/////////////////////////////////////////////////

func checkUserAndRespond(producer sarama.SyncProducer, db *gorm.DB, message *sarama.ConsumerMessage) {
	// Unmarshal Kafka message
	var msg Message
	err := json.Unmarshal(message.Value, &msg)
	if err != nil {
		log.Printf("Failed to unmarshal Kafka message: %v", err)
		return
	}

	// Perform action based on the message
	if msg.Action == "check_user" {
		userID := msg.Payload
		var user models.User
		err := db.First(&user, "id = ?", userID).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// User not found
				msg.Payload = "false"
			} else {
				log.Printf("Failed to retrieve user: %v", err)
				return
			}
		} else {
			// User found
			msg.Payload = fmt.Sprintf("true,%s,%s", user.Username, user.FirstName) // Respond with user details
		}

		// Marshal and send response
		responseData, err := json.Marshal(msg)
		if err != nil {
			log.Printf("Failed to marshal response data: %v", err)
			return
		}

		_, _, err = producer.SendMessage(&sarama.ProducerMessage{
			Topic: "user_responses",
			Value: sarama.ByteEncoder(responseData),
		})
		if err != nil {
			log.Printf("Failed to send response to Kafka: %v", err)
		}
	}
}
