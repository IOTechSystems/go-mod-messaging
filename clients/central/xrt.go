// Copyright (C) 2023 IOTech Ltd

package central

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/edgexfoundry/go-mod-core-contracts/v3/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/v3/errors"
	"github.com/edgexfoundry/go-mod-core-contracts/v3/xrtmodels"
	"github.com/edgexfoundry/go-mod-messaging/v3/clients/central/utils"
	"github.com/edgexfoundry/go-mod-messaging/v3/clients/interfaces"
	"github.com/edgexfoundry/go-mod-messaging/v3/messaging"
	"github.com/edgexfoundry/go-mod-messaging/v3/pkg/types"
)

const (
	clientName = "central-go-client"
)

// xrtClient implements the client of MQTT management API, https://docs.iotechsys.com/edge-xrt21/mqtt-management/mqtt-management.html
type xrtClient struct {
	lc              logger.LoggingClient
	requestMap      utils.RequestMap
	messageBus      messaging.MessageClient
	requestTopic    string
	replyTopic      string
	responseTimeout time.Duration

	clientOptions *ClientOptions
}

type MessageHandler func(message types.MessageEnvelope)

type Subscription struct {
	topicChannel   types.TopicChannel
	messageHandler MessageHandler
}

type ClientOptions struct {
	*CommandOptions
	*DiscoveryOptions
	*StatusOptions
}

// CommandOptions provides the config for sending the request to manage components
type CommandOptions struct {
	CommandTopic string
}

// DiscoveryOptions provides the config for sending the discovery request like discovery:trigger, device:scan
type DiscoveryOptions struct {
	DiscoveryTopic          string
	DiscoveryMessageHandler MessageHandler
	DiscoveryDuration       time.Duration
	DiscoveryTimeout        time.Duration
}

// StatusOptions provides the config for subscribing the XRT status
type StatusOptions struct {
	StatusTopic          string
	StatusMessageHandler MessageHandler
}

func NewXrtClient(ctx context.Context, messageBus messaging.MessageClient, requestTopic string, replyTopic string,
	responseTimeout time.Duration, lc logger.LoggingClient, clientOptions *ClientOptions) (interfaces.XrtClient, errors.EdgeX) {
	client := &xrtClient{
		lc:              lc,
		requestMap:      utils.NewRequestMap(),
		messageBus:      messageBus,
		requestTopic:    requestTopic,
		replyTopic:      replyTopic,
		responseTimeout: responseTimeout,
		clientOptions:   clientOptions,
	}

	err := initSubscriptions(ctx, client, clientOptions, lc)
	if err != nil {
		return client, errors.NewCommonEdgeX(errors.Kind(err), "failed to init subscriptions", err)
	}

	return client, nil
}

func NewClientOptions(commandOptions *CommandOptions, discoveryOptions *DiscoveryOptions, statusOptions *StatusOptions) *ClientOptions {
	return &ClientOptions{
		CommandOptions:   commandOptions,
		DiscoveryOptions: discoveryOptions,
		StatusOptions:    statusOptions,
	}
}

func NewCommandOptions(commandTopic string) *CommandOptions {
	return &CommandOptions{
		CommandTopic: commandTopic,
	}
}

func NewDiscoveryOptions(discoveryTopic string, discoveryMessageHandler MessageHandler, discoveryDuration, discoveryTimeout time.Duration) *DiscoveryOptions {
	return &DiscoveryOptions{
		DiscoveryTopic:          discoveryTopic,
		DiscoveryMessageHandler: discoveryMessageHandler,
		DiscoveryDuration:       discoveryDuration,
		DiscoveryTimeout:        discoveryTimeout,
	}
}

func NewStatusOptions(statusTopic string, statusMessageHandler MessageHandler) *StatusOptions {
	return &StatusOptions{
		StatusTopic:          statusTopic,
		StatusMessageHandler: statusMessageHandler,
	}
}

func (c *xrtClient) SetResponseTimeout(responseTimeout time.Duration) {
	c.responseTimeout = responseTimeout
}

