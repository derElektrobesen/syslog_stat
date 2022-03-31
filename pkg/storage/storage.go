package storage

import (
	"fmt"
	"time"

	"github.com/tarantool/go-tarantool"

	"github.com/derElektrobesen/syslog_stat/pkg/types"
	"github.com/rs/zerolog"
)

type Config struct {
	TntAddr             string
	TntTimeout          time.Duration
	TntReconnectTimeout time.Duration
	TntUser             string
	TntPass             string

	TntFuncName string
}

type Storage struct {
	cfg    Config
	logger zerolog.Logger
	tnt    *tarantool.Connection
}

func NewStorage(lg zerolog.Logger, cfg Config) (*Storage, error) {
	tntOpts := tarantool.Opts{
		Timeout:       cfg.TntTimeout,
		Reconnect:     cfg.TntReconnectTimeout,
		MaxReconnects: 0,
		User:          cfg.TntUser,
		Pass:          cfg.TntPass,
	}

	client, err := tarantool.Connect(cfg.TntAddr, tntOpts)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to %s: %w", cfg.TntAddr, err)
	}

	return &Storage{
		cfg:    cfg,
		logger: lg,
		tnt:    client,
	}, nil
}

func (s *Storage) Store(msg types.LogMessage) {
	_, err := s.tnt.Call(
		s.cfg.TntFuncName,
		[]interface{}{
			msg.RemoteAddr,
			msg.RemoteHost,
			msg.Timestamp.Unix(),
		},
	)

	if err != nil {
		// TODO: limit number of log lines
		s.logger.Warn().Err(err).Msg("unable to send stat in tnt")
	}
}
