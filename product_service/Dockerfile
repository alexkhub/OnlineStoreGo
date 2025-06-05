FROM golang:1.23.2-alpine as builder

WORKDIR /app
COPY  go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /example cmd/example.go

FROM alpine:3
COPY /configs/config.yaml /configs/config.yaml
COPY --from=builder example /bin/example

ENTRYPOINT ["/bin/example"]
EXPOSE 8083

CMD ["./example"]