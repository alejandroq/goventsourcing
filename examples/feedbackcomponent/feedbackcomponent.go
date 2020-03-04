package feedbackcomponent

import (
	"context"
	"fmt"

	"github.com/alejandroq/goventsourcing/pkg/eventsourcingiface"
)

//Service ...
type Service struct {
	bus eventsourcingiface.EventBus
}

//Feedback ...
type Feedback struct {
	Contents string `json:"contents"`
}

//Start ...
func (s *Service) Start(bus eventsourcingiface.EventBus) {
	fmt.Println("[INFO] started feedback component service")
	s.bus = bus
}

//Apply ...
func (s *Service) Apply(ctx context.Context, event eventsourcingiface.Event) {
	fmt.Println("[DEBUG]", event)
	if event == nil {
		return
	}
	fmt.Printf("[INFO] applying change from %s sequence number %v.\n", event.GetMetadata().GetOriginStreamName(), event.GetLocalSequenceID())
}
