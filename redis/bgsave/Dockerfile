FROM alpine:3.12.0
MAINTAINER Mark Day <mark.day@aistemos.com>

COPY ./bgsave /usr/local/bin/bgsave
RUN chmod 0777 /usr/local/bin/bgsave

ENTRYPOINT ["/usr/local/bin/bgsave"]
