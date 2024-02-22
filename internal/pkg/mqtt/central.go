// Copyright (C) 2023 IOTech Ltd

package mqtt

import (
	pahoMqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/edgexfoundry/go-mod-messaging/v3/pkg/types"
)

const (
	UnsubscribeOperation = "Unsubscribe"
)

func (mc *Client) PublishBinaryData(data []byte, topic string) error {
	optionsReader := mc.mqttClient.OptionsReader()

	return getTokenError(
		mc.mqttClient.Publish(
			topic,
			optionsReader.WillQos(),
			optionsReader.WillRetained(),
			data),
		optionsReader.ConnectTimeout(),
		PublishOperation,
		"Unable to publish message")
}

func (mc *Client) SubscribeBinaryData(topics []types.TopicChannel, messageErrors chan error) error {
	optionsReader := mc.mqttClient.OptionsReader()

	for _, topic := range topics {
		handler := newBinaryDataMessageHandler(topic.Messages)
		qos := optionsReader.WillQos()

		// Since the MQTT client might try to subscribe to the same topic and get the error 'not currently connected and ResumeSubs not set',
		// we need to unsubscribe the topic before subscribing to prevent the error.
		token := mc.mqttClient.Unsubscribe(topic.Topic)
		err := getTokenError(token, optionsReader.ConnectTimeout(), UnsubscribeOperation, "Failed to unsubscribe")
		if err != nil {
			return err
		}
		token = mc.mqttClient.Subscribe(topic.Topic, qos, handler)
		err = getTokenError(token, optionsReader.ConnectTimeout(), SubscribeOperation, "Failed to create subscription")
		if err != nil {
			return err
		}

		mc.existingSubscriptions[topic.Topic] = existingSubscription{
			topic:   topic.Topic,
			qos:     qos,
			handler: handler,
			errors:  messageErrors,
		}

	}

	return nil
}

// newBinaryDataMessageHandler creates a function which propagates the received messages to the proper channel.
func newBinaryDataMessageHandler(messageChannel chan<- types.MessageEnvelope) pahoMqtt.MessageHandler {
	return func(client pahoMqtt.Client, message pahoMqtt.Message) {
		// Use MessageEnvelope.Payload to store the binary data instead of unmarshalling binary to MessageEnvelope
		messageEnvelope := types.NewMessageEnvelopeForRequest(message.Payload(), nil)
		messageEnvelope.ReceivedTopic = message.Topic()
		messageChannel <- messageEnvelope
	}
}
