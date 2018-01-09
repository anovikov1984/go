all: install-deps-tests run-tests

run-tests: install-deps-tests
	@echo "Running tests"
	bash ./scripts/run-tests.sh

install-deps-tests:
	@echo "Install dependencies for tests"
	# go get -u github.com/pubnub/go
	go get -u golang.org/x/net/context
	
	go get github.com/satori/go.uuid
	# cd ${GOPATH}/src/github.com/satori/go.uuid\
	
	# git --git-dir ${GOPATH}/src/github.com/satori/go.uuid/.git stash 
	# git --git-dir ${GOPATH}/src/github.com/satori/go.uuid/.git checkout tags/v1.1.0

	# cd ${GOPATH}/src/github.com/pubnub	
	go get -u github.com/stretchr/testify


