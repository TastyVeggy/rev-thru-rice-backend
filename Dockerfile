FROM golang:1.23 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o main .

FROM debian:bookworm-slim

WORKDIR /root/

COPY --from=build /app/main .
EXPOSE ${PORT}

CMD ["./main"]
