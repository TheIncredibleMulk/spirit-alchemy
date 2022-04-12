package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/docs/v1"
	"google.golang.org/api/option"
)

// If having trouble see the following link:
// https://developers.google.com/docs/api/quickstart/go

// accepts structs and arrrays and marshalls data into readable json so you can pass it into jq for easier groking
func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", " ")
	return string(s)
}

// Retrieves a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Requests a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		log.Fatalf("Unable to cache OAuth token: %v", err)
	}
	json.NewEncoder(f).Encode(token)
}

func main() {
	ctx := context.Background()
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/documents.readonly")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := docs.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Docs client: %v", err)
	}

	if err != nil {
		log.Fatalf("Unable to Retrieve Docs client: %v", err)
	}

	// Prints the title of the requested doc:
	// https://docs.google.com/document/d/195j9eDD3ccgjQRttHhJPymLJUCOUjs-jmwTrekvdjFE/edit
	// docId := "195j9eDD3ccgjQRttHhJPymLJUCOUjs-jmwTrekvdjFE"
	// https://docs.google.com/document/d/1x4VEgs0fa7igmr1H6xozBWodlm3h4SbEqlby_Lh2VVc/edit
	docId := "1x4VEgs0fa7igmr1H6xozBWodlm3h4SbEqlby_Lh2VVc"
	doc, err := srv.Documents.Get(docId).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from document: %v", err)
	}

	fmt.Printf("The title of the doc is: %s\n\n", doc.Title)

	var bookmarks []int

	fmt.Println("Len of Content: ", len(doc.Body.Content))
	for ic, c := range doc.Body.Content {
		if c != nil {
			if c.Paragraph != nil {
				if c.Paragraph.Elements != nil {
					for _, e := range c.Paragraph.Elements {
						if e != nil {
							if e.TextRun != nil {
								if e.TextRun.TextStyle.Link != nil {
									fmt.Println("BookmarkId: ", e.TextRun.TextStyle.Link.BookmarkId)
									bookmarks = append(bookmarks, ic)
								}
								fmt.Printf("Content: %+s", e.TextRun.Content)
								lastchar := e.TextRun.Content[len(e.TextRun.Content)-1]
								if lastchar != 0x000A {
									fmt.Println()
								}
								// fmt.Printf("Fontsize: %+v   ", e.TextRun.TextStyle.FontSize.Magnitude)

							}
						}
					}
				}
			}
		}
	}
	fmt.Println(bookmarks)

	// for _, c := range content {
	// 	for _, e := range c.Paragraph.Elements {
	// 		fmt.Printf("Elements content: %+v", prettyPrint(e.TextRun.Content))
	// 	}
	// }
	// fmt.Printf("%s\n", doc.Body.Content)
}
