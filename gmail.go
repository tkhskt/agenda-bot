package main

import (
	"fmt"
	"time"
	"log"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

func main() {

	config := oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		Endpoint:     google.Endpoint,
		RedirectURL:  "urn:ietf:wg:oauth:2.0:oob",          //今回はリダイレクトしないためこれ
		Scopes:       []string{"https://mail.google.com/"}, //必要なスコープを追加
	}

	expiry, _ := time.Parse("2006-01-02", "2018-06-16")
	token := oauth2.Token{
		AccessToken:  "",
		TokenType:    "Bearer",
		RefreshToken: "",
		Expiry:       expiry,
	}

	client := config.Client(oauth2.NoContext, &token)

	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve gmail Client %v", err)
	}

	ums := gmail.NewUsersMessagesService(srv)

	r, err := srv.Users.Messages.List("me").Do()
	if err != nil {
		log.Fatalf("Unable to get labels. %v", err)
	}

	for _, v := range r.Messages {
		ms, err := ums.Get("me", v.Id).Do()
		if err != nil {
			log.Fatalf("%v", err)
		}
		fmt.Println(ms.Snippet)
		fmt.Println(ms.Payload.Filename)
	}

}
