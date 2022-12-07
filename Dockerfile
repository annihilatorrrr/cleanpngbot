FROM golang:1.19.4-alpine3.16 as builder
WORKDIR /cleanpngbot
RUN apk update && apk upgrade --available && sync
COPY . .
RUN go build -ldflags="-w -s" .
RUN rm -rf *.go && rm -rf go.*
FROM alpine:3.17.0
RUN apk update && apk upgrade --available && sync
COPY --from=builder /cleanpngbot/cleanpngbot /cleanpngbot
ENTRYPOINT ["/cleanpngbot"]
