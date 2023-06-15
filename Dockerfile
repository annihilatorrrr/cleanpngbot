FROM golang:1.20.5-alpine3.17 as builder
WORKDIR /cleanpngbot
RUN apk update && apk upgrade --available && sync && apk add --no-cache --virtual .build-deps upx
COPY . .
RUN go build -ldflags="-w -s" .
RUN upx /cleanpngbot/cleanpngbot
FROM alpine:3.18.2
RUN apk update && apk upgrade --available && sync
COPY --from=builder /cleanpngbot/cleanpngbot /cleanpngbot
ENTRYPOINT ["/cleanpngbot"]
