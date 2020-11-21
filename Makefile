test:
	go test -v -race -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
check:
	gofmt -l -w -s .
	go vet
	go test -v -race -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
coveralls:
	go test -v -covermode=count -coverprofile=coverage.out
	go get github.com/mattn/goveralls
	${GOPATH}/bin/goveralls -coverprofile=coverage.out
build:
	go build -o release/simpleslack