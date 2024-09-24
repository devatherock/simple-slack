docker_tag=latest
skip_pull=false
lint_yaml=false
go_version=1.22

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
run-api:
	go build -o bin/ ./...
	./bin/app
build-all:
	gofmt -l -w -s .
	go vet ./...
	go test -v ./... -tags test
	mkdir -p bin
	go build -o bin/ ./...	
docker-build-plugin:
	docker build -t devatherock/simple-slack:$(docker_tag) \
	    --build-arg GO_VERSION=$(go_version) \
	    -f build/Plugin.Dockerfile .
integration-test-plugin:
ifneq ($(skip_pull), true)
	docker pull devatherock/simple-slack:$(docker_tag)
endif
	go test -v ./... -tags integration
docker-build-api:
	docker build -t devatherock/simple-slack-api:$(docker_tag) \
	    --build-arg GO_VERSION=$(go_version) \
	    -f build/Api.Dockerfile .
integration-test-api:
ifneq ($(skip_pull), true)
	docker pull devatherock/simple-slack-api:$(docker_tag)
endif
	DOCKER_TAG=$(docker_tag) docker compose -f build/docker-compose.yml up --wait
	go test -v ./... -tags api
	docker-compose -f build/docker-compose.yml down
deploy:
	docker run --rm \
        -e FLY_API_TOKEN \
        -v $(CURDIR):/work \
        -w /work \
        flyio/flyctl:v0.2.93 deploy --image devatherock/simple-slack-api:$(docker_tag)
