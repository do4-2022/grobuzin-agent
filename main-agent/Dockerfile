FROM alpine:3.19

RUN apk add go

WORKDIR /build

COPY . .

WORKDIR /build/main-agent

RUN CGO_ENABLED=0 go build -o main-agent