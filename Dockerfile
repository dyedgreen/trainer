# Build Go App
FROM golang:1.11.2-alpine as build

WORKDIR /go/src/trainer
COPY ./internal ./internal

RUN go install ./internal

# Run environment
FROM alpine

WORKDIR /app
COPY --from=build /go/bin/internal ./app
COPY ./web ./web

EXPOSE 80
ENTRYPOINT ["./app"]
