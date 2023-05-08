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
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
	twilio "github.com/twilio/twilio-go"

	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type TwilioMessage struct {
	Request   string `json:"request"`
	Response  string `json:"response"`
	EventName string `json:"event_name"`
}

func sendMessageToTwilio(event TwilioMessage) error {
	accountSid := os.Getenv("YOUR_ACCOUNT_SID")
	authToken := os.Getenv("YOUR_AUTH_TOKEN")
	from := os.Getenv("YOUR_TWILIO_PRODUCT") + ":" + os.Getenv("YOUR_TWILIO_NUMBER")
	toPhone := os.Getenv("YOUR_TWILIO_PRODUCT") + ":" + os.Getenv("RECIPIENT_NUMBER")
	message := event.EventName + " || " + event.Request + " || " + event.Response

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username:   accountSid,
		Password:   authToken,
		AccountSid: accountSid,
	})

	params := openapi.CreateMessageParams{}
	params.SetTo(toPhone)
	params.SetFrom(from)
	params.SetBody(message)

	_, err := client.Api.CreateMessage(&params)
	if err != nil {
		return fmt.Errorf("Request to Twilio returned an error: %v", err)
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

	var event TwilioMessage
	err := json.Unmarshal([]byte(body.Body), &event)
	message := "nok"
	if err == nil {
		err = sendMessageToTwilio(event)
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
