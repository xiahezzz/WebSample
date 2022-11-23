FROM golang:1.18-alpine as builder

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go env -w GOPROXY="https://goproxy.cn,direct"\
    && go env -w GO111MODULE="on"\
    && go mod download \
    && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app ./main.go

#FROM alpine:latest as prod
FROM ubuntu:latest as prod
RUN apt-get update -y && apt-get install -y locales && rm -rf /var/lib/apt/lists/* \
	&& localedef -i en_US -c -f UTF-8 -A /usr/share/locale/locale.alias en_US.UTF-8\
    && apt-get update -y\
    && apt-get install -y musl\
    && apt-get install -y iptables\
    && apt-get install -y net-tools
ENV LANG en_US.utf8

WORKDIR /usr/src/app/

COPY --from=builder /usr/local/bin/app /usr/src/app/

EXPOSE 8080

ENV PROTEST=""

CMD ["./app","${PROTEST}"]