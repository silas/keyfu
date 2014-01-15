STATIC=$(shell find static -type f -not -name '*.go' -exec echo '{}.go' \;)
OPTS="-tags='static' -v"

static: $(STATIC)
	sed -i.bak 's|package main|package static|g' static/*.go
	sed -i.bak 's|go_bindata|Data|g' static/*.go
	find static -type f -name '*.bak' -delete

build: static
	go build $(OPTS)

install: static
	go install $(OPTS)

%.go: %
	go-bindata -toc $<

test:
	go test -v

cover:
	go test -coverprofile=.coverage.out
	go tool cover -html=.coverage.out

clean:
	rm -f keyfu static/*.go *.test

.PHONY: clean cover test
