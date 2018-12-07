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
	"strings"
	"time"
	"travel-agent-backend/models"
	"travel-agent-backend/utils"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gmail "google.golang.org/api/gmail/v1"
)

//EmailReader reads unread emails
func EmailReader(configuration models.Config, db *sql.DB) {
	fmt.Println("Starting email reader forever service...")

	queryMap := make(map[string]string)
	programTime := time.Now()

	b, err := ioutil.ReadFile("/Users/kunal/go/src/travel-agent-backend/credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.GmailModifyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	client := getClient(config)
	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	user := configuration.Email.User
	r, err := srv.Users.Labels.List(user).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve labels: %v", err)
	}

	//Case not possible
	if len(r.Labels) == 0 {
		log.Println("No labels found.")
		return
	}
	// fmt.Println("Labels:")
	// for _, l := range r.Labels {
	// 	fmt.Println("- %s\n", l.Name, l.MessagesUnread, l.Id)
	// }

	queryMap["from"] = configuration.Email.EmailID
	queryMap["label"] = "unread"

	query := generateRequestQuery(queryMap)
	for {
		if time.Since(programTime).Seconds() >= configuration.GmailAPIInterval {
			programTime = time.Now()
			req := srv.Users.Messages.List(user).Q(query)

			res, err := req.Do()
			if err != nil {
				log.Fatalf("Error in fetching message list")
			}
			//loop over email messages
			for i := range res.Messages {
				id := res.Messages[i].Id

				singleMsg := srv.Users.Messages.Get(user, id).Format("full")

				message, err := singleMsg.Do()
				if err != nil {
					log.Fatalf("Err in getting msg, in do: %v", err)
				}

				//loop over email body parts
				for i := 1; i < len(message.Payload.Parts); i++ {
					fileName := message.Payload.Parts[i].Filename
					fmt.Println("Filename??", fileName)
					if fileName == "" {
						continue
					}
					fileType := strings.Split(fileName, ".")[1]
					if fileType != "pdf" {
						continue
					}
					fmt.Println("---", message.Payload.Parts[i])

					fmt.Println("---", user)
					fmt.Println("---", message.Id)
					fmt.Println("---", message.Payload.Parts[i].PartId)
					fmt.Println("---", message.Payload.Parts[i].Body.AttachmentId)
					fmt.Println("-------END------")

					//verify email address

					attachment, err := srv.Users.Messages.Attachments.Get(user, message.Id, message.Payload.Parts[i].Body.AttachmentId).Do()
					if err != nil {
						log.Fatalf("Error in finding attachment data do! %v", err)
					}
					// fmt.Println("Printing attachment: ", attachment)

					//verify attachment type

					decoded, err := base64.URLEncoding.DecodeString(attachment.Data)
					if err != nil {
						log.Fatalf("Error in decoding attachment data do! %v", err)
					}

					err = utils.WriteToFile("attachments/"+fileName, decoded)
					if err != nil {
						log.Fatalf("Error in writing attachment data to file %v", err)
					}
					//parse email attachment
				}
				err = markEmailAsRead(srv, configuration, id)
				if err != nil {
					log.Fatalf("Error in marking email as read %v", err)
				}
			}
		}
	}
}

//Generates a gmail query by which gmail list api is to be called
func generateRequestQuery(queryMap map[string]string) string {
	var queryString string
	for key, value := range queryMap {
		queryString += key + ":" + value + ","
	}
	queryString = queryString[:len(queryString)-1]
	return queryString
}

func markEmailAsRead(srv *gmail.Service, configuration models.Config, msgID string) error {
	var markMessageRead = new(gmail.ModifyMessageRequest)
	markMessageRead.RemoveLabelIds = []string{"UNREAD"}
	req := srv.Users.Messages.Modify(configuration.Email.User, msgID, markMessageRead)
	_, err := req.Do()
	return err
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
