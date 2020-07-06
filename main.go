package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	client "github.com/sdan/shiftsms/client"
	tw "github.com/sfreiberg/gotwilio"
)

// var twilio *tw.Twilio
// var lastUserHandle *string
// var lastUserName *string
// var lastUserID *string

func main() {
	zerr := os.Unsetenv("lastUserName")
	aerr := os.Unsetenv("lastUserID")
	berr := os.Unsetenv("lastUserHandle")
	if zerr != nil {
		fmt.Println(zerr)
	}
	if aerr != nil {
		fmt.Println(aerr)
	}
	if berr != nil {
		fmt.Println(berr)
	}

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
	m.HandleFunc("/webhook/twilio", SMSHandler).Methods("POST")

	//Start Server
	server := &http.Server{
		Handler: m,
	}
	server.Addr = ":9090"
	server.ListenAndServe()
}

// type SMSMessage struct {
// 	ToCountry string
// }

type lastUserState struct {
	lastUserHandle string
	lastUserName   string
	lastUserID     string
}

//Twilio Consumer
func SMSHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("SMS Recieved")
	request.ParseForm()
	// for k, v := range request.Form {
	// 	fmt.Println("key:", k)
	// 	fmt.Println("value:", v)
	// }
	fmt.Println("msg:", request.Form["Body"][0])
	go SendDM(request.Form["Body"][0])
	// body, _ := ioutil.ReadAll(request.Body)

	// var load SMSMessage
	// err := json.Unmarshal(body, &load)
	// if err != nil {
	// 	fmt.Println("An error occured: " + err.Error())
	// }

	// // fmt.Println(request.Body)
	// // var insms string
	// // err := json.Unmarshal(body, &insms)
	// // if err != nil {
	// // 	fmt.Println("An error occured: " + err.Error())
	// // }
	// fmt.Println(load)

}

// type DMStuct struct {
// 	event SendEvent
// }

// type SendEvent struct {
// 	typeMessage    string
// 	message_create SendMessageCreate
// }

// type SendMessageCreate struct {
// 	target       map[string]interface{}
// 	message_data map[string]interface{}
// }

func SendDM(text string) {
	path := "https://api.twitter.com/1.1/direct_messages/events/new.json"
	httpClient := client.CreateClient()
	// {"event": {"type": "message_create", "message_create": {"target": {"recipient_id": "RECIPIENT_USER_ID"}, "message_data": {"text": "Hello World!"}}}}
	// var sendthis DMStuct
	lastUserID := os.Getenv("lastUserID")

	payloadtext := `{"event": {"type": "message_create", "message_create": {"target": {"recipient_id": "RECIPIENT_USER_ID"}, "message_data": {"text": "INPUTTEXTHERE"}}}}`
	payloadtext = strings.Replace(payloadtext, "RECIPIENT_USER_ID", lastUserID, 1)
	payloadtext = strings.Replace(payloadtext, "INPUTTEXTHERE", text, 1)
	var jsonStr = []byte(payloadtext)
	// sendthis := DMStuct

	// jsonValue, _ := json.Marshal(sendthis)

	resp, err := httpClient.Post(path, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		fmt.Println("oopsies", err)
	}
	fmt.Println("respoe:", resp)
}

//Twitter Consumer
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
		dmtext := load.DirectMessageEvent[0].(map[string]interface{})["message_create"].(map[string]interface{})["message_data"].(map[string]interface{})["text"]
		tostring := load.DirectMessageEvent[0].(map[string]interface{})["message_create"].(map[string]interface{})["target"].(map[string]interface{})["recipient_id"].(string)
		fromstring := load.DirectMessageEvent[0].(map[string]interface{})["message_create"].(map[string]interface{})["sender_id"].(string)
		tostringName := load.Metadata[fromstring].(map[string]interface{})["name"].(string)
		tostringScreenName := load.Metadata[fromstring].(map[string]interface{})["screen_name"].(string)
		// fromstring := load.DirectMessageEvent[0].(map[string]interface{})["message_create"].(map[string]interface{})["sender_id"].(string)
		if load.UserId == tostring {
			fmt.Println(tostringName+" (@"+tostringScreenName+")", "sent you: ", dmtext)
			// *lastUserName = tostringName
			// *lastUserID = tostring
			// *lastUserHandle = tostringScreenName

			// databasePass = os.Getenv("DATABASE_PASS")
			// fmt.Printf("Database Password: %s\n", databasePass)

			err := os.Setenv("lastUserName", tostringName)
			aerr := os.Setenv("lastUserID", fromstring)
			berr := os.Setenv("lastUserHandle", tostringScreenName)
			if err != nil {
				fmt.Println(err)
			}
			if aerr != nil {
				fmt.Println(aerr)
			}
			if berr != nil {
				fmt.Println(berr)
			}

			// 16502756958
			accountSid := os.Getenv("TWILIO_SID")
			authToken := os.Getenv("TWILIO_AUTH")
			twilio := tw.NewTwilioClient(accountSid, authToken)
			from := "+"
			to := "+"
			fmt.Println("sent sms", dmtext.(string))
			twilio.SendSMS(from, to, tostringName+"(@"+tostringScreenName+"): "+dmtext.(string), "", "")
		} else {
			fmt.Println("user sent via sms")
		}
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
