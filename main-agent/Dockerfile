FROM alpine:3.19

RUN apk add go

WORKDIR /build

COPY . .

RUN CGO_ENABLED=0 go build -o main-agent