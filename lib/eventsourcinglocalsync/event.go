package eventsourcinglocalsync

import (
	"time"

	"github.com/alejandroq/goventsourcing/pkg/eventsourcingiface"
)

type metadata struct {
	OriginStreamName string `json:"originStreamName"`
	TraceID          string `json:"traceID"`
}

type event struct {
	EventID          string                           `json:"eventID"`
	TransactionID    string                           `json:"transactionID"`
	LocalSequenceID  int                              `json:"sequenceID"`
	GlobalSequenceID int                              `json:"globalSequenceID"`
	Timestamp        time.Time                        `json:"timestamp"`
	Type             string                           `json:"type"`
	Metadata         eventsourcingiface.EventMetadata `json:"metadata"`
	Body             string                           `json:"body"`
	Version          int                              `json:"version"`
}

func (e event) GetTransactionID() string {
	return e.TransactionID
}

func (e event) GetLocalSequenceID() int {
	return e.LocalSequenceID
}

func (e event) GetGlobalSequenceID() int {
	return e.GlobalSequenceID
}

func (e event) GetTimestamp() time.Time {
	return e.Timestamp
}

func (e event) GetEventID() string {
	return e.EventID
}

func (e event) GetType() string {
	return e.Type
}

func (e event) GetMetadata() eventsourcingiface.EventMetadata {
	return e.Metadata
}

func (e event) GetBody() string {
	return e.Body
}

func (e event) GetVersion() int {
	return e.Version
}

func (e event) SetEventID(i string) eventsourcingiface.Event {
	e.EventID = i
	return e
}

func (e event) SetType(i string) eventsourcingiface.Event {
	e.Type = i
	return e
}

func (e event) SetMetadata(i eventsourcingiface.EventMetadata) eventsourcingiface.Event {
	e.Metadata = i
	return e
}

func (e event) SetBody(i string) eventsourcingiface.Event {
	e.Body = i
	return e
}

func (e event) SetVersion(i int) eventsourcingiface.Event {
	e.Version = i
	return e
}

func (m metadata) GetOriginStreamName() string {
	return m.OriginStreamName
}

func (m metadata) GetTraceID() string {
	return m.TraceID
}

func (m metadata) SetOriginStreamName(i string) eventsourcingiface.EventMetadata {
	m.OriginStreamName = i
	return m
}

func (m metadata) SetTraceID(i string) eventsourcingiface.EventMetadata {
	m.TraceID = i
	return m
}
