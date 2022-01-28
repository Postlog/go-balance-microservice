FROM golang:alpine
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download
RUN go mod verify

COPY . .
RUN go build -o server ./cmd/

CMD ["./server", "--config", "config/prod.json"]
