package eventsourcingiface

import (
	"context"
	"time"
)

//Event are individual records in the ledger of the event bus
//that make up the persistency layer and are dynamically
//allocated to actors.
type Event interface {
	//EventBus generates the following 4 attributes upon
	//event write.
	GetTransactionID() string
	GetLocalSequenceID() int
	GetGlobalSequenceID() int
	GetTimestamp() time.Time

	//Service layer defined getter setter attributes.
	GetEventID() string
	GetType() string
	GetMetadata() EventMetadata
	GetBody() string
	GetVersion() int

	SetEventID(string) Event
	SetType(string) Event
	SetMetadata(EventMetadata) Event
	SetBody(string) Event
	SetVersion(int) Event
}

//EventMetadata are optional identifiers appended to Events
//which occasions functionality downstream.
//can be the difference between orchestration and choreography:
//https://stackoverflow.com/questions/4127241/orchestration-vs-choreography
type EventMetadata interface {
	GetOriginStreamName() string
	GetTraceID() string

	SetOriginStreamName(string) EventMetadata
	SetTraceID(string) EventMetadata
}

//Subscriber applies messages. Implementing domains, would
//benefit from an identity state (a zero value so to speak)
//and associative principles.
//Apply returns void as-is intended to introduce side-effects
//in the bounded context. The context is passed along for resource
//cleaning if need-be.
//StartWith is a flexible setup method.
//It is generally reccomended that the
type Subscriber interface {
	StartWith(EventBus)
	Apply(context.Context, Event)
}

//EventBus centralizes access to an event source persistency layer
//and therefore should be disciplined as the abstraction for said.
type EventBus interface {
	//Subscribe to a stream name with a Subscriber
	Subscribe(string, Subscriber) error

	//Write a message to the EventBus
	Write(string, Event) error

	//Read from the EventBus starting at N position and consume J records
	Read(string, int, int) ([]Event, error)

	//Conviniently create a new event
	NewEvent() Event

	//Conviniently create new event metadata
	NewEventMetadata() EventMetadata
}
