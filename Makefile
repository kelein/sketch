Branch?=`git rev-parse --abbrev-ref HEAD`
SHA1=`git rev-parse --short HEAD`
Date=`date +%Y-%m-%dT%H:%M:%S`
Version=$(Branch)@$(SHA1)@$(Date)
User=`whoami`@`hostname`

LDFLAGS=-ldflags "-X 'sketch/pkg/version.Version=${Version}'  \
	-X 'sketch/pkg/version.Revision=${SHA1}' 	   			 \
	-X 'sketch/pkg/version.Branch=${Branch}' 	   			 \
	-X 'sketch/pkg/version.BuildDate=${Date}' 	   			 \
	-X 'sketch/pkg/version.BuildUser=${User}'"

default: build

app:
	@echo "\033[32mcurrent build flag version: ${Version}\033[0m"
	@go build ${LDFLAGS} -o bin/sketch ./cmd/sketch/sketch.go
	@echo "\033[32mbinary file output target at: bin/sketch\033[0m"

app-linux:
	@echo "\033[32mcurrent build flag version: ${Version}\033[0m"
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o bin/sketch ./cmd/sketch/sketch.go
	@echo "\033[32mbinary file output target at: bin/sketch\033[0m"

build: app

run:
	@go run cmd/sketch/sketch.go

image:
	@make app-linux
	@docker build . -t ${TAG}
	@docker push ${TAG}

tools:
	go get -v github.com/Masterminds/glide
	go get -v honnef.co/go/tools/cmd/staticcheck
	go get -v honnef.co/go/tools/cmd/gosimple
	go get -v honnef.co/go/tools/cmd/unused
	go get -v github.com/gordonklaus/ineffassign
	go get -v github.com/fzipp/gocyclo
	go get -v github.com/golang/lint/golint
	go get -v github.com/pierrre/gotestcover
	go get -v github.com/client9/misspell/cmd/misspell
	go get -v github.com/mfridman/tparse

vet:
	@if test -n '$(shell go vet `glide nv` 2>&1 | egrep -v "^vendor\/|.*_easyjson\.go|.*\.pb\.go|^exit\ status")'; then \
		echo '$(shell go vet `glide nv` 2>&1 | egrep -v "^vendor\/|.*_easyjson\.go|.*\.pb\.go|^exit\ status")'; \
		exit 1; \
	fi

lint:
	@if test -n '$(shell golint `glide nv` 2>&1 | egrep -v ".*_easyjson\.go|.*\.pb\.go|.*lock_linux\.go")'; then \
		echo '$(shell golint `glide nv` 2>&1 | egrep -v ".*_easyjson\.go|.*\.pb\.go|.*lock_linux\.go")'; \
	fi

staticcheck:
	staticcheck $(shell glide nv)

gosimple:
	gosimple $(shell glide nv)

unused:
	unused $(shell glide nv)

ineffassign:
	@for f in `find . -type d -depth 1 |egrep -v "git|hook|vendor|easyjson\.go"`; do \
		ineffassign $$f; \
	done

misspell:
	misspell $(shell find . -maxdepth 1 -mindepth 1 -type d |egrep -v "^\.$|git|hook|vendor|admin")
	misspell $(shell find ./admin -type f -not -path "./admin/web/*")

gocyclo:
	@gocyclo -over 20 $(shell find . -name "*.go" |egrep -v "pb\.go|_test\.go|vendor|easyjson\.go")

check: vet lint staticcheck gosimple unused misspell

doc:
	godoc -http=:6060

cover:
	go test -v -coverprofile=cover.out ./...
	go tool cover -html=cover.out
	@rm -f cover.out

fmt:
	gofmt -s -w .

clean:
	@rm -f ./bin/*
	@echo "\033[32mclean done\033[0m"

.PHONY: vet clean linux doc fmt build test cover check unused gocyclo gosimple staticcheck rpm upload
