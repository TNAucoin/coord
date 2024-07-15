FROM golang:1.22 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /worker cmd/worker/worker.go

FROM ubuntu:22.04 as runner

RUN apt-get update && rm -rf /var/lib/apt/lists/*

COPY --from=builder /worker /worker

CMD chmod +x /worker

CMD ["/worker"]
