# docker build -t alpine-aws-cli --build-arg CF=test.rtf .
FROM alpine:3.10

LABEL maintainer="Mark Day<mark.day@aistemos.com>"

RUN apk update

# Install base packages
RUN apk add --no-cache --virtual .build-deps
RUN apk add bash

# Install build packages
RUN apk add make && apk add curl && apk add openssh

# Install git
RUN apk add git

# Set timezone to UTC by default
RUN ln -sf /usr/share/zoneinfo/Etc/UTC /etc/localtime

# Install aws-cli
RUN apk -Uuv add groff less python py-pip
RUN pip install awscli
RUN apk --purge -v del py-pip
RUN rm /var/cache/apk/*

# add credentials file
ADD credentials /root/.aws/credentials

# copy down the config file
ARG CF
RUN aws s3 cp s3://aistemos-cipher-configs/${CF} /config-src/${CF}

CMD ["/bin/bash"]
ENTRYPOINT ["sh", "-c", "cp /config-src/* /config/", "--"]
