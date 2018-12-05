package services

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"travel-agent-backend/models"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gmail "google.golang.org/api/gmail/v1"
)

func Start(configuration models.Config, db *sql.DB) {
	fmt.Println("Starting forever service...")

	b, err := ioutil.ReadFile("C:/Users/Alisha Virani/go/src/travel-agent-backend/credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	user := "me"
	r, err := srv.Users.Labels.List(user).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve labels: %v", err)
	}
	fmt.Println("listing lists ==>", r)
	if len(r.Labels) == 0 {
		log.Println("No labels found.")
		return
	}
	// fmt.Println("Labels:")
	// for _, l := range r.Labels {
	// 	fmt.Println("- %s\n", l.Name, l.MessagesUnread, l.Id)
	// }

	req := srv.Users.Messages.List(user).Q("from:av.alisha99@gmail.com,label:unread")

	res, err := req.Do()
	if err != nil {
		fmt.Println("Error in fetching message list")
		return
	}

	// fmt.Println("printing unread messages!!!", *res)

	for i := range res.Messages {
		id := res.Messages[i].Id
		// fmt.Println("Printing id", id)

		singleMsg := srv.Users.Messages.Get(user, id).Format("full")

		message, err := singleMsg.Do()
		if err != nil {
			log.Panic("Err in getting msg, in do: ", err)
		}
		// fmt.Println("Printing single msg", message)

		// fmt.Println("Printting message payload parts: ", message.Payload.Parts)

		//i <
		for i := 1; i < len(message.Payload.Parts); i++ {
			fmt.Println("---", message.Payload.Parts[i])

			fmt.Println("---", user)
			fmt.Println("---", message.Id)
			fmt.Println("---", message.Payload.Parts[i].PartId)
			fmt.Println("---", message.Payload.Parts[i].Body.AttachmentId)
			fmt.Println("-------END------")

			attachment, err := srv.Users.Messages.Attachments.Get(user, message.Id, message.Payload.Parts[i].Body.AttachmentId).Do()
			if err != nil {
				log.Panic("Error in finding attachment data do!", err)
			}
			// fmt.Println("Printing attachment: ", attachment)

			decoded, err := base64.URLEncoding.DecodeString(attachment.Data)
			if err != nil {
				log.Panic("Error in decoding attachment data do!", err)
			}

			f, err := os.Create(message.Payload.Parts[i].Filename)
			if err != nil {
				panic(err)
			}
			defer f.Close()

			if _, err := f.Write(decoded); err != nil {
				panic(err)
			}
			if err := f.Sync(); err != nil {
				panic(err)
			}
			//parse email attachment, check headers

		}
	}
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
