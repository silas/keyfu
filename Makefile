setup:
	go get
	go get github.com/mitchellh/gox
	go get github.com/tools/godep
	go get github.com/jteeuwen/go-bindata/...
	godep restore

static.go:
	go-bindata -o=./static.go static

build: static.go
	go build

release: static.go
	rm -fr ./build
	mkdir -p ./build
	gox -osarch="darwin/amd64" -osarch="linux/amd64" -osarch="linux/386" -output="./build/keyfu_{{.OS}}_{{.Arch}}"
	find ./build -type f -exec zip {}.zip {} \;

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
