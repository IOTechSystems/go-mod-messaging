// Copyright (C) 2023 IOTech Ltd

package pkg

import (
	"fmt"
	"github.com/edgexfoundry/go-mod-messaging/v3/pkg/types"
)

func (n NoopClient) PublishBinaryData(data []byte, topic string) error {
	return fmt.Errorf("not supported PublishBinaryData func")
}

func (n NoopClient) SubscribeBinaryData(topics []types.TopicChannel, messageErrors chan error) error {
	return fmt.Errorf("not supported SubscribeBinaryData func")
}
