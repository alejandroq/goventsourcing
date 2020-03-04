package eventsourcinglocalsync

import (
	"context"
	"fmt"
	"time"

	"github.com/alejandroq/goventsourcing/pkg/eventsourcingiface"
	"github.com/google/uuid"
)

//EventBus implements eventsourcingiface.EventBus
type EventBus struct {
	ctx                 context.Context
	events              map[string][]eventsourcingiface.Event
	subscribers         map[string][]eventsourcingiface.Subscriber
	localSequenceCount  map[string]int
	globalSequenceCount int
}

//New eventsourcingiface.EventBus implementation
func New() eventsourcingiface.EventBus {
	return &EventBus{
		context.Background(),
		make(map[string][]eventsourcingiface.Event),
		make(map[string][]eventsourcingiface.Subscriber),
		make(map[string]int),
		0,
	}
}

//Subscribe subscribes a subscriber to a stream
func (eb *EventBus) Subscribe(sn string, s eventsourcingiface.Subscriber) error {
	eb.subscribers[sn] = append(eb.subscribers[sn], s)
	s.Start(eb)
	return nil
}

//Write an event to 2 streams:
//- 1 to the general stream name
//- 2 to the sub-stream: ie. <stream name>-<event id>
func (eb *EventBus) Write(sn string, m eventsourcingiface.Event) error {
	//generate a unique transaction ID
	tid := uuid.New().String()

	//increment global sequence IDs
	eb.globalSequenceCount = eb.globalSequenceCount + 1

	write := func(sn string) {
		//increment local sequence IDs per the particular stream
		eb.localSequenceCount[sn] = eb.localSequenceCount[sn] + 1

		//generate an event and apply client provided parameters.
		var e eventsourcingiface.Event = event{
			TransactionID:    tid,
			LocalSequenceID:  eb.localSequenceCount[sn],
			GlobalSequenceID: eb.globalSequenceCount,
			Timestamp:        time.Now(),
		}

		e = e.SetEventID(m.GetEventID()).SetType(m.GetType())
		e = e.SetBody(m.GetBody()).SetVersion(m.GetVersion())
		if m.GetMetadata() != nil {
			e = e.SetMetadata(m.GetMetadata())
		} else {
			e = e.SetMetadata(eb.NewEventMetadata())
		}

		//append event to the ledger
		eb.events[sn] = append(eb.events[sn], e.(event))

		//apply new event to subscribed subscribers
		for _, s := range eb.subscribers[sn] {
			s.Apply(eb.ctx, e)
		}
	}

	//write to the general stream
	write(sn)
	//write to the sub-stream: ie. <stream name>-<event id>
	write(fmt.Sprintf("%s-%s", sn, m.GetEventID()))

	return nil
}

//Read events from a stream.
//Limits of -1 will return all records following a given position. This could be used for calculating
//aggregations and projections to dispose upon materialized data views (or ephemeral data stores).
func (eb *EventBus) Read(sn string, pos int, limit int) (rs []eventsourcingiface.Event, err error) {
	defer func() {
		//recover from out of bound errors
		if r := recover(); r != nil {
			//return an empty result set if requiring recovery
			rs = []eventsourcingiface.Event{}
			return
		}
	}()

	es := eb.events[sn]

	if limit == -1 && pos > 0 {
		return es[pos:], nil
	}

	if limit == -1 && pos <= 0 {
		return es, nil
	}

	if len(es[pos:]) > limit {
		return es[pos : pos+limit], nil
	}

	return es[pos:], nil
}

//NewEvent creates a new eventsourcingiface.Event
func (eb *EventBus) NewEvent() eventsourcingiface.Event {
	return event{}
}

//NewEventMetadata creates a new eventsourcingiface.EventMetadata
func (eb *EventBus) NewEventMetadata() eventsourcingiface.EventMetadata {
	return metadata{}
}
