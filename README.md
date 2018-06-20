# agenda-bot

メールで共有されたアジェンダをslackのagendaチャンネルに流すbot

ついでに週報の催促もする

1分おきにgmailのinboxをチェックしてagendaが上がってないかチェックします

## Install

`go get github.com/gericass/agenda-bot`

## Setup

1. `token.json`の作成

```
{
  "clientId": "google APIのclientId",
  "clientSecret": "google APIのclientSecret",
  "accessToken": "google APIのaccessToken",
  "refreshToken": "google APIのrefreshToken",
  "slackToken": "slackのtoken",
  "agendaChannel": "slackのagendaを上げるchannelId"
  "shuhoChannel": "slackの週報リマインドをするchannelId"
}
```

2. `filename.txt`の作成

`touch filename.txt`

## Usage

`go run agenda.go`

