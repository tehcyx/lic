default: build

MAJOR_VERSION=0
MINOR_VERSION=1
PATCH_VERSION=0



build: test cover
	go build -i -o bin/lic -ldflags "-X pkg/lic/cmd.Version=${MAJOR_VERSION}.${MINOR_VERSION}.${PATCH_VERSION}" ./cmd/lic

test:
	go test ./...

cover:
	go test ./... -cover

clean:
	rm -rf bin