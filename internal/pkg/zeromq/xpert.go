// Copyright (C) 2023 IOTech Ltd

package zeromq

import (
	"fmt"

	"github.com/edgexfoundry/go-mod-messaging/v2/pkg/types"
)

func (client *zeromqClient) PublishBinaryData(data []byte, topic string) error {
	return fmt.Errorf("not supported PublishBinaryData func")
}

func (client *zeromqClient) SubscribeBinaryData(topics []types.TopicChannel, messageErrors chan error) error {
	return fmt.Errorf("not supported SubscribeBinaryData func")
}
