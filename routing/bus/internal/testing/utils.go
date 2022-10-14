// Copyright (c) 2022 Contributors to the Eclipse Foundation
//
// See the NOTICE file(s) distributed with this work for additional
// information regarding copyright ownership.
//
// This program and the accompanying materials are made available under the
// terms of the Eclipse Public License 2.0 which is available at
// https://www.eclipse.org/legal/epl-2.0, or the Apache License, Version 2.0
// which is available at https://www.apache.org/licenses/LICENSE-2.0.
//
// SPDX-License-Identifier: EPL-2.0 OR Apache-2.0

package test

import (
	"context"
	"reflect"
	"testing"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/eclipse-kanto/azure-connector/config"
	"github.com/eclipse-kanto/azure-connector/routing/message/handlers"

	"github.com/stretchr/testify/assert"
)

const (
	fieldName           = "name"
	fieldSubscriber     = "subscriber"
	fieldSubscribeTopic = "subscribeTopic"
	fieldPublisher      = "publisher"
	fieldPublishTopic   = "publishTopic"
	fieldHandlerFunc    = "handlerFunc"
	fieldHandlers       = "handlers"
)

type dummySubscriber struct{}

func (s dummySubscriber) Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error) {
	return make(chan *message.Message), nil
}
func (s dummySubscriber) Close() error { return nil }

// NewDummySubscriber instantiates a new dummy Watermill subscriber.
func NewDummySubscriber() message.Subscriber {
	return dummySubscriber{}
}

type dummyMessageHandler struct {
	handleName string
	topics     string
	initErr    error
	handleErr  error
}

func (h *dummyMessageHandler) Init(connInfo *config.RemoteConnectionInfo) error {
	return h.initErr
}

func (h *dummyMessageHandler) HandleMessage(msg *message.Message) ([]*message.Message, error) {
	if h.handleErr != nil {
		return nil, h.handleErr
	}
	msg.Metadata["handler_name"] = h.handleName
	return []*message.Message{msg}, nil
}

func (h *dummyMessageHandler) Name() string {
	return h.handleName
}

func (h *dummyMessageHandler) Topics() string {
	return h.topics
}

// NewDummyCommandHandler instantiates a new dummy command handler.
func NewDummyCommandHandler(handlerName string, initErr error, handleErr error) handlers.CommandHandler {
	return &dummyMessageHandler{
		handleName: handlerName,
		initErr:    initErr,
		handleErr:  handleErr,
	}
}

// NewDummyTelemetryHandler instantiates a new dummy telemetry handler.
func NewDummyTelemetryHandler(handlerName string, topic string, initErr error) handlers.TelemetryHandler {
	return &dummyMessageHandler{
		handleName: handlerName,
		topics:     topic,
		initErr:    initErr,
	}
}

// AssertRouterHandler asserts a Watermill router handler.
func AssertRouterHandler(t *testing.T, expectedHandlerName, expectedSubcribeTopic, expectedPublishTopic string, refHandler reflect.Value) {
	handlerName := refHandler.FieldByName(fieldName)
	assert.Equal(t, expectedHandlerName, handlerName.String())
	subscribeTopic := refHandler.FieldByName(fieldSubscribeTopic)
	assert.Equal(t, expectedSubcribeTopic, subscribeTopic.String())
	subscriber := refHandler.FieldByName(fieldSubscriber)
	assert.Equal(t, false, subscriber.IsZero())
	publishTopic := refHandler.FieldByName(fieldPublishTopic)
	assert.Equal(t, expectedPublishTopic, publishTopic.String())
	publisher := refHandler.FieldByName(fieldPublisher)
	assert.Equal(t, false, publisher.IsZero())
	handlerFunc := refHandler.FieldByName(fieldHandlerFunc)
	assert.Equal(t, false, handlerFunc.IsZero())
}

// AssertNoRouterHandlers asserts no handlers are registered in a Watermill router.
func AssertNoRouterHandlers(t *testing.T, router *message.Router) {
	refRouterPtr := reflect.ValueOf(router)
	refRouter := reflect.Indirect(refRouterPtr)
	refHandlers := refRouter.FieldByName(fieldHandlers)
	assert.Equal(t, 0, refHandlers.Len())
}
