FROM golang:1.19.3-alpine3.16 as builder
WORKDIR /cleanpngbot
RUN apk update && apk upgrade --available && sync
COPY . .
RUN go build -ldflags="-w -s" .
RUN rm -rf *.go && rm -rf go.*
FROM alpine:3.16.3
RUN apk update && apk upgrade --available && sync
COPY --from=builder /cleanpngbot/cleanpngbot /cleanpngbot
ENTRYPOINT ["/cleanpngbot"]
