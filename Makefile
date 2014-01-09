STATIC=$(shell find static -type f -not -name '*.go' -exec echo '{}.go' \;)

build: $(STATIC)
	sed -i.bak 's|package main|package static|g' static/*.go
	sed -i.bak 's|go_bindata|Data|g' static/*.go
	find static -type f -name '*.bak' -delete
	go build -tags='static' -v

%.go: %
	go-bindata -toc $<

test:
	go test -v

clean:
	rm static/*.go
