test-coverage:
	go test -v -coverprofile=./coverage/gotest.profile ./... 2>&1 | tee ./coverage/gotest.log
	go tool cover -html=./coverage/gotest.profile -o ./coverage/index.html

