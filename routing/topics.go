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

package routing

import (
	"fmt"
	"net/url"

	"github.com/eclipse/ditto-clients-golang/protocol"
)

const (
	keyMessageID       = "$.mid"
	keyContentType     = "$.ct"
	keyContentEncoding = "$.ce"
	contentType        = "application/json"
	contentEncoding    = "utf-8"

	remoteCloudTopicFmt     = "devices/%s/messages/devicebound/#"
	remoteTelemetryTopicFmt = "devices/%s/messages/events/%s"

	localCmdTopicLongFmt  = "command//%s:%s/req/%s/%s"
	localCmdTopicShortFmt = "c//%s:%s/q/%s/%s"
)

// CreateRemoteCloudTopic constructs the remote MQTT topic for receiving C2D messages from an Azure IoT Hub device.
func CreateRemoteCloudTopic(deviceID string) string {
	return fmt.Sprintf(remoteCloudTopicFmt, deviceID)
}

// CreateTelemetryTopic constructs the MQTT topic for sending telemetry data to an Azure IoT Hub device.
func CreateTelemetryTopic(deviceID, msgID string) string {
	msgProps := make(url.Values, 3)
	msgProps[keyContentType] = []string{contentType}
	msgProps[keyContentEncoding] = []string{contentEncoding}
	if msgID != "" {
		msgProps[keyMessageID] = []string{msgID}
	}
	return fmt.Sprintf(remoteTelemetryTopicFmt, deviceID, msgProps.Encode())
}

// CreateLocalCmdTopicLong constructs the local MQTT topic for receiving C2D messages from an Azure IoT Hub device.
func CreateLocalCmdTopicLong(env *protocol.Envelope) string {
	return createLocalCmdTopic(localCmdTopicLongFmt, env)
}

// CreateLocalCmdTopicShort constructs the local MQTT topic for receiving C2D messages from an Azure IoT Hub device.
func CreateLocalCmdTopicShort(env *protocol.Envelope) string {
	return createLocalCmdTopic(localCmdTopicShortFmt, env)
}

func createLocalCmdTopic(template string, env *protocol.Envelope) string {
	return fmt.Sprintf(template,
		env.Topic.Namespace,
		env.Topic.EntityName,
		env.Headers.CorrelationID(),
		env.Topic.Action,
	)
}
