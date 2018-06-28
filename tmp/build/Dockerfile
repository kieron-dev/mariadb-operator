FROM alpine:3.6

RUN adduser -D mysql-operator
USER mysql-operator

ADD tmp/_output/bin/mysql-operator /usr/local/bin/mysql-operator
