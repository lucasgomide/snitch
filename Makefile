.PHONY: test-travis test build

test:
	@go test `go list ./... | grep -v /vendor/`

test-travis:
	./travis-test.sh

build:
	GOARCH="amd64" GOOS="linux" go build -o ./build/snitch_linux/snitch
	GOARCH="amd64" GOOS="darwin" go build -o ./build/snitch_darwin/snitch
