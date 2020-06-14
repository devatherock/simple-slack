FROM alpine

LABEL maintainer="devatherock@gmail.com"
LABEL io.github.devatherock.version="0.4.0"

COPY release/simpleslack /bin/simpleslack

CMD ["/bin/simpleslack"]
