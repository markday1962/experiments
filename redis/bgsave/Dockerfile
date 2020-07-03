FROM alpine:3.12.0
MAINTAINER Mark Day <mark.day@aistemos.com>

COPY ./aist-bgsave /usr/local/bin/bgsave
RUN chmod 0777 /usr/local/bin/bgsave
RUN ls -al /usr/local/bin/

ENTRYPOINT ["sh", "-c", "/usr/local/bin/bgsave"]