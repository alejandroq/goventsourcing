package eventsourcinglocalasync

import (
	"context"
	"fmt"
	"sync"
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
	ctx                 context.Context
	events              map[stream][]event
	channels            map[stream][]*chan event // doesn't this have issues as not actor? so multi-subscribers would lead to issue? fan out.
	localSequenceCount  map[stream]int
	globalSequenceCount int
	channelsMux         sync.Mutex
	writeMux            sync.Mutex
}

// TODO partial streams and global streams? - additional channels? yes per cadence "string-x". how to enforce? automatically derived unique stream name as a derivative of a global? how?

//New concrete implementation of EventBus
func New(ctx context.Context) eventsourcingiface.EventBus {
	return &EventBus{
		ctx,
		make(map[stream][]event),
		make(map[stream][]*chan event),
		make(map[stream]int),
		0,
		sync.Mutex{},
		sync.Mutex{},
	}
}

//Subscribe subscribes a subscriber to a stream
func (eb *EventBus) Subscribe(stream string, s subscriber) error {
	eb.channelsMux.Lock()
	err := eb.subscribe(stream, s)
	eb.channelsMux.Unlock()
	return err
}

func (eb *EventBus) subscribe(stream string, s subscriber) error {
	c := make(chan event)
	eb.channels[stream] = append(eb.channels[stream], &c)
	go func(c <-chan event, s subscriber) {
		fmt.Println("[INFO] this routine is running")
		for {
			select {
			case e := <-c:
				fmt.Println("[INFO]", e)
				s.Apply(e)
			case <-eb.ctx.Done():
				fmt.Println("[INFO] this routine is clossing")
				return
			}
		}
	}(c, s)
	return nil
}

//Write an Event to a stream
func (eb *EventBus) Write(stream string, message message) error {
	eb.writeMux.Lock()
	e, err := eb.write(stream, message)
	eb.writeMux.Unlock()
	eb.writeToChannels(stream, e)
	return err
}

func (eb *EventBus) write(stream string, message message) (*event, error) {
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
	return &e, nil
}

func (eb *EventBus) writeToChannels(stream string, event *event) {
	eb.channelsMux.Lock()
	for _, c := range eb.channels[stream] {
		*c <- *event
	}
	eb.channelsMux.Unlock()
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
