package main

import (
	"github.com/alejandroq/goventsourcing/examples/eventsourcinglocalsync"
	"github.com/alejandroq/goventsourcing/examples/feedbackcomponent"
)

func main() {
	eb := eventsourcinglocalsync.New()
	fs := feedbackcomponent.Service{}
	_ = eb.Subscribe("Published Feedback", &fs)
}
