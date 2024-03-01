
FROM golang:alpine AS builder

RUN apk update && apk add --no-cache 'git=~2'

ENV GO111MODULE=on
WORKDIR /app
COPY . .

RUN go mod download

RUN go build -o main cmd/api/main.go

FROM alpine:3

WORKDIR /


COPY --from=builder /app/main .

ENV PORT 8080
ENV GIN_MODE release
EXPOSE 8080

CMD ["./main"]
