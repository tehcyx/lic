default: build

LIC_VERSION=$(shell cat VERSION)

build: test cover
	go build -i -o bin/lic -ldflags "-X pkg/lic/cmd.Version=${LIC_VERSION}" ./cmd/lic

test:
	go test ./...

cover:
	go test ./... -cover

clean:
	rm -rf bin