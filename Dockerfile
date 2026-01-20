# ========================
# Build stage
# ========================
FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# generate ent
RUN go generate ./ent

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o app ./apps/cmd/server

# ========================
# Runtime stage
# ========================
FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/app .
COPY --from=builder /app/migrations ./migrations

EXPOSE 3098

CMD ["/app/app"]
