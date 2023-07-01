FROM golang:1.20.5-alpine3.18 as builder
WORKDIR /cleanpngbot
RUN apk update && apk upgrade --available && sync && apk add --no-cache --virtual .build-deps
COPY . .
RUN go build -ldflags="-w -s" .
FROM alpine:3.18.2
RUN apk update && apk upgrade --available && sync
COPY --from=builder /cleanpngbot/cleanpngbot /cleanpngbot
ENTRYPOINT ["/cleanpngbot"]
