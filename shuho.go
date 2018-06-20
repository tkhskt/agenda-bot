package main

import (
	"fmt"
	"github.com/nlopes/slack"
)

func remindShuho() {
	token, err := getToken()
	if err != nil {
		fmt.Errorf("shuho token error: %v", err)
	}
	client := slack.New(token.SlackToken)
	_, _, err = client.PostMessage(token.ShuhoChannel, "みんなもう週報出した？", slack.PostMessageParameters{})
	if err != nil {
		fmt.Errorf("shuho post error: %v", err)
	}
}
