FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG SUPABASE_URL
ARG SUPABASE_SERVICE_ROLE_KEY

RUN CGO_ENABLED=0 GOOS=linux go build -o server

FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 8080

ENV PORT=8080

CMD ["./server"]
