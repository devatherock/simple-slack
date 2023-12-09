ARG GO_VERSION=1.19
FROM golang:${GO_VERSION}-alpine3.18 AS build

COPY . /home/workspace
WORKDIR /home/workspace

RUN go build -o bin/simpleslack ./...


FROM alpine:3.19.0

LABEL maintainer="devatherock@gmail.com"

COPY --from=build /home/workspace/bin/simpleslack /bin/simpleslack

CMD ["/bin/simpleslack"]
