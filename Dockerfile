FROM alpine

LABEL maintainer="devatherock@gmail.com"

COPY release/simpleslack /bin/simpleslack

CMD ["/bin/simpleslack"]
