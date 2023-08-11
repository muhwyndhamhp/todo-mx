GO_MODULE := github.com/muhwyndhamhp/todo-mx
VERSION  ?= $(shell git describe --tags --abbrev=0)
LDFLAGS   := -X "$(GO_MODULE)/config.Version=$(VERSION)"


run:
	@air -c .air.toml --build.cmd "go build -ldflags \"$(LDFLAGS)\" -o ./tmp/main ."
