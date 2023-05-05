package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

var twilioWebhookURL = os.Getenv("TWILIO_WEBHOOK_URL")

// TwilioMessage struct
type TwilioMessage struct {
	EventName string `json:"event_name"`
	Request   string `json:"request"`
	Response  string `json:"response"`
}

func sendMessageToTwilio(twilioMessage TwilioMessage) error {
	jsonMessage, err := json.Marshal(twilioMessage)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", twilioWebhookURL, bytes.NewBuffer(jsonMessage))
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
		return fmt.Errorf("Request to Twilio returned an error: %v", resp.Status)
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
	lambda.Start(HandleLambdaEvent)
}
