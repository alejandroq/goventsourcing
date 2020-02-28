package eventsourcinglocalsync

import (
	"time"

	"github.com/alejandroq/goventsourcing/pkg/eventsourcingiface"
)

type stream = string
type event = eventsourcingiface.Event
type message = eventsourcingiface.Message
type subscriber = eventsourcingiface.Subscriber

//EventBus implementing an eventsourcingiface.EventBus.
//The local implementation does not require a polling or
//a polling relationship to the persistency layer as it will
//encapsulated in this version.
type EventBus struct {
	events              map[stream][]event
	subscribers         map[stream][]subscriber
	localSequenceCount  map[stream]int
	globalSequenceCount int
}

//New concrete implementation of EventBus
func New() eventsourcingiface.EventBus {
	return &EventBus{
		make(map[stream][]event),
		make(map[stream][]subscriber),
		make(map[stream]int),
		0,
	}
}

//Subscribe subscribes a subscriber to a stream
func (eb *EventBus) Subscribe(stream string, subscriber subscriber) error {
	eb.subscribers[stream] = append(eb.subscribers[stream], subscriber)
	return nil
}

//Write an Event to a stream
func (eb *EventBus) Write(stream string, message message) error {
	eb.globalSequenceCount = eb.globalSequenceCount + 1
	eb.localSequenceCount[stream] = eb.localSequenceCount[stream] + 1
	e := event{
		TransactionID:    "uniqueid",
		LocalSequenceID:  eb.localSequenceCount[stream],
		GlobalSequenceID: eb.globalSequenceCount,
		Timestamp:        time.Now(),
		Message:          message,
	}
	eb.events[stream] = append(eb.events[stream], e)
	eb.writeToEachSubscriber(stream, e)
	return nil
}

func (eb EventBus) writeToEachSubscriber(stream string, event event) {
	for _, s := range eb.subscribers[stream] {
		s.Apply(event)
	}
}

//Read Events from a stream
func (eb *EventBus) Read(stream string, position int, limit int) (rs []event, err error) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()
	es := eb.events[stream]
	if len(es) == 0 || limit == 0 {
		rs = []event{}
	} else if len(es) > limit {
		for i := 0; i < limit; i++ {
			rs = append(rs, es[position+i])
		}
	} else {
		rs = es
	}
	return rs, nil
}

//Subscription is an implementation of eventsourcingiface.Subscription
type Subscription struct {
	eb *EventBus
}
