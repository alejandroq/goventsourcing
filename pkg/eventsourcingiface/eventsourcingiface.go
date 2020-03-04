//Package eventsourcingiface serves to provide interface abstractions
//over event sourcing requirements. Specific implementations of said
//abstractions could facilitate PostgreSQL as a source for example.
//An explicit goal of this package is to be general enough to fit
//numerous cases and be declarative enough to help bridge the complexity
//gap of the event sourcing paradigm for newcomers. In this, it functions
//as a specification.
//
//As these ifaces are meant to be core and minimal, convinience
//methods should be in "extender" interfaces in a boosting package.
//
//When selecting a data store for your event sourced bounded context,
//the following rules should be facilitated:
//1. optimistic concurrency
//2. sequence
//3. immutability
//
//
//If you have any feedback on how to improve eventsourcingiface and this
//interpretation of the specification, be sure to create an issue or
//tweet the author at @redpause.
package eventsourcingiface

import (
	"context"
	"time"
)

//Event are individual records in the ledger stored in
//the event store that make up the persistency layer.
//These events are distributed to actors in whichever way
//implementations choose to do so.
type Event interface {
	//Event bus should generate the following 4 attributes upon an event write
	GetTransactionID() string
	GetLocalSequenceID() int
	GetGlobalSequenceID() int
	GetTimestamp() time.Time

	//Service layer defined getter setter attributes for building an event
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

//EventMetadata are optional identifiers appended to Events which
//may bear functionality downstream. For example, the trace ID can
//simplify tracking a request across numerous components (or services).
type EventMetadata interface {
	GetOriginStreamName() string
	GetTraceID() string

	SetOriginStreamName(string) EventMetadata
	SetTraceID(string) EventMetadata
}

//Subscriber applies side-effects per events.
type Subscriber interface {
	//Setup lifecycle method
	Start(EventBus)

	//Apply side-effects given an event
	Apply(context.Context, Event)
}

//EventBus centralizes access to an event source persistency layer.
//It's use should be disciplined.
type EventBus interface {
	//Subscribe to a stream name with a Subscriber
	Subscribe(string, Subscriber) error

	//Write an event to the event bus
	Write(string, Event) error

	//Read from a stream in the event bus starting at N position and consume J records.
	//A limit of -1 should return all events from N position till the end for a given stream.
	Read(string, int, int) ([]Event, error)

	//Conviniently create a new event that meets the Event interface
	NewEvent() Event

	//Conviniently create new event metadata that meets the EventMetadata interface
	NewEventMetadata() EventMetadata
}
