FROM golang:1.17-alpine as builder
RUN mkdir /src
WORKDIR /src
ADD . .
RUN go get github.com/pilu/fresh
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "-extldflags '-static' -s -w -X main.version=$(cat VERSION)" -o main .
CMD ["fresh"]

FROM alpine
COPY --from=builder /src/kubeops /app/
WORKDIR /app
EXPOSE 8080
#ENTRYPOINT ["/app/kubeops"]
CMD ["/main"]
