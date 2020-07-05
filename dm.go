package main

import (
	"fmt"
	"log"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func main() {
	config := oauth1.NewConfig("3cmfM8AoiNtkvU3oSc9eJ0eit", "RO4tuRpSeux5WpeivJwDrRHzkg7p26iXk2jpKrkk3F39tmNgTT")
	token := oauth1.NewToken("800117847774986240-KrZIBz46nJLvqj6gnahQrn0EegGwIL6", "KAA1uDLpKeaRKemUuq9oxSffhVkTkZmALvJ54w8Y6FjQc")

	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)

	// List most recent 10 Direct Messages
	messages, _, err := client.DirectMessages.EventsList(
		&twitter.DirectMessageEventsListParams{Count: 20},
	)
	fmt.Println("User's DIRECT MESSAGES:")
	if err != nil {
		log.Fatal(err)
	}
	for _, event := range messages.Events {
		fmt.Printf("%+v\n", event)
		fmt.Printf("  %+v\n", event.Message)
		fmt.Printf("  %+v\n", event.Message.Data)
	}

	fmt.Println("DONE")

	// Show Direct Message event
	// event, _, err := client.DirectMessages.EventsShow("1066903366071017476", nil)
	// fmt.Printf("DM Events Show:\n%+v, %v\n", event.Message.Data, err)

	// Create Direct Message event
	/*
		event, _, err = client.DirectMessages.EventsNew(&twitter.DirectMessageEventsNewParams{
			Event: &twitter.DirectMessageEvent{
				Type: "message_create",
				Message: &twitter.DirectMessageEventMessage{
					Target: &twitter.DirectMessageTarget{
						RecipientID: "2856535627",
					},
					Data: &twitter.DirectMessageData{
						Text: "testing",
					},
				},
			},
		})
		fmt.Printf("DM Event New:\n%+v, %v\n", event, err)
	*/

	// Destroy Direct Message event
	//_, err = client.DirectMessages.EventsDestroy("1066904217049133060")
	//fmt.Printf("DM Events Delete:\n err: %v\n", err)
}
