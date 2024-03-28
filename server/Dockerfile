FROM golang:1.21.4 AS builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /main cmd/chatserver/main.go

FROM scratch

COPY --from=builder /main /main

EXPOSE 80

CMD ["/main"]
