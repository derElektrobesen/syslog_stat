package main

import (
	"os"
	"time"

	_ "github.com/derElektrobesen/syslog_stat/pkg/deps"
	"github.com/derElektrobesen/syslog_stat/pkg/storage"
	"github.com/derElektrobesen/syslog_stat/pkg/syslogsrv"
	"github.com/rs/zerolog"
)

func main() {
	lg := newLogger()

	s := syslogsrv.NewServer(lg, syslogsrv.Config{
		ListenAddr: "0.0.0.0:8888",
		NRoutines:  1,
	})

	st, err := storage.NewStorage(lg, storage.Config{
		TntAddr:             "127.0.0.1:1234",
		TntTimeout:          10 * time.Millisecond,
		TntReconnectTimeout: 200 * time.Millisecond,
		TntUser:             "user",
		TntPass:             "pass",
		TntFuncName:         "store",
	})

	if err != nil {
		lg.Error().Err(err).Msg("unable to connect to tarantool")
		return
	}

	s.Run(st.Store)
}

func newLogger() zerolog.Logger {
	return zerolog.New(os.Stderr).With().Timestamp().Logger()
}
