FROM golang:1.23

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o task ./cmd

CMD ["sh", "-c", "./task $PORT"]