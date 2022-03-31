.MAIN: build

EASYJSON_FILES := $(shell grep -rl --include="*.go" 'easyjson:json' pkg/)

.PHONY: build
build: bin/syslog_stat

clean:
	@find . -name '*_easyjson.go' -delete -print

vendor:
	go mod vendor

generate: bin/easyjson $(EASYJSON_FILES:.go=_easyjson.go)

bin/syslog_stat: generate
	go build -o syslog_stat main.go

bin/easyjson:
	go build -o bin/easyjson github.com/mailru/easyjson/easyjson

%_easyjson.go: %.go
	bin/easyjson $<
