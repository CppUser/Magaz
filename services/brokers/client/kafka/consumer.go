package kafka

import (
	"fmt"
	"github.com/IBM/sarama"
	"log"
)

func (kc *Client) ConsumeResponses(service string) {
	topic := fmt.Sprintf("%s_responses", service)
	partitionConsumer, err := kc.Consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to start consumer for topic %s: %v", topic, err)
	}

	defer partitionConsumer.Close()

	for message := range partitionConsumer.Messages() {
		var msg Message
		err := json.Unmarshal(message.Value, &msg)
		if err != nil {
			log.Printf("failed to unmarshal message: %v", err)
			continue
		}

		// Check if there's a waiting channel for the response
		kc.mux.Lock()
		if responseChan, ok := kc.responseChans[msg.RequestID]; ok {
			// Send the response to the waiting channel
			responseChan <- &msg
			close(responseChan)
			delete(kc.responseChans, msg.RequestID)
		}
		kc.mux.Unlock()
	}
}
