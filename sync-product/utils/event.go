package utils

import "time"

type EventMessage struct {
	Header struct {
		Version       string        `json:"version,omitempty"`
		Timestamp     time.Time     `json:"timestamp,omitempty"`
		OrgService    string        `json:"orgService,omitempty"`
		From          string        `json:"from,omitempty"`
		Channel       string        `json:"channel,omitempty"`
		Broker        string        `json:"broker,omitempty"`
		UseCase       string        `json:"useCase,omitempty"`
		Session       string        `json:"session,omitempty"`
		Transaction   string        `json:"transaction,omitempty"`
		Communication string        `json:"communication,omitempty"`
		GroupTags     []interface{} `json:"groupTags,omitempty"`
		Identity      struct {
			Device int64 `json:"device,omitempty"`
		} `json:"identity,omitempty"`
		BaseAPIVersion string `json:"baseApiVersion,omitempty"`
		SchemaVersion  string `json:"schemaVersion,omitempty"`
		InstanceData   string `json:"instanceData,omitempty"`
	} `json:"header"`
	Body interface {
	} `json:"body"`
}