// sendXrtRequest sends general request to XRT
func (c *xrtClient) sendXrtRequest(ctx context.Context, requestId string, request interface{}, response interface{}) errors.EdgeX {
	return c.sendXrtRequestWithTimeout(ctx, c.requestTopic, requestId, request, response, c.responseTimeout)
}

// sendXrtDiscoveryRequest sends discovery request to XRT
func (c *xrtClient) sendXrtDiscoveryRequest(ctx context.Context, requestId string, request interface{}, response interface{}) errors.EdgeX {
	if c.clientOptions == nil || c.clientOptions.DiscoveryOptions == nil {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "please provide DiscoveryOptions for the discovery request", nil)
	}
	timeout := time.Duration(c.responseTimeout.Nanoseconds() + c.clientOptions.DiscoveryDuration.Nanoseconds() + c.clientOptions.DiscoveryTimeout.Nanoseconds())
	return c.sendXrtRequestWithTimeout(ctx, c.requestTopic, requestId, request, response, timeout)
}

// sendXrtCommandRequest sends command request to XRT
func (c *xrtClient) sendXrtCommandRequest(ctx context.Context, requestId string, request interface{}, response interface{}) errors.EdgeX {
	if c.clientOptions == nil || c.clientOptions.CommandOptions == nil {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "please provide CommandOptions for the command request", nil)
	}
	return c.sendXrtRequestWithTimeout(ctx, c.clientOptions.CommandOptions.CommandTopic, requestId, request, response, c.responseTimeout)
}

func (c *xrtClient) sendXrtRequestWithTimeout(ctx context.Context, requestTopic string, requestId string, request interface{}, response interface{}, responseTimeout time.Duration) errors.EdgeX {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return errors.NewCommonEdgeXWrapper(err)
	}

	// Before publishing the request, we should create responseChan to receive the response from XRT
	c.requestMap.Add(requestId)

	err = c.messageBus.PublishBinaryData(jsonData, requestTopic)
	if err != nil {
		return errors.NewCommonEdgeX(errors.Kind(err), "failed to send the XRT request", err)
	}

	cmdResponseBytes, err := utils.FetchXRTResponse(ctx, requestId, c.requestMap, responseTimeout)
	if err != nil {
		return errors.NewCommonEdgeXWrapper(err)
	}

	err = json.Unmarshal(cmdResponseBytes, response)
	if err != nil {
		return errors.NewCommonEdgeX(errors.KindServerError, "failed to JSON decoding command response: %v", err)
	}

	// handle error result from the XRT
	var commonResponse xrtmodels.CommonResponse
	err = json.Unmarshal(cmdResponseBytes, &commonResponse)
	if err != nil {
		return errors.NewCommonEdgeX(errors.KindServerError, "failed to JSON decoding command response: %v", err)
	}
	if commonResponse.Result.Error() != nil {
		return errors.NewCommonEdgeXWrapper(commonResponse.Result.Error())
	}
	return nil
}

// sendXrtRequestWithSubTimeout publish the xrt request and wait for responses from multiple xrt nodes for the specific subscribe timeout
func (c *xrtClient) sendXrtRequestWithSubTimeout(ctx context.Context, requestTopic string, requestId string, request any,
	response any, subscribeTimeout time.Duration) errors.EdgeX {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return errors.NewCommonEdgeXWrapper(err)
	}

	// Before publishing the request, we should create responseChan to receive the response from XRT
	c.requestMap.Add(requestId)

	err = c.messageBus.PublishBinaryData(jsonData, requestTopic)
	if err != nil {
		return errors.NewCommonEdgeX(errors.Kind(err), "failed to send the XRT request", err)
	}

	edgexErr := utils.FetchXRTResWithSubTimeout(ctx, requestId, c.requestMap, subscribeTimeout, response)
	if edgexErr != nil {
		return errors.NewCommonEdgeXWrapper(edgexErr)
	}

	return nil
}

