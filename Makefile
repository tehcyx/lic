default: build

LIC_VERSION=$(shell cat VERSION)
GIT_COMMIT=$(shell git rev-list -1 HEAD)

build: test cover
	go build -ldflags "-X github.com/tehcyx/lic/pkg/lic/cmd.Version=${LIC_VERSION} -X github.com/tehcyx/lic/pkg/lic/cmd.GitCommit=${GIT_COMMIT}" -i -o bin/lic ./cmd/lic

install: build
	go install

test:
	go test ./...

cover:
	go test ./... -cover

clean:
	rm -rf bin