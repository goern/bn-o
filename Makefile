# This is just a wrapper for the old school guys

OUT_DIR=_output
OUT_PKG_DIR=Godeps/_workspace/pkg
BN_O_VERSION=0.1.0

.PHONY: all build
all build: main.go deps
	CGO_ENABLED=0 go build --ldflags '-extldflags "-static" -X github.com/goern/bn-o/version=$(BN_O_VERSION)'

deps: Godeps/Godeps.json
	godep restore

.PHONY: test
test:
	go test -v github.com/goern/bn-o

.PHONY: image
image: build test
	strip bn-o
	docker build --rm --tag goern/bn-o:$(BN_O_VERSION) -f Dockerfile .

.PHONY: clean
clean:
	rm -rf bn-o

.PHONY: clean-image
clean-image:
	docker rmi goern/bn-o:$(BN_O_VERSION)

.PHONY: tags
tags:
	gotags -tag-relative=true -R=true -sort=true -f="tags" -fields=+l .
