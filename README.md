# agenda-bot

メールで共有されたアジェンダをslackのagendaチャンネルに流すbot

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
  "slackChannel": "slackのChannelId"
}
```

2. `filename.txt`の作成

`touch filename.txt`

