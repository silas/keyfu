setup:
	go get
	go get github.com/kr/godep
	go get github.com/jteeuwen/go-bindata/...
	godep restore

static.go:
	go-bindata -o=./static.go static

build: static.go
	go build

install: static.go
	go install

test:
	go test -v

cover:
	go test -coverprofile=.coverage.out
	go tool cover -html=.coverage.out

clean:
	rm -f keyfu static.go *.test

.PHONY: clean cover test
