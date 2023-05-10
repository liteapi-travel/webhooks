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
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
)

type IftttMessage struct {
	EventName string `json:"event_name"`
	Request   string `json:"request"`
	Response  string `json:"response"`
}

type PayloadMessage struct {
	EventName string `json:"value1"`
	Request   string `json:"value2"`
	Response  string `json:"value3"`
}

func sendMessageToIfttt(iftttMessage IftttMessage) error {

	payloadMessage := PayloadMessage(iftttMessage)

	jsonMessage, err := json.Marshal(payloadMessage)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", os.Getenv("IFTTT_WEBHOOK_URL"), bytes.NewBuffer(jsonMessage))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("Request to Slack returned an error: %v", resp.Status)
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

	var event IftttMessage
	err := json.Unmarshal([]byte(body.Body), &event)
	message := "nok"
	if err == nil {
		err = sendMessageToIfttt(event)
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
