package main

import (
	"encoding/json"
	"net/http"

	"github.com/alejandroq/goventsourcing/examples/eventsourcinglocalsync"
	"github.com/alejandroq/goventsourcing/examples/feedbackcomponent"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func main() {
	sn := "PublishedFeedback"
	eb := eventsourcinglocalsync.New()
	fs := feedbackcomponent.Service{}
	_ = eb.Subscribe(sn, &fs)

	r := gin.Default()
	r.POST("/", func(c *gin.Context) {
		var f feedbackcomponent.Feedback
		c.Bind(&f)

		//create a new event with
		b, _ := json.Marshal(f)
		e := eventsourcinglocalsync.New().NewEvent().SetType("SentFeedback")
		e = e.SetEventID(uuid.New().String()).SetBody(string(b))
		_ = eb.Write(sn, e)

		c.JSON(http.StatusCreated, gin.H{
			"status": "created",
			"data":   f,
		})
	})
	_ = r.Run()
}
