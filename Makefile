all: build

GIT_COMMIT:=$(shell git rev-list -1 HEAD)
GIT_LAST_TAG:=$(shell git describe --abbrev=0 --tags)
GIT_EXACT_TAG:=$(shell git name-rev --name-only --tags HEAD)
UNAME_S := $(shell uname -s)
XARGS:=xargs -r
ifeq ($(UNAME_S),Darwin)
    XARGS:=xargs
endif

COMMANDS_PATH:=github.com/sladyn98/hawk-eye/commands
LDFLAGS:=-X ${COMMANDS_PATH}.GitCommit=${GIT_COMMIT} \
	-X ${COMMANDS_PATH}.GitLastTag=${GIT_LAST_TAG} \
	-X ${COMMANDS_PATH}.GitExactTag=${GIT_EXACT_TAG}

build:
	go generate
	go build -ldflags "$(LDFLAGS)" .

# produce a build debugger friendly
debug-build:
	go generate
	go build -ldflags "$(LDFLAGS)" -gcflags=all="-N -l" .

install:
	go generate
	go install -ldflags "$(LDFLAGS)" .

releases:
	go generate
	gox -ldflags "$(LDFLAGS)" -output "dist/{{.Dir}}_{{.OS}}_{{.Arch}}"

test:
	go test -v -bench=. ./...
