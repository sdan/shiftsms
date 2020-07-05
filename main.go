package main

import (
	"fmt"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func main() {
	config := oauth1.NewConfig("3cmfM8AoiNtkvU3oSc9eJ0eit", "RO4tuRpSeux5WpeivJwDrRHzkg7p26iXk2jpKrkk3F39tmNgTT")
	token := oauth1.NewToken("800117847774986240-KrZIBz46nJLvqj6gnahQrn0EegGwIL6", "KAA1uDLpKeaRKemUuq9oxSffhVkTkZmALvJ54w8Y6FjQc")
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	params := &twitter.StreamUserParams{
		StallWarnings: twitter.Bool(true),
	}
	stream, streamErr := client.Streams.User(params)
	if streamErr != nil {
		fmt.Println("we've got a problem")
	}
	demux := twitter.NewSwitchDemux()

	demux.DM = func(dm *twitter.DirectMessage) {
		fmt.Println(dm.SenderID)
	}

	for message := range stream.Messages {
		demux.Handle(message)
	}
}
