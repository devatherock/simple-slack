FROM alpine

COPY release/simpleslack /bin/simpleslack

CMD ["/bin/simpleslack"]