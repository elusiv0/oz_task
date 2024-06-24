FROM golang:1.22.2

WORKDIR /app


COPY . .
RUN go mod download
RUN go mod tidy

RUN go build -o oz_task ./cmd/main.go

CMD ["/oz_task"]