FROM golang:1.22.2-bullseye as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -trimpath -ldflags "-w -s" -o app

FROM debian:bullseye-slim as deploy
RUN apt-get update
COPY --from=builder /app/app .

CMD ["./app"]

FROM golang:1.22.2 as dev
WORKDIR /app
RUN go install github.com/air-verse/air@latest
CMD ["air"]
