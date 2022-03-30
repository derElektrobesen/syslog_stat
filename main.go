package main

import (
	"encoding/json"
	"log"
	"time"

	"gopkg.in/mcuadros/go-syslog.v2"
)

type logContent struct {
	RemoteAddr    string `json:"remote_addr"`
	Request       string `json:"request"`
	HTTPStatus    int    `json:"status"`
	HTTPReferer   string `json:"http_referrer"`
	HTTPUserAgent string `json:"http_user_agent"`
}

type logMessage struct {
	logContent
	RemoteHost string
	Timestamp  time.Time
}

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
			msg := logMessage{
				RemoteHost: logParts["hostname"].(string),
				Timestamp:  logParts["timestamp"].(time.Time),
			}

			content := logParts["content"].(string)
			err := json.Unmarshal([]byte(content), &msg.logContent)
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

func handleLogMessage(msg logMessage) {
	log.Printf("%+v", msg)
}
