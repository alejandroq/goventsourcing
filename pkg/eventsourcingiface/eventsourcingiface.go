package eventsourcingiface

import (
	"context"
	"time"
)

//Context is an EventBus context with ReadableWritable characteristics.
type Context interface {
	context.Context
	EventBusFactory
}

//Event message appends to a global event bus for a bounded context
//for auditability purposes and to engage clients such as simulations
//and subscriptions.
//Contracts with implicit clients should be upheld while it makes
//sense to do so. As this is typically an archectectural pain point for
//event sourcing, strategies such as anti-corruption layers can assist
//in the event of needing to maintain older contracts with dependent
//downstream clients.
//A possible solution for the above, between commands and other event types,
//a nameing convention should include a version string for explicit dependency;
//therefore extensible anti-corruption layers can translate n+1 versions to n if
//feasible. The interface is meant to provide
//struct tag flexibility captured in concrete implementations.
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
//orchestration vs choreography (?)
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
//in the bounded context.
//StartStartWith is flexible and dependent upon the implementation
//of the triggering EventBus. It is generally reccomended that the
//Context be generated in the EventBus for resouce cleanup, etc.
type Subscriber interface {
	StartWith(Context)
	Apply(Event)
}

//EventBus centralizes access to an event source persistency layer
//and therefore should be disciplined as the abstraction for said.
//Subscriptions can derive from the EventBus.
//Commands are event categories or types, typically tagged with type
//strings that include the substring of `command:<event>` and are meant
//to enact side-effects.
type EventBus interface {
	//Subscribe to a stream name with a Subscriber
	Subscribe(string, Subscriber) error

	//Write a message to the EventBus
	Write(string, Event) error

	//Read from the EventBus starting at N position and consume J records
	Read(string, int, int) ([]Event, error)
}

//EventBusFactory encapsulates EventBus and a factory methods for
//for generating interface compliant Event and EventMetadata.
type EventBusFactory interface {
	NewEvent() Event
	NewEventMetadata() EventMetadata
	EventBus
}
