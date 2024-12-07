ARG GO_VERSION=1.22
FROM golang:${GO_VERSION}-alpine3.20 AS build

COPY . /home/workspace
WORKDIR /home/workspace

RUN go build -o bin/ ./cmd/app


FROM alpine:3.21.0

COPY --from=build /home/workspace/bin/app /bin/app

ENTRYPOINT ["/bin/app"]