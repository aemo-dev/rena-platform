FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY rena/backend/go.mod rena/backend/go.sum ./
RUN GOTOOLCHAIN=auto go mod download

COPY rena/backend/*.go ./

RUN GOTOOLCHAIN=auto CGO_ENABLED=0 GOOS=linux go build -o server

FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 8080

ENV PORT=8080

CMD ["./server"]
