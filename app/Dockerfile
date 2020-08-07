FROM golang:alpine AS builder
RUN apk add gcc libc-dev git
COPY ./* /go/src/app/
WORKDIR /go/src/app/
RUN go get ./...
ARG test=yes
RUN if [ "$test" = "yes" ]; then go test || exit 1 ; fi
RUN go build -o nevernote

FROM alpine:latest
WORKDIR /app
COPY --from=builder /go/src/app/nevernote /app/nevernote
CMD ["/app/nevernote"]
