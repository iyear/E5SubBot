FROM golang:alpine as builder

WORKDIR /root

RUN apk update && apk add git \
    && git clone https://github.com/iyear/E5SubBot.git \
    && cd E5SubBot && go build

FROM alpine:latest

ENV TIME_ZONE=Asia/Shanghai

RUN apk update && apk add tzdata \
    && ln -snf /usr/share/zoneinfo/$TIME_ZONE /etc/localtime && echo $TIME_ZONE > /etc/timezone

WORKDIR /root

COPY --from=builder /root/E5SubBot/main /root

CMD [ "/root/main" ]
