docker_tag=latest
skip_pull=false
lint_yaml=false

clean:
	rm -f coverage.out
	rm -f coverage.html
	rm -f test-report.json
	rm -rf bin
	go clean -testcache
check:
ifeq ($(lint_yaml), true)
	yamllint -d relaxed . --no-warnings
endif
	gofmt -l -w -s .
	go vet ./...
	go test -v ./... -tags test -race -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
coveralls:
	go test -v ./... -tags test -covermode=count -coverprofile=coverage.out -json > test-report.json
	go install github.com/mattn/goveralls
	${GOPATH}/bin/goveralls -coverprofile=coverage.out
build-all:
	gofmt -l -w -s .
	go vet ./...
	go test -v ./... -tags test
	mkdir -p bin
	go build -o bin/ ./...
docker-build:
	docker build -t devatherock/simple-slack:$(docker_tag) \
	    -f build/Dockerfile .
integration-test:
ifneq ($(skip_pull), true)
	docker pull devatherock/simple-slack:$(docker_tag)
endif
	go test -v ./... -tags integration