package eventsourcinglocalsync

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/alejandroq/goventsourcing/pkg/eventsourcingiface"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

//mock component service
type mock struct {
	sideeffected bool
	lastmessage  string
}

type body struct {
	Hello string `json:"hello"`
}

func (ms *mock) StartWith(bus eventsourcingiface.EventBus) {
	ms.lastmessage = "ENABLED"
}

func (ms *mock) Apply(ctx context.Context, event eventsourcingiface.Event) {
	//apply should be idempotent as message duplicity cannot be guaranteed nil
	if strings.Contains(event.GetType(), "example:command") {
		ms.sideeffected = true
	}
	var b body
	_ = json.Unmarshal([]byte(event.GetBody()), &b)
	ms.lastmessage = b.Hello
}

func TestSubscriberLifeCycle(t *testing.T) {
	eb := New()
	sn := "PublishedOrder"
	ms := mock{false, ""}
	_ = eb.Subscribe(sn, &ms)
	assert.Equal(t, "ENABLED", ms.lastmessage)
}

func TestWriteReadEvent(t *testing.T) {
	eb := New()
	sn := "PublishedOrder"
	ms := mock{false, ""}
	_ = eb.Subscribe(sn, &ms)

	b, _ := json.Marshal(body{"World"})
	e := eb.NewEvent().SetType("example:command").SetBody(string(b))
	_ = eb.Write(sn, e)

	//a write to "PublishedOrder" should cause a service side-effect.
	assert.True(t, ms.sideeffected)

	//read from "PublishedOrder" and verify that expected records exist.
	rs, _ := eb.Read(sn, 0, 5)
	assert.Greater(t, len(rs), 0)
	assert.Equal(t, "World", ms.lastmessage)
}

func TestIncrementingLocalGlobalSequenceIDs(t *testing.T) {
	eb := New()
	sn := "PublishedOrder"
	ms := mock{false, ""}
	_ = eb.Subscribe(sn, &ms)
	id := uuid.New().String()

	for i := 0; i < 3; i++ {
		b, _ := json.Marshal(body{"World"})
		e := eb.NewEvent().SetEventID(fmt.Sprintf("%s-%v", id, i)).SetType("example:command").SetBody(string(b))
		_ = eb.Write(sn, e)
	}

	rs, _ := eb.Read(sn, 0, -1)
	assert.Equal(t, 3, len(rs))

	assert.Less(t, rs[0].GetGlobalSequenceID(), rs[1].GetGlobalSequenceID())
	assert.Less(t, rs[0].GetLocalSequenceID(), rs[1].GetLocalSequenceID())

	//read from a sub-stream and assert the records exist as expected
	rs, _ = eb.Read(fmt.Sprintf("%s-%s-1", sn, id), 0, 5)
	assert.Equal(t, 1, len(rs))
}

func TestReadLimit(t *testing.T) {
	eb := New()
	sn := "PublishedOrder"
	ms := mock{false, ""}
	_ = eb.Subscribe(sn, &ms)
	b, _ := json.Marshal(body{"World"})
	e := eb.NewEvent().SetEventID(uuid.New().String()).SetBody(string(b))

	//write 3 times redundantly.
	_ = eb.Write(sn, e)
	_ = eb.Write(sn, e)
	_ = eb.Write(sn, e)

	//assert can only read 2 provided explicit limit
	rs, _ := eb.Read(sn, 0, 2)
	assert.Equal(t, 2, len(rs))

	//assert can read all provided an explicit limit of -1
	rs, _ = eb.Read(sn, 0, -1)
	assert.Equal(t, 3, len(rs))
}

func TestOneToManySubscribersToStream(t *testing.T) {
	eb := New()
	sn := "PublishedOrder"

	sub1 := mock{false, ""}
	sub2 := mock{false, ""}
	sub3 := mock{false, ""}

	//subscribe 3 subscribers
	_ = eb.Subscribe(sn, &sub1)
	_ = eb.Subscribe(sn, &sub2)
	_ = eb.Subscribe(sn, &sub3)

	//publish a faux event
	b, _ := json.Marshal(body{"World"})
	e := eb.NewEvent().SetEventID(uuid.New().String()).SetType("example:command").SetBody(string(b))
	eb.Write(sn, e)

	//assert on known sideeffects
	assert.Equal(t, true, sub1.sideeffected)
	assert.Equal(t, true, sub2.sideeffected)
	assert.Equal(t, true, sub3.sideeffected)
}

func TestOriginStreamTraceIDIntegrity(t *testing.T) {
	eb := New()
	sn := "PublishedOrder"
	ms := mock{false, ""}
	_ = eb.Subscribe(sn, &ms)

	tid := uuid.New().String()
	b, _ := json.Marshal(body{"World"})
	m := eb.NewEventMetadata().SetOriginStreamName(sn).SetTraceID(tid)
	e := eb.NewEvent().SetEventID(uuid.New().String()).SetBody(string(b))
	e = e.SetMetadata(m)
	eb.Write(sn, e)

	rs, _ := eb.Read(sn, 0, -1)
	assert.Equal(t, tid, rs[0].GetMetadata().GetTraceID())
}
