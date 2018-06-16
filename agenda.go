package main

import (
	"fmt"
	"time"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"io/ioutil"
	"encoding/json"
	"encoding/base64"
	"strings"
	"errors"
	"os"
	"github.com/nlopes/slack"
	"github.com/jasonlvhit/gocron"
)

type token struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	SlackToken   string `json:"slackToken"`
	SlackChannel string `json:"slackChannel"`
}

type file struct {
	filename string
	dec      []byte
}

func decode(str string) ([]byte, error) {
	if strings.ContainsAny(str, "+/") {
		return nil, errors.New("invalid base64url encoding")
	}
	str = strings.Replace(str, "-", "+", -1)
	str = strings.Replace(str, "_", "/", -1)
	for len(str)%4 != 0 {
		str += "="
	}
	return base64.StdEncoding.DecodeString(str)
}

func getToken() (*token, error) {
	t, err := ioutil.ReadFile("./token.json")
	if err != nil {
		return nil, err
	}
	var tk *token
	json.Unmarshal(t, &tk)

	return tk, nil

}

func (tk *token) getGmailService() (*gmail.Service, error) {
	config := oauth2.Config{
		ClientID:     tk.ClientId,
		ClientSecret: tk.ClientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  "urn:ietf:wg:oauth:2.0:oob",          //今回はリダイレクトしないためこれ
		Scopes:       []string{"https://mail.google.com/"}, //必要なスコープを追加
	}

	expiry, _ := time.Parse("2006-01-02", "2018-06-16")
	token := oauth2.Token{
		AccessToken:  tk.AccessToken,
		TokenType:    "Bearer",
		RefreshToken: tk.RefreshToken,
		Expiry:       expiry,
	}

	client := config.Client(oauth2.NoContext, &token)

	srv, err := gmail.New(client)
	if err != nil {
		return nil, err
	}
	return srv, nil
}

func (fl *file) createFile() error {
	filePath := "./file/" + fl.filename
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write(fl.dec); err != nil {
		return err
	}
	if err := f.Sync(); err != nil {
		return err
	}
	return nil
}

func getFileFromMessage(part *gmail.MessagePart, srv *gmail.Service, ms *gmail.Message) ([]byte, error) {
	umas := gmail.NewUsersMessagesAttachmentsService(srv)
	atc := umas.Get("me", ms.Id, part.Body.AttachmentId)
	dt, err := atc.Do()
	if err != nil {
		return nil, err
	}
	dec, err := decode(dt.Data)
	if err != nil {
		return nil, err
	}
	return dec, nil
}

func (fl *file) saveFileName() error {
	file, err := os.Create(`filename.txt`)
	if err != nil {
		return err
	}
	defer file.Close()
	file.Write(([]byte)(fl.filename))
	return nil
}

func (fl *file) isLatestFile() (bool, error) {
	b, err := ioutil.ReadFile("filename.txt") // just pass the file name
	if err != nil {
		return false, err
	}
	if string(b) != fl.filename {
		return true, nil
	}

	return false, nil
}

func searchFileFromMessage(part []*gmail.MessagePart, ms *gmail.Message, srv *gmail.Service) (*file, error) {
	for _, v := range part {
		if strings.Contains(v.Filename, "アジェンダ") {
			dec, err := getFileFromMessage(v, srv, ms)
			if err != nil {
				return nil, err
			}
			fl := &file{v.Filename, dec}
			return fl, nil
		}
	}
	return nil, nil
}

func (fl *file) handleFile(token *token) error {
	fileIsLatest, err := fl.isLatestFile()
	if err != nil {
		fmt.Errorf("file latest error: %v", err)
	}
	if fileIsLatest {
		err = fl.saveFileName()
		if err != nil {
			fmt.Errorf("save file name error: %v", err)
		}
		err = fl.createFile()
		if err != nil {
			fmt.Errorf("create file error: %v", err)
		}
		err = fl.postSlack(token)
		if err != nil {
			fmt.Errorf("slack error: %v", err)
		}
	}
	return nil
}

func (fl *file) postSlack(token *token) error {
	api := slack.New(token.SlackToken)
	file := slack.FileUploadParameters{
		File:     "./file/" + fl.filename,
		Filetype: "pdf",
		Title:    fl.filename,
		Channels: []string{token.SlackChannel},
	}
	_, err := api.UploadFile(file)
	if err != nil {
		return err
	}
	return nil
}

func task() {
	defer fmt.Println("done task")
	token, err := getToken()
	if err != nil {
		fmt.Errorf("token error: %v", err)
	}
	srv, err := token.getGmailService()
	if err != nil {
		fmt.Errorf("getService error: %v", err)
	}

	ums := gmail.NewUsersMessagesService(srv)

	r, err := srv.Users.Messages.List("me").Do()
	if err != nil {
		fmt.Errorf("get message list error: %v", err)
	}

	for _, v := range r.Messages {
		ms, err := ums.Get("me", v.Id).Do()
		if err != nil {
			fmt.Errorf("%v", err)
		}
		fl, err := searchFileFromMessage(ms.Payload.Parts, ms, srv)
		if err != nil {
			fmt.Errorf("search file error: %v", err)
		}
		if fl != nil {
			err = fl.handleFile(token)
			if err != nil {
				fmt.Errorf("handle file error: %v", err)
			}
			return
		}
	}
}

func main() {
	gocron.Every(10).Minutes().Do(task)
	<-gocron.Start()
}
