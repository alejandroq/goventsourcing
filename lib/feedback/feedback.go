package feedback

import (
	"fmt"

	"github.com/alejandroq/goventsourcing/pkg/eventsourcingiface"
)

type context = eventsourcingiface.Context
type eventbus = eventsourcingiface.EventBus
type event = eventsourcingiface.Event

//Feedback domain async handles given user feedback
type Feedback struct {
	ctx context
}

//WithContext is the initial step in the service lifecycle
//for context allocation and log aggregation.
func (f *Feedback) WithContext(ctx context) error {
	fmt.Println("[INFO] Staring Feedback")
	f.ctx = ctx
	return nil
}

//Apply events for the Feedback service's handling
func (f *Feedback) Apply(event event) {}
