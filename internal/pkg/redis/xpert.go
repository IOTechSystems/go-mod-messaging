// Copyright (C) 2023 IOTech Ltd

package redis

import (
	"github.com/edgexfoundry/go-mod-messaging/v2/internal/pkg"
	"github.com/edgexfoundry/go-mod-messaging/v2/pkg/types"
)

func (c Client) PublishBinaryData(data []byte, topic string) error {
	if c.publishClient == nil {
		return pkg.NewMissingConfigurationErr("PublishHostInfo", "Unable to create a connection for publishing")
	}

	if topic == "" {
		// Empty topics are not allowed for Redis
		return pkg.NewInvalidTopicErr("", "Unable to publish to the invalid topic")
	}

	topic = convertToRedisTopicScheme(topic)
	return c.publishClient.SendBinaryData(topic, data)
}

func (c Client) SubscribeBinaryData(topics []types.TopicChannel, messageErrors chan error) error {
	if c.subscribeClient == nil {
		return pkg.NewMissingConfigurationErr("SubscribeHostInfo", "Unable to create a connection for subscribing")
	}

	err := c.validateTopics(topics)
	if err != nil {
		return err
	}

	for i := range topics {
		go func(topic types.TopicChannel) {
			topicName := convertToRedisTopicScheme(topic.Topic)
			messageChannel := topic.Messages
			for {
				message, err := c.subscribeClient.ReceiveBinaryData(topicName)
				if err != nil {
					messageErrors <- err
					continue
				}

				message.ReceivedTopic = convertFromRedisTopicScheme(message.ReceivedTopic)

				messageChannel <- *message
			}
		}(topics[i])
	}

	return nil
}

func (g *goRedisWrapper) SendBinaryData(topic string, data []byte) error {
	_, err := g.wrappedClient.Publish(topic, data).Result()
	if err != nil {
		return err
	}
	return nil
}

func (g *goRedisWrapper) ReceiveBinaryData(topic string) (*types.MessageEnvelope, error) {
	subscription := g.getSubscription(topic)
	data, err := subscription.ReceiveMessage()
	if err != nil {
		return nil, err
	}
	// Use MessageEnvelope.Payload to store the binary data instead of unmarshalling binary to MessageEnvelope
	messageEnvelope := types.NewMessageEnvelopeForRequest([]byte(data.Payload), nil)
	messageEnvelope.ReceivedTopic = data.Channel
	return &messageEnvelope, nil
}
