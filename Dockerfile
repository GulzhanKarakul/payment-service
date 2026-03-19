# ========= BUILD STAGE ============
FROM golang:1.21-alpine AS builder
# install git for private dependencies
RUN apk add --no-cache git
WORKDIR /app
# cache dependencies for faster build
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
# build static binary for linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
		-ldflags="-w -s" -o main ./cmd/main.go

# ========== FINAL STAGE ============
FROM alpine:latest
# add SSL certificates and time zone data
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
# add user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
COPY --from=builder /app/main .
EXPOSE 8080
USER appuser
CMD ["./main"]