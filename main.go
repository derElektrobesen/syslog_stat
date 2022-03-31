package main

import (
	"log"
	"os"

	_ "github.com/derElektrobesen/syslog_stat/pkg/deps"
	"github.com/derElektrobesen/syslog_stat/pkg/syslogsrv"
	"github.com/derElektrobesen/syslog_stat/pkg/types"
	"github.com/rs/zerolog"
)

func main() {
	lg := newLogger()

	s := syslogsrv.NewServer(lg, syslogsrv.Config{
		ListenHostPort: "0.0.0.0:8888",
		NRoutines:      1,
	})

	s.Run(handleLogMessage)
}

func newLogger() zerolog.Logger {
	return zerolog.New(os.Stderr).With().Timestamp().Logger()
}

func handleLogMessage(msg types.LogMessage) {
	log.Printf("%+v", msg)
}
