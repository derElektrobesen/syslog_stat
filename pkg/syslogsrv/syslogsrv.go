package syslogsrv

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"

	"github.com/derElektrobesen/syslog_stat/pkg/types"
	"gopkg.in/mcuadros/go-syslog.v2"
	"gopkg.in/mcuadros/go-syslog.v2/format"
)

type Config struct {
	ListenHostPort string
	NRoutines      int
}

type Server struct {
	cfg    Config
	logger zerolog.Logger
}

type Handler func(types.LogMessage)

func NewServer(lg zerolog.Logger, cfg Config) *Server {
	return &Server{
		cfg:    cfg,
		logger: lg,
	}
}

func (s *Server) Run(handler Handler) error {
	ch := make(syslog.LogPartsChannel)

	server := syslog.NewServer()
	server.SetFormat(syslog.RFC3164)
	server.SetHandler(syslog.NewChannelHandler(ch))

	if s.cfg.NRoutines < 1 {
		return fmt.Errorf("bad number of routines: %d", s.cfg.NRoutines)
	}

	if err := server.ListenUDP(s.cfg.ListenHostPort); err != nil {
		return fmt.Errorf("unable to listen %s: %w", s.cfg.ListenHostPort, err)
	}

	if err := server.Boot(); err != nil {
		return fmt.Errorf("unable to boot: %w", err)
	}

	for i := 0; i < s.cfg.NRoutines; i++ {
		go func() {
			for logParts := range ch {
				s.handleMessages(logParts, handler)
			}
		}()
	}

	s.logger.Info().Msg("starting the server")

	server.Wait()

	return nil
}

func (s *Server) handleMessages(logParts format.LogParts, handle Handler) {
	msg := types.LogMessage{
		RemoteHost: logParts["hostname"].(string),
		Timestamp:  logParts["timestamp"].(time.Time),
	}

	content := logParts["content"].(string)
	err := msg.LogContent.UnmarshalJSON([]byte(content))
	if err != nil {
		s.logger.Warn().
			Str("content", content).
			Err(err).
			Msg("unable to unmashal content")
	} else {
		handle(msg)
	}
}
