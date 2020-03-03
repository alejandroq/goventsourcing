package feedback

import (
	"fmt"

	"github.com/alejandroq/goventsourcing/pkg/eventsourcingiface"
)

//Service ...
type Service struct {
	ctx eventsourcingiface.Context
}

//StartWith ...
func (s *Service) StartWith(ctx eventsourcingiface.Context) {
	fmt.Println("[INFO] started feedback service")
	s.ctx = ctx
}

//Apply ...
func (s *Service) Apply(event eventsourcingiface.Event) {
	fmt.Printf("[INFO] applying change from %s sequence number %v.\n", event.GetMetadata().GetOriginStreamName(), event.GetLocalSequenceID())
}
