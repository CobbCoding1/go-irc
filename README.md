# Go-irc
simple irc client, written in go, specifically for use for my Twitch bot.

## Getting started

### Download the package
```bash
git clone git@github.com:CobbCoding1/go-irc.git $GOPATH/src/github.com/CobbCoding1/goirc
```

### import into project
```Go
import (
    "github.com/CobbCoding1/goirc"
)
```


### Usage
```
client := goirc.Init(domain_name //irc.twitch.tv for example, port, password, username, nickname)
client.Connect() //connect to the specified domain
client.Disconnect() //handles disconnect, should be called immediately after connect
client.Join(server //cobbcoding for example)

// main loop
for {
    msg := client.GetData() //get any messages being sent
    // do stuff with message
}
```
