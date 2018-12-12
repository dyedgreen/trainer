# Build Go App
FROM golang:1.11.2-alpine as build

WORKDIR /go/src/server
COPY ./server .

RUN go install .

# Run environment
FROM alpine

WORKDIR /app
COPY --from=build /go/bin/server ./server

# Expose ports
EXPOSE 80

# Run app
ENTRYPOINT ["./server"]
