FROM golang:1.20.4-alpine3.17 as builder
WORKDIR /cleanpngbot
RUN apk update && apk upgrade --available && sync && apk add --no-cache --virtual .build-deps upx
COPY . .
RUN go build -ldflags="-w -s" .
RUN rm -rf *.go && rm -rf go.* && upx /cleanpngbot/cleanpngbot && apk --purge del .build-deps
FROM alpine:3.18.0
RUN apk update && apk upgrade --available && sync
COPY --from=builder /cleanpngbot/cleanpngbot /cleanpngbot
ENTRYPOINT ["/cleanpngbot"]
