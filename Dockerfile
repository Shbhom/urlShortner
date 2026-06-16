# Stage 1: Frontend Builder
FROM node:22-alpine AS frontend-builder

WORKDIR /app/Dashboard

COPY Dashboard/package*.json ./
RUN npm ci

COPY Dashboard/ ./
RUN npm run build

# Stage 2: Go Builder

FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

COPY --from=frontend-builder /app/Dashboard/dist ./internal/server/ui/dist

#CGO -> binary is statically linked no external packages required
RUN CGO_ENABLED=0 GOOS=linux go build -o url-shortener ./cmd/server/main.go

# Stage 2 Runner
FROM alpine:latest AS runner

RUN apk --no-cache add ca-certificates

RUN adduser -D nonroot

WORKDIR /app

COPY --from=builder /app/url-shortener  .

USER nonroot

ENTRYPOINT ["./url-shortener"]
