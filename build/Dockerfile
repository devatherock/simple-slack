ARG GO_VERSION=1.18
FROM golang:${GO_VERSION}-alpine3.14 AS build

COPY . /home/workspace
WORKDIR /home/workspace

RUN mkdir -p bin && \
    go build -o bin/simpleslack ./...


FROM alpine:3.14.10

LABEL maintainer="devatherock@gmail.com"

COPY --from=build /home/workspace/bin/simpleslack /bin/simpleslack

CMD ["/bin/simpleslack"]
