.PHONY: test-travis test

test:
	@go test `go list ./... | grep -v /vendor/`

test-travis:
	@go test -coverprofile=hook.coverprofile ./hook
	@go test -coverprofile=tsuru.coverprofile ./tsuru
	@gover
	@goveralls -coverprofile=gover.coverprofile -service travis-ci
