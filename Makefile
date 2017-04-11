.PHONY: test-travis test

test:
	@go test `go list ./... | grep -v /vendor/`

test-travis:
	@go test -coverprofile=hook.coverprofile ./hook
	@go test -coverprofile=tsuru.coverprofile ./tsuru
	@gover
	@goveralls -coverprofile=gover.coverprofile -service travis-ci

build:
	GOARCH="amd64" GOOS="linux" go build -o ./build/snitch_linux/snitch
	GOARCH="amd64" GOOS="darwin" go build -o ./build/snitch_darwin/snitch
