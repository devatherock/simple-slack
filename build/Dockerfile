ARG GO_VERSION=1.20
FROM golang:${GO_VERSION}-alpine3.19 AS build

COPY . /home/workspace
WORKDIR /home/workspace

RUN go build -o bin/ ./cmd/plugin


FROM alpine:3.19.1

LABEL maintainer="devatherock@gmail.com"

COPY --from=build /home/workspace/bin/plugin /bin/plugin

CMD ["/bin/plugin"]
