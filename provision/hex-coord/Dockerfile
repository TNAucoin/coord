FROM golang:1.22 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /coord cmd/hex-coord/main.go

FROM ubuntu:22.04 as runner

RUN apt-get update && rm -rf /var/lib/apt/lists/*

COPY --from=builder /coord /coord

CMD chmod +x /river

CMD ["/coord"]
