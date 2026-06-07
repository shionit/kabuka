test:
	go test -v ./...

test-integration:
	go test -v -tags integration ./...

test-coverage:
	go test -v -coverprofile=./gotest.profile ./... 2>&1
	go tool cover -html=./gotest.profile -o ./coverage.html