func initSubscriptions(ctx context.Context, xrtClient *xrtClient, clientOptions *ClientOptions, lc logger.LoggingClient) errors.EdgeX {
	subscriptions := createSubscriptions(xrtClient, clientOptions)
	messageErrors := make(chan error)
	// Create goroutine to handle MessageBus errors
	go func() {
		for {
			select {
			case <-ctx.Done():
				lc.Info("Exiting waiting for MessageBus errors")
				return
			case msgErr := <-messageErrors:
				lc.Errorf("error receiving message from bus, %s", msgErr.Error())
			}
		}
	}()
	// Create goroutines to receive message for each subscription
	for _, sub := range subscriptions {
		go func(subscription Subscription) {
			lc.Infof("Waiting for messages from the MessageBus on the '%s' topic", subscription.topicChannel.Topic)
			for {
				select {
				case <-ctx.Done():
					lc.Infof("Exiting waiting for MessageBus '%s' topic messages", subscription.topicChannel.Topic)
					return
				case message := <-subscription.topicChannel.Messages:
					lc.Debugf("Received message from the topic %s", subscription.topicChannel.Topic)
					subscription.messageHandler(message)
				}
			}
		}(sub)
		err := xrtClient.messageBus.SubscribeBinaryData([]types.TopicChannel{sub.topicChannel}, messageErrors)
		if err != nil {
			return errors.NewCommonEdgeX(errors.Kind(err), fmt.Sprintf("subscribe to topic '%s' failed", sub.topicChannel.Topic), err)
		}
		lc.Debugf("Subscribed to %s", sub.topicChannel.Topic)
	}
	return nil
}

func commandReplyHandler(requestMap utils.RequestMap, lc logger.LoggingClient) MessageHandler {
	return func(message types.MessageEnvelope) {
		var response xrtmodels.BaseResponse
		err := json.Unmarshal(message.Payload, &response)
		if err != nil {
			lc.Warnf("failed to parse XRT reply, message:%s, err: %v", string(message.Payload), err)
			return
		}
		resChan, ok := requestMap.Get(response.RequestId)
		if !ok {
			lc.Debugf("deprecated response from the XRT, it might be caused by timeout or unknown error, topic: %s, message:%s", message.ReceivedTopic, string(message.Payload))
			return
		}

		resChan <- message.Payload
	}
}

func createSubscriptions(xrtClient *xrtClient, clientOptions *ClientOptions) []Subscription {
	var subscriptions []Subscription
	subscriptions = append(subscriptions, subscription(xrtClient.replyTopic, commandReplyHandler(xrtClient.requestMap, xrtClient.lc)))

	if clientOptions == nil {
		return subscriptions
	}
	if clientOptions.DiscoveryOptions != nil {
		if clientOptions.DiscoveryOptions.DiscoveryTopic != "" && clientOptions.DiscoveryOptions.DiscoveryMessageHandler != nil {
			subscriptions = append(subscriptions, subscription(clientOptions.DiscoveryOptions.DiscoveryTopic, clientOptions.DiscoveryOptions.DiscoveryMessageHandler))
		}
	}
	if clientOptions.StatusOptions != nil {
		if clientOptions.StatusOptions.StatusTopic != "" && clientOptions.StatusOptions.StatusMessageHandler != nil {
			subscriptions = append(subscriptions, subscription(clientOptions.StatusOptions.StatusTopic, clientOptions.StatusOptions.StatusMessageHandler))
		}
	}
	return subscriptions
}

func subscription(topic string, messageHandler MessageHandler) Subscription {
	return Subscription{
		topicChannel: types.TopicChannel{
			Topic:    topic,
			Messages: make(chan types.MessageEnvelope),
		},
		messageHandler: messageHandler,
	}
}

func (c *xrtClient) Close() errors.EdgeX {
	err := c.messageBus.Disconnect()
	if err != nil {
		return errors.NewCommonEdgeXWrapper(err)
	}
	return nil
}
