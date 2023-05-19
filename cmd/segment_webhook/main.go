// server.go
//
// Use this sample code to handle webhook events in your integration.
//
// 1) Create a new Go module
//   go mod init example.com/liteAPI/webhooks/example
//
// 2) Paste this code into a new file (server.go)
//
// 3) Install dependencies
//   go get -u github.com/gin-gonic/gin
//
// 4) Run the server on http://127.0.0.1:8080
//   go run server.go

package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
	"github.com/segmentio/analytics-go/v3"
)

type SegmentMessage struct {
	EventName string `json:"event_name"`
	Request   string `json:"request"`
	Response  string `json:"response"`
}

func sendMessageToSegment(segmentMessage SegmentMessage) error {

	client := analytics.New(os.Getenv("SEGMENT_WRITE_KEY"))

	properties := analytics.Properties{}

	properties.Set("request", segmentMessage.Request)
	properties.Set("response", segmentMessage.Response)

	// Enqueues a track event that will be sent asynchronously.
	client.Enqueue(analytics.Track{
		UserId:     "Lite-Api",
		Event:      segmentMessage.EventName,
		Properties: properties,
	})
	// Flushes any queued messages and closes the client.
	err := client.Close()
	if err != nil {
		return err
	}

	return nil
}

// MyResponse struc
type MyResponse struct {
	Message string `json:"success"`
}

type functionBody struct {
	Body string `json:"body"`
}

func HandleLambdaEvent(body functionBody) (MyResponse, error) {

	var event SegmentMessage
	err := json.Unmarshal([]byte(body.Body), &event)

	message := "nok"
	if err == nil {
		err = sendMessageToSegment(event)
		if err == nil {
			message = "ok"
		}
	}
	return MyResponse{Message: message}, err
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(".env file couldn't be loaded")
	}

	lambda.Start(HandleLambdaEvent)
}
