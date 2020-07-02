# docker build -t 111177312954.dkr.ecr.eu-west-1.amazonaws.com/aistemos/alpine-awscli:0.1.1 .
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

ENTRYPOINT ["sh", "-c", "aws s3 cp s3://aistemos-classifiers-w2v-cpc/s300_w10_sg1_mc5_hs1.pkl /var/w2v_cpc/", "--"]
