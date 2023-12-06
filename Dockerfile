FROM golang:1.21.5

WORKDIR /app

COPY . .

RUN go mod download
RUN go mod tidy

CMD ["go", "run", "main.go"]