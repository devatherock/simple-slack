ARG GO_VERSION=1.18
FROM golang:${GO_VERSION}-alpine3.17 AS build

COPY . /home/workspace
WORKDIR /home/workspace

RUN go build -o bin/simpleslack ./...


FROM alpine:3.17.4

LABEL maintainer="devatherock@gmail.com"

COPY --from=build /home/workspace/bin/simpleslack /bin/simpleslack

CMD ["/bin/simpleslack"]
