FROM golang:1.25-alpine AS builder
RUN mkdir /src
WORKDIR /src
COPY . .
RUN go build -ldflags "-s -w -X main.version=$(cat VERSION)" -o main .

FROM alpine:3.22
COPY --from=builder /src/main /app/
COPY bootstrap_*.yaml /app/
WORKDIR /app
EXPOSE 8080
ENTRYPOINT ["/app/main"]
