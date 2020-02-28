package eventsourcinglocalsync

import (
	"strings"
	"testing"

	"github.com/alejandroq/goventsourcing/pkg/eventsourcingiface"
	"github.com/stretchr/testify/assert"
)

type context = eventsourcingiface.Context
type metadata = eventsourcingiface.MessageMetadata

type mock struct {
	ctx          context
	sideeffected bool
}

func (ms *mock) WithContext(ctx context) {
	ms.ctx = ctx
}

func (ms *mock) Apply(event event) {
	if strings.Contains(event.Type, "EventBusTested") {
		ms.sideeffected = true
	}
}

func TestEventBus_Subscribe(t *testing.T) {
	eb := New()
	sn := "EventBusTested"
	ms := mock{nil, false}
	err := eb.Subscribe(sn, &ms)
	assert.Nil(t, err)

	m := eventsourcingiface.Message{
		Type:     "EventBusTested:command-uniqueid",
		Metadata: metadata{},
		Body:     nil,
		Version:  0,
	}

	err = eb.Write(sn, m)
	assert.Nil(t, err)

	//this stream doesn't exist, therefore return empty
	es, err := eb.Read("ButtonClicked", 0, 2)
	assert.Nil(t, err)
	assert.Len(t, es, 0)

	//this stream returns less than the limit return 1
	es, err = eb.Read(sn, 0, 2)
	assert.Nil(t, err)
	assert.Len(t, es, 1)

	//the given position doesn't exist, therefore return empty
	_ = eb.Write(sn, m)
	es, err = eb.Read(sn, 3, 1)
	assert.Nil(t, err)
	assert.Len(t, es, 0)
}

func TestEventBus_ReadWrite(t *testing.T) {
	eb := New()
	sn := "EventBusTested"
	ms := mock{nil, false}
	err := eb.Subscribe(sn, &ms)
	assert.Nil(t, err)

	m := eventsourcingiface.Message{
		Type:     "EventBusTested:command-uniqueid",
		Metadata: metadata{},
		Body:     nil,
		Version:  0,
	}

	_ = eb.Write(sn, m)
	_ = eb.Write(sn, m)
	_ = eb.Write(sn, m)

	es, err := eb.Read(sn, 0, 2)
	assert.Nil(t, err)
	assert.Len(t, es, 2)

	//assert that globalsequenceid's increment as expected
	assert.Equal(t, es[0].GlobalSequenceID, 1)
	assert.Equal(t, es[1].GlobalSequenceID, 2)

	//assert the mock services sideeffect was called
	assert.True(t, ms.sideeffected)
}
