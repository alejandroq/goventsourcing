package main

import (
	"github.com/alejandroq/goventsourcing/lib/eventsourcinglocalsync"
	"github.com/alejandroq/goventsourcing/lib/feedback"
)

func main() {
	eb := eventsourcinglocalsync.New()
	fb := feedback.Feedback{}
	eb.Subscribe("FeedbackSubmitted", &fb)
}
