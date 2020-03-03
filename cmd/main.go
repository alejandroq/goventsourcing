package main

import (
	"github.com/alejandroq/goventsourcing/example/eventsourcinglocalsync"
	"github.com/alejandroq/goventsourcing/example/feedbackcomponent"
)

func main() {
	eb := eventsourcinglocalsync.New()
	fs := feedbackcomponent.Service{}
	_ = eb.Subscribe("Published Feedback", &fs)
}
