FROM golang:latest AS builder

WORKDIR /app
COPY . .

RUN go mod download

RUN go build -o main cmd/api/main.go

FROM ubuntu:latest

RUN apt-get update && apt-get install -y wget

RUN wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb && \
    apt install -y ./google-chrome-stable_current_amd64.deb && \
    rm google-chrome-stable_current_amd64.deb

WORKDIR /

COPY --from=builder /app/main .

ENV PORT 8080
ENV GIN_MODE release
EXPOSE 8080

CMD ["./main"]
