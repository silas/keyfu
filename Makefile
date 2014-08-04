build: static
	go build

static:
	go-bindata -o=./static.go src static

save:
	godep save -copy=false

release: static
	rm -fr ./build
	gox -osarch="darwin/amd64" -osarch="linux/amd64" -osarch="linux/386" -output="./build/keyfu_{{.OS}}_{{.Arch}}/keyfu"
	find ./build -type f -exec zip -j {}.zip {} \;
	find ./build -type f -name keyfu -delete
	find ./build -type d -mindepth 1 -exec mv {}/keyfu.zip {}.zip \; -delete

install: static
	go install

test:
	go test -v

cover:
	go test -coverprofile=.coverage.out
	go tool cover -html=.coverage.out

setup:
	go get
	go get github.com/mitchellh/gox
	go get github.com/tools/godep
	go get github.com/jteeuwen/go-bindata/...
	godep restore

clean:
	rm -f build keyfu static.go *.test

.PHONY: clean cover static test
