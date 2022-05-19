package main

import (
	"context"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/docs/v1"
	"google.golang.org/api/option"
)

type Recipe struct {
	Title   string
	Content string
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

	var bookmarks = make(map[int]string)
	var content string

	fmt.Println("Len of Content: ", len(doc.Body.Content))
	for ic, c := range doc.Body.Content {
		if c != nil {
			if c.Paragraph != nil {
				if c.Paragraph.Elements != nil {
					for _, e := range c.Paragraph.Elements {
						if e != nil {
							if e.TextRun != nil {
								if e.TextRun.TextStyle.Link != nil {
									// fmt.Println("BookmarkId: ", e.TextRun.TextStyle.Link.BookmarkId)
									// bookmarks = append(bookmarks, ic)
									bookmarks[ic] = e.TextRun.TextStyle.Link.BookmarkId
								}
								// fmt.Printf("Content: %+s", e.TextRun.Content)
								content = fmt.Sprintf("%+s%+s", content, e.TextRun.Content)
								// lastchar := e.TextRun.Content[len(e.TextRun.Content)-1]
								// if lastchar != 0x000A {
								// 	fmt.Println()
								// }
								// fmt.Printf("Fontsize: %+v   ", e.TextRun.TextStyle.FontSize.Magnitude)

							}
						}
					}
				}
			}
		}
	}
	fmt.Println(content)
	fmt.Println(bookmarks)
	fmt.Println("//==========================================================//", "//==========================================================//")
	fmt.Println()

	// for _, c := range content {
	// 	for _, e := range c.Paragraph.Elements {
	// 		fmt.Printf("Elements content: %+v", prettyPrint(e.TextRun.Content))
	// 	}
	// }
	// fmt.Printf("%s\n", doc.Body.Content)

	//==========================================================// html template parsing logic //==========================================================//

	f, err := os.Create("testTemplate.html")
	if err != nil {
		panic(err)
	}

	data := Recipe{
		Title:   doc.Title,
		Content: content,
	}

	tmpl := template.Must(template.ParseFiles("template/testTemplate.gohtml"))
	tmpl.Execute(f, data)
	fmt.Println(tmpl)

}
