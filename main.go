package main

import (
	"encoding/json"
	"log"
	"time"

	_ "github.com/derElektrobesen/syslog_stat/pkg/deps"
	"github.com/derElektrobesen/syslog_stat/pkg/types"
	"gopkg.in/mcuadros/go-syslog.v2"
)

func main() {
	channel := make(syslog.LogPartsChannel)
	handler := syslog.NewChannelHandler(channel)

	server := syslog.NewServer()
	server.SetFormat(syslog.RFC3164)
	server.SetHandler(handler)

	listenHostPort := "0.0.0.0:8888"

	if err := server.ListenUDP(listenHostPort); err != nil {
		log.Println("Unable to listen %s: %s", listenHostPort, err)
		return
	}

	if err := server.Boot(); err != nil {
		log.Println("Unable to boot: %s", err)
		return
	}

	go func(channel syslog.LogPartsChannel) {
		for logParts := range channel {
			msg := types.LogMessage{
				RemoteHost: logParts["hostname"].(string),
				Timestamp:  logParts["timestamp"].(time.Time),
			}

			content := logParts["content"].(string)
			err := json.Unmarshal([]byte(content), &msg.LogContent)
			if err != nil {
				log.Println("unable to unmarshal content %s: %s", content, err)
				continue
			}

			handleLogMessage(msg)
		}
	}(channel)

	log.Println("Starting server")

	server.Wait()
}

func handleLogMessage(msg types.LogMessage) {
	log.Printf("%+v", msg)
}
