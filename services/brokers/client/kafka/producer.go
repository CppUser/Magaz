package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"log"
)

type Producer struct {
	Producer sarama.SyncProducer
}

func NewProducer(url []string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Retry.Max = 5
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Version = sarama.V3_6_0_0

	producer, err := sarama.NewSyncProducer(url, config)
	if err != nil {
		return nil, err
	}
	return &Producer{Producer: producer}, nil
}

func (p *Producer) SendMessage(serviceName, action, payload string) error {
	message := Message{
		RequestID: uuid.New().String(),
		Service:   serviceName,
		Action:    action,
		Payload:   payload,
	}

	msgBt, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	topic := fmt.Sprintf("%s_requests", serviceName)

	kmsg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(msgBt),
	}

	partition, offset, err := p.Producer.SendMessage(kmsg)
	if err != nil {
		return err
	}
	log.Printf("Message sent to partition %d at offset %d\n", partition, offset)
	return nil
}

func (kc *Client) SendMessageWithResponse(service, action, payload string) (*Message, error) {
	requestID := uuid.New().String()
	message := Message{
		RequestID: requestID,
		Service:   service,
		Action:    action,
		Payload:   payload,
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}

	topic := fmt.Sprintf("%s_requests", service)
	kafkaMessage := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(messageBytes),
	}

	partition, offset, err := kc.Producer.SendMessage(kafkaMessage)
	if err != nil {
		return nil, fmt.Errorf("failed to send message to Kafka: %w", err)
	}

	log.Printf("Message sent to partition %d at offset %d", partition, offset)

	// Prepare a response channel and store it using the RequestID
	responseChan := make(chan *Message, 1)
	kc.mux.Lock()
	kc.responseChans[requestID] = responseChan
	kc.mux.Unlock()

	// Wait for the response from the appropriate response topic
	select {
	case response := <-responseChan:
		return response, nil
	}
}
