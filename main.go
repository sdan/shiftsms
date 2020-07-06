package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	client "github.com/sdan/shiftsms/client"
	tw "github.com/sfreiberg/gotwilio"
)

// var twilio *tw.Twilio

func main() {
	//Load env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		fmt.Println("Error loading .env file")
	}

	fmt.Println("Starting Server")

	if args := os.Args; len(args) > 1 && args[1] == "-register" {
		go client.RegisterWebhook()
	} else {
		fmt.Println("No registration")
	}

	//Create a new Mux Handler
	m := mux.NewRouter()
	//Listen to the base url and send a response
	m.HandleFunc("/", func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(200)
		fmt.Fprintf(writer, "Server is up and running")
	})
	//Listen to crc check and handle
	m.HandleFunc("/webhook/twitter", CrcCheck).Methods("GET")
	m.HandleFunc("/webhook/twitter", WebhookHandler).Methods("POST")

	//Start Server
	server := &http.Server{
		Handler: m,
	}
	server.Addr = ":9090"
	server.ListenAndServe()
}

func WebhookHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Handler called")
	//Read the body of the tweet
	body, _ := ioutil.ReadAll(request.Body)
	// fmt.Println(string(body))
	//Initialize a webhok load obhject for json decoding
	var load client.WebhookLoad
	err := json.Unmarshal(body, &load)
	if err != nil {
		fmt.Println("An error occured: " + err.Error())
	}

	if load.DirectMessageEvent == nil {
		fmt.Println("Not a DM")
	} else {

		// fmt.Println("LOAD: ", load)
		// fmt.Println("UserID: ", load.UserId)
		// fmt.Println("DM: ", load.DirectMessageEvent[0].(map[string]interface{})["message_create"].(map[string]interface{})["message_data"].(map[string]interface{})["text"])
		// fmt.Println("To: ", load.DirectMessageEvent[0].(map[string]interface{})["message_create"].(map[string]interface{})["target"].(map[string]interface{})["recipient_id"])
		// fmt.Println("From: ", load.DirectMessageEvent[0].(map[string]interface{})["message_create"].(map[string]interface{})["sender_id"])
		// fmt.Println("MD: ", load.Metadata)
		// for k, v := range load.Metadata {
		// 	fmt.Println("key", k)
		// 	fmt.Println("val", v)
		// }

		// tostring := load.DirectMessageEvent[0].(map[string]interface{})["message_create"].(map[string]interface{})["target"].(map[string]interface{})["recipient_id"].(string)
		// fromstring := load.DirectMessageEvent[0].(map[string]interface{})["message_create"].(map[string]interface{})["sender_id"].(string)

		// // if typeErr != nil {
		// // 	fmt.Error(typeErr)
		// // }
		// fmt.Println("\n\n\n\n")
		// fmt.Println(load.Metadata[fromint].(map[string]interface{})["name"])
		// fmt.Println(load.Metadata[toint].(map[string]interface{})["name"])

		dmtext := load.DirectMessageEvent[0].(map[string]interface{})["message_create"].(map[string]interface{})["message_data"].(map[string]interface{})["text"]
		tostring := load.DirectMessageEvent[0].(map[string]interface{})["message_create"].(map[string]interface{})["target"].(map[string]interface{})["recipient_id"].(string)
		// fromstring := load.DirectMessageEvent[0].(map[string]interface{})["message_create"].(map[string]interface{})["sender_id"].(string)
		if load.UserId == tostring {
			fmt.Println(dmtext)
		}
		// 16502756958
		accountSid := os.Getenv("TWILIO_SID")
		authToken := os.Getenv("TWILIO_AUTH")
		twilio := tw.NewTwilioClient(accountSid, authToken)
		from := "+16502756958"
		to := "+14087469016"
		fmt.Println("sent sms", dmtext.(string))
		twilio.SendSMS(from, to, dmtext.(string), "", "")

		// fmt.Println("User1: ", load.Metadata.User1.Handle)
		// fmt.Println("User2: ", load.Metadata.User2.Handle)
		// // var user1 client.User
		// // var user2 client.User
		// fmt.Println("md int: ", load.Metadata.(map[int]interface{}))

		// //Check if it was a tweet_create_event and tweet was in the payload and it was not tweeted by the bot
		// if len(load.TweetCreateEvent) < 1 || load.UserId == load.TweetCreateEvent[0].User.IdStr {
		// 	return
		// }
		// //Send Hello world as a reply to the tweet, replies need to begin with the handles
		// //of accounts they are replying to
		// _, err = client.SendTweet("@"+load.TweetCreateEvent[0].User.Handle+" Hello World", load.TweetCreateEvent[0].IdStr)
		// if err != nil {
		// 	fmt.Println("An error occured:")
		// 	fmt.Println(err.Error())
		// } else {
		// 	fmt.Println("Tweet sent successfully")
		// }
	}
}

func CrcCheck(writer http.ResponseWriter, request *http.Request) {
	//Set response header to json type
	writer.Header().Set("Content-Type", "application/json")
	//Get crc token in parameter
	token := request.URL.Query()["crc_token"]
	if len(token) < 1 {
		fmt.Fprintf(writer, "No crc_token given")
		return
	}

	//Encrypt and encode in base 64 then return
	h := hmac.New(sha256.New, []byte(os.Getenv("CONSUMER_SECRET")))
	h.Write([]byte(token[0]))
	encoded := base64.StdEncoding.EncodeToString(h.Sum(nil))
	//Generate response string map
	response := make(map[string]string)
	response["response_token"] = "sha256=" + encoded
	//Turn response map to json and send it to the writer
	responseJson, _ := json.Marshal(response)
	fmt.Fprintf(writer, string(responseJson))
}
