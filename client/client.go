package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/dghubble/oauth1"
)

//Struct to parse webhook load
type WebhookLoad struct {
	UserId             string                 `json:"for_user_id,omitempty"`
	DirectMessageEvent []interface{}          `json:"direct_message_events,omitempty"`
	Metadata           map[string]interface{} `json:"users,omitempty"`
}

//Struct to parse DM
type DirectMessage struct {
	Id          int
	Type        string
	RecipientID User
	SenderID    User
	Message     string
}

// type Metadata struct {
// 	User1 User
// 	User2 User
// }

//Struct to parse user
type User struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Handle string `json:"screen_name"`
}

func CreateClient() *http.Client {
	//Create oauth client with consumer keys and access token
	config := oauth1.NewConfig(os.Getenv("CONSUMER_KEY"), os.Getenv("CONSUMER_SECRET"))
	token := oauth1.NewToken(os.Getenv("ACCESS_TOKEN_KEY"), os.Getenv("ACCESS_TOKEN_SECRET"))

	return config.Client(oauth1.NoContext, token)
}

func RegisterWebhook() {
	fmt.Println("Registering webhook...")
	httpClient := CreateClient()

	//Set parameters
	path := "https://api.twitter.com/1.1/account_activity/all/" + os.Getenv("WEBHOOK_ENV") + "/webhooks.json"
	values := url.Values{}
	values.Set("url", os.Getenv("APP_URL")+"/webhook/twitter")

	//Make Oauth Post with parameters
	resp, _ := httpClient.PostForm(path, values)
	defer resp.Body.Close()
	//Parse response and check response
	body, _ := ioutil.ReadAll(resp.Body)
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		panic(err)
	}
	fmt.Println(data)
	fmt.Println("Webhook id of " + data["id"].(string) + " has been registered")
	SubscribeWebhook()
}

func SubscribeWebhook() {
	fmt.Println("Subscribing webapp...")
	client := CreateClient()
	path := "https://api.twitter.com/1.1/account_activity/all/" + os.Getenv("WEBHOOK_ENV") + "/subscriptions.json"
	resp, _ := client.PostForm(path, nil)
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	//If response code is 204 it was successful
	if resp.StatusCode == 204 {
		fmt.Println("Subscribed successfully")
	} else if resp.StatusCode != 204 {
		fmt.Println("Could not subscribe the webhook. Response below:")
		fmt.Println(string(body))
	}
}

// func SendTweet(tweet string, reply_id string) (*Tweet, error) {
// 	fmt.Println("Sending tweet as reply to " + reply_id)
// 	//Initialize tweet object to store response in
// 	var responseTweet Tweet
// 	//Add params
// 	params := url.Values{}
// 	params.Set("status", tweet)
// 	params.Set("in_reply_to_status_id", reply_id)
// 	//Grab client and post
// 	client := CreateClient()
// 	resp, err := client.PostForm("https://api.twitter.com/1.1/statuses/update.json", params)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()
// 	//Decode response and send out
// 	body, _ := ioutil.ReadAll(resp.Body)
// 	fmt.Println(string(body))
// 	err = json.Unmarshal(body, &responseTweet)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &responseTweet, nil
// }
