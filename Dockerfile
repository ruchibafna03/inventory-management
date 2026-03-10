# ── Stage 1: Build React frontend ────────────────────────────────────────────
FROM node:18-alpine AS frontend
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci --silent
COPY frontend/ ./
RUN npm run build
# output → /app/web  (vite outDir: '../web')

# ── Stage 2: Build Go backend ─────────────────────────────────────────────────
FROM golang:1.23-alpine AS backend
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server ./cmd/server

# ── Stage 3: Final image ──────────────────────────────────────────────────────
FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=backend /app/server .
COPY --from=backend /app/migrations ./migrations/
COPY --from=frontend /app/web ./web/
ENV PORT=8080
ENV WEB_ROOT=/app/web
EXPOSE 8080
CMD ["./server"]
