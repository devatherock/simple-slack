FROM alpine

COPY release/simpleslack /bin/simpleslack

ENTRYPOINT ["/bin/simpleslack"]