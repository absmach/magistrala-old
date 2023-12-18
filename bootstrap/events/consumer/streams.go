// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package consumer

import (
	"context"
	"encoding/json"
	"time"

	"github.com/absmach/magistrala/bootstrap"
	"github.com/absmach/magistrala/pkg/clients"
	"github.com/absmach/magistrala/pkg/events"
)

const (
	thingRemove     = "thing.remove"
	thingDisconnect = "policy.delete"

	channelPrefix = "group."
	channelUpdate = channelPrefix + "update"
	channelDelete = channelPrefix + "delete"
)

type eventHandler struct {
	svc bootstrap.Service
}

// NewEventHandler returns new event store handler.
func NewEventHandler(svc bootstrap.Service) events.EventHandler {
	return &eventHandler{
		svc: svc,
	}
}

func (es *eventHandler) Handle(ctx context.Context, event events.Event) error {
	msg, err := event.Encode()
	if err != nil {
		return err
	}

	switch msg["operation"] {
	case thingRemove:
		rte := decodeRemoveThing(msg)
		err = es.svc.RemoveConfigHandler(ctx, rte.id)
	case thingDisconnect:
		dte := decodeDisconnectThing(msg)
		err = es.svc.DisconnectThingHandler(ctx, dte.channelID, dte.thingID)
	case channelUpdate:
		uce := decodeUpdateChannel(msg)
		err = es.handleUpdateChannel(ctx, uce)
	case channelDelete:
		dce := decodeDeleteChannel(msg)
		err = es.svc.RemoveChannelHandler(ctx, dce.id)
	}
	if err != nil {
		return err
	}

	return nil
}

func decodeRemoveThing(event map[string]interface{}) removeEvent {
	status := read(event, "status", "")
	st, err := clients.ToStatus(status)
	if err != nil {
		return removeEvent{}
	}
	switch st {
	case clients.EnabledStatus:
		return removeEvent{}
	case clients.DisabledStatus:
		return removeEvent{
			id: read(event, "id", ""),
		}
	default:
		return removeEvent{}
	}
}

func decodeUpdateChannel(event map[string]interface{}) updateChannelEvent {
	strmeta := read(event, "metadata", "{}")
	var metadata map[string]interface{}
	if err := json.Unmarshal([]byte(strmeta), &metadata); err != nil {
		metadata = map[string]interface{}{}
	}

	return updateChannelEvent{
		id:        read(event, "id", ""),
		name:      read(event, "name", ""),
		metadata:  metadata,
		updatedAt: readTime(event, "updated_at", time.Now()),
		updatedBy: read(event, "updated_by", ""),
	}
}

func decodeDeleteChannel(event map[string]interface{}) removeEvent {
	return removeEvent{
		id: read(event, "id", ""),
	}
}

func decodeDisconnectThing(event map[string]interface{}) disconnectEvent {
	return disconnectEvent{
		channelID: read(event, "chan_id", ""),
		thingID:   read(event, "thing_id", ""),
	}
}

func (es *eventHandler) handleUpdateChannel(ctx context.Context, uce updateChannelEvent) error {
	channel := bootstrap.Channel{
		ID:        uce.id,
		Name:      uce.name,
		Metadata:  uce.metadata,
		UpdatedAt: uce.updatedAt,
		UpdatedBy: uce.updatedBy,
	}
	return es.svc.UpdateChannelHandler(ctx, channel)
}

func read(event map[string]interface{}, key, def string) string {
	val, ok := event[key].(string)
	if !ok {
		return def
	}

	return val
}

func readTime(event map[string]interface{}, key string, def time.Time) time.Time {
	val, ok := event[key].(time.Time)
	if !ok {
		return def
	}

	return val
}
