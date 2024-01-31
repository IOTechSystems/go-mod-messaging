// Copyright (C) 2023 IOTech Ltd

package redis

import "github.com/edgexfoundry/go-mod-messaging/v3/pkg/types"

func (r *SubscriptionRedisClientMock) PublishBinaryData(topic string, data []byte) error {
	return nil
}

func (r *SubscriptionRedisClientMock) ReceiveBinaryData(topic string) (*types.MessageEnvelope, error) {
	return &types.MessageEnvelope{}, nil
}

func (r *SubscriptionRedisClientMock) SendBinaryData(topic string, data []byte) error {
	return nil
}
