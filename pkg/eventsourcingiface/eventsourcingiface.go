package eventsourcingiface

import (
	"context"
	"time"
)

//MessageMetadata are optional identifiers appended to Events
//which occasions functionality downstream.
type MessageMetadata struct {
	OriginStreamName string `json:"originStreamName"`
	TraceID          string `json:"traceID"`
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
//feasible.
type Event struct {
	TransactionID    string    `json:"transactionID"`
	LocalSequenceID  int       `json:"sequenceID"`
	GlobalSequenceID int       `json:"globalSequenceID"`
	Timestamp        time.Time `json:"timestamp"`
	Message
}

//Message are written to a stream in an EventBus; during this process
//they are fleshed out into Events.
type Message struct {
	Type     string                 `json:"type"`
	Metadata MessageMetadata        `json:"metadata"`
	Body     map[string]interface{} `json:"body"`
	Version  int                    `json:"version"`
}

//Context is an EventBus context with ReadableWritable characteristics.
type Context interface {
	context.Context
	ReadableWritableEventBus
}

//Subscriber applies messages. Implementing domains, would
//benefit from an identity state (a zero value so to speak)
//and associative principles.
//Apply returns void as-is intended to introduce side-effects
//in the bounded context.
//StartWithContext is flexible and dependent upon the implementation
//of the triggering EventBus. It is generally reccomended that the
//Context be generated in the EventBus for resouce cleanup, etc.
type Subscriber interface {
	WithContext(Context)
	Apply(Event)
}

//EventBus centralizes access to an event source persistency layer
//and therefore should be disciplined as the abstraction for said.
//Subscriptions can derive from the EventBus.
//Commands are event categories or types, typically tagged with type
//strings that include the substring of `command:<event>` and are meant
//to enact side-effects.
type EventBus interface {
	SubscribableEventBus
	ReadableWritableEventBus
}

//SubscribableEventBus is an interface to subscribe to an EventBus
type SubscribableEventBus interface {
	//Subscribe to a stream name with a Subscriber
	Subscribe(string, Subscriber) error
}

//ReadableWritableEventBus is an interface to write and read from an EventBus.
type ReadableWritableEventBus interface {
	//Write a message to a stream
	Write(string, Message) error

	//Read from a stream starting at N position and consume J records
	Read(string, int, int) ([]Event, error)
}
