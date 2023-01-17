FROM golang:1.18-alpine as builder
RUN mkdir /src
WORKDIR /src
ADD . .
RUN go build -ldflags "-s -w -X main.version=$(cat VERSION)" -o main .

FROM alpine
COPY --from=builder /src/main /app/
COPY bootstrap_*.yaml /app/
WORKDIR /app
EXPOSE 8080
ENTRYPOINT ["/app/main"]
