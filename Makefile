test-coverage:
	go test -v -coverprofile=./gotest.profile ./... 2>&1
	go tool cover -html=./gotest.profile -o ./coverage.html

