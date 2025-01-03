FROM golang:1.23 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o main .

FROM debian:bookworm-slim


ENV SEED_DATA_DIR=db/seed_data
ENV PORT=8080

WORKDIR /root/

COPY --from=build /app/main .
RUN mkdir -p ${SEED_DATA_DIR}
COPY --from=build /app/${SEED_DATA_DIR}/* ${SEED_DATA_DIR}

RUN apt-get update && apt-get install -y curl bash && \
    curl -o /wait-for-it.sh https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh && \
    chmod +x /wait-for-it.sh

EXPOSE ${PORT}
CMD ["/wait-for-it.sh", "postgres:5432", "--", "./main"]
