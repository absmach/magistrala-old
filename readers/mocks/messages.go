// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package mocks

import (
	"encoding/json"
	"sync"

	"github.com/absmach/magistrala/pkg/transformers/senml"
	"github.com/absmach/magistrala/readers"
)

var _ readers.MessageRepository = (*messageRepositoryMock)(nil)

type messageRepositoryMock struct {
	mutex    sync.Mutex
	messages map[string][]readers.Message
}

// NewMessageRepository returns mock implementation of message repository.
func NewMessageRepository(chanID string, messages []readers.Message) readers.MessageRepository {
	repo := map[string][]readers.Message{
		chanID: messages,
	}

	return &messageRepositoryMock{
		mutex:    sync.Mutex{},
		messages: repo,
	}
}

func (repo *messageRepositoryMock) ReadAll(chanID string, rpm readers.PageMetadata) (readers.MessagesPage, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	if rpm.Format != "" && rpm.Format != "messages" {
		return readers.MessagesPage{}, nil
	}

	var query map[string]interface{}
	meta, err := json.Marshal(rpm)
	if err != nil {
		return readers.MessagesPage{}, err
	}
	if err := json.Unmarshal(meta, &query); err != nil {
		return readers.MessagesPage{}, err
	}

	var msgs []readers.Message
	for _, m := range repo.messages[chanID] {
		mes := m.(senml.Message)

		ok := true

		for name := range query {
			switch name {
			case "subtopic":
				if rpm.Subtopic != mes.Subtopic {
					ok = false
				}
			case "publisher":
				if rpm.Publisher != mes.Publisher {
					ok = false
				}
			case "name":
				if rpm.Name != mes.Name {
					ok = false
				}
			case "protocol":
				if rpm.Protocol != mes.Protocol {
					ok = false
				}
			case "v":
				if mes.Value == nil {
					ok = false
				}

				val, okQuery := query["comparator"]
				if okQuery {
					switch val.(string) {
					case readers.LowerThanKey:
						if mes.Value != nil &&
							*mes.Value >= rpm.Value {
							ok = false
						}
					case readers.LowerThanEqualKey:
						if mes.Value != nil &&
							*mes.Value > rpm.Value {
							ok = false
						}
					case readers.GreaterThanKey:
						if mes.Value != nil &&
							*mes.Value <= rpm.Value {
							ok = false
						}
					case readers.GreaterThanEqualKey:
						if mes.Value != nil &&
							*mes.Value < rpm.Value {
							ok = false
						}
					case readers.EqualKey:
					default:
						if mes.Value != nil &&
							*mes.Value != rpm.Value {
							ok = false
						}
					}
				}
			case "vb":
				if mes.BoolValue == nil ||
					(mes.BoolValue != nil &&
						*mes.BoolValue != rpm.BoolValue) {
					ok = false
				}
			case "vs":
				if mes.StringValue == nil ||
					(mes.StringValue != nil &&
						*mes.StringValue != rpm.StringValue) {
					ok = false
				}
			case "vd":
				if mes.DataValue == nil ||
					(mes.DataValue != nil &&
						*mes.DataValue != rpm.DataValue) {
					ok = false
				}
			case "from":
				if mes.Time < rpm.From {
					ok = false
				}
			case "to":
				if mes.Time >= rpm.To {
					ok = false
				}
			}

			if !ok {
				break
			}
		}

		if ok {
			msgs = append(msgs, m)
		}
	}

	numOfMessages := uint64(len(msgs))

	if rpm.Offset >= numOfMessages {
		return readers.MessagesPage{}, nil
	}

	if rpm.Limit < 1 {
		return readers.MessagesPage{}, nil
	}

	end := rpm.Offset + rpm.Limit
	if rpm.Offset+rpm.Limit > numOfMessages {
		end = numOfMessages
	}

	return readers.MessagesPage{
		PageMetadata: rpm,
		Total:        uint64(len(msgs)),
		Messages:     msgs[rpm.Offset:end],
	}, nil
}
