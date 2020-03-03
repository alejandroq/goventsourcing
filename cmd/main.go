package main

import (
	"github.com/alejandroq/goventsourcing/lib/eventsourcinglocalsync"
	"github.com/alejandroq/goventsourcing/lib/feedback"
)

func main() {
	eb := eventsourcinglocalsync.New()
	fs := feedback.Service{}
	_ = eb.Subscribe("Published Feedback", &fs)
}
