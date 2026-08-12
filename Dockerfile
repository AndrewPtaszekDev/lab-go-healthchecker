FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

FROM alpine:latest

RUN adduser -D user
USER user

WORKDIR /user/

# copies JUST the exe 'main' from prev stage
COPY --from=builder /app/main .

EXPOSE 23234

CMD ["./main"]
