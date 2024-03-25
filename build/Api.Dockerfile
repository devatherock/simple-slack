ARG GO_VERSION=1.20
FROM golang:${GO_VERSION}-alpine3.19 AS build

COPY . /home/workspace
WORKDIR /home/workspace

RUN go build -o bin/ ./cmd/app


FROM alpine:3.19.1

COPY --from=build /home/workspace/bin/app /bin/app

ENTRYPOINT ["/bin/app"]