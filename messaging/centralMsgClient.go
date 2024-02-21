// Copyright (C) 2024 IOTech Ltd

package messaging

// CentralMsgClient is a low-level MessageClient interface implemented by Edge Central
type CentralMsgClient interface {
	// Unsubscribe ends the subscription from the provided topics
	Unsubscribe(topic string) error
}
