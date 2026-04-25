FROM golang:1.26.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /bin/app ./src/cmd/backend

FROM alpine:3.21

RUN adduser -D nonroot
USER nonroot

COPY --from=builder /bin/app /app

EXPOSE 80

CMD ["/app"]
