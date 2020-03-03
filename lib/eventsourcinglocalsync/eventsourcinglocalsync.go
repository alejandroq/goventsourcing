package eventsourcinglocalsync

import (
	"fmt"
	"time"

	"github.com/alejandroq/goventsourcing/pkg/eventsourcingiface"
	"github.com/google/uuid"
)

//EventBus subscribes, reads and writes to the event store.
type EventBus struct {
	events              map[string][]eventsourcingiface.Event
	subscribers         map[string][]eventsourcingiface.Subscriber
	localSequenceCount  map[string]int
	globalSequenceCount int
}

//New eventsourcingiface.eventbusfactory
func New() eventsourcingiface.EventBusFactory {
	return &EventBus{
		make(map[string][]eventsourcingiface.Event),
		make(map[string][]eventsourcingiface.Subscriber),
		make(map[string]int),
		0,
	}
}

//Subscribe subscribes a subscriber to a stream
func (eb *EventBus) Subscribe(sn string, s eventsourcingiface.Subscriber) error {
	eb.subscribers[sn] = append(eb.subscribers[sn], s)
	return nil
}

//Write an event to 2 streams - 1 to the plain stream name and
//the other is written to a sub-stream: <stream name>-<event id>
func (eb *EventBus) Write(sn string, m eventsourcingiface.Event) error {
	tid := uuid.New().String()

	//increment global sequence IDs
	eb.globalSequenceCount = eb.globalSequenceCount + 1

	write := func(sn string) {
		eb.localSequenceCount[sn] = eb.localSequenceCount[sn] + 1

		//create event and begin copying types, bodies, etc from client provided events
		var e eventsourcingiface.Event = event{
			TransactionID:    tid,
			LocalSequenceID:  eb.localSequenceCount[sn],
			GlobalSequenceID: eb.globalSequenceCount,
			Timestamp:        time.Now(),
		}
		e = e.SetEventID(m.GetEventID())
		e = e.SetType(m.GetType()).SetMetadata(m.GetMetadata())
		e = e.SetBody(m.GetBody()).SetVersion(m.GetVersion())

		//append event to stream
		eb.events[sn] = append(eb.events[sn], e.(event))

		//apply new event to subscribers of a stream name
		for _, s := range eb.subscribers[sn] {
			s.Apply(e)
		}
	}

	//write to the global stream name as well as the event ID based stream
	write(sn)
	fmt.Println("[INFO] writing to stream " + fmt.Sprintf("%s-%s", sn, m.GetEventID()))
	write(fmt.Sprintf("%s-%s", sn, m.GetEventID()))

	return nil
}

//Read events from a stream. Projections can be computed across state if the limit is -1.
func (eb *EventBus) Read(sn string, pos int, limit int) (rs []eventsourcingiface.Event, err error) {
	rs = []eventsourcingiface.Event{}

	defer func() {
		//recover from out of bound errors
		if r := recover(); r != nil {
			//return an empty result set if out of bounds
			rs = []eventsourcingiface.Event{}
			return
		}
	}()

	es := eb.events[sn]

	if len(es) > limit && limit != -1 {
		for i := 0; i < limit; i++ {
			rs = append(rs, es[pos+i])
		}
	} else {
		rs = es
	}

	fmt.Println(eb.events)

	return rs, nil
}

//NewEvent creates a new event
func (eb *EventBus) NewEvent() eventsourcingiface.Event {
	return event{}
}

//NewEventMetadata creates new metadata
func (eb *EventBus) NewEventMetadata() eventsourcingiface.EventMetadata {
	return metadata{}
}
