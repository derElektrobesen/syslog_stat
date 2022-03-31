package types

import "time"

//easyjson:json
type LogContent struct {
	RemoteAddr    string `json:"remote_addr"` // user ip addr
	Request       string `json:"request"`
	HTTPStatus    int    `json:"status"`
	HTTPReferer   string `json:"http_referrer"`
	HTTPUserAgent string `json:"http_user_agent"`
}

//easyjson:json
type LogMessage struct {
	LogContent
	RemoteHost string // hostname which handled user request
	Timestamp  time.Time
}
