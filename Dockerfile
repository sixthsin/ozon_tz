FROM golang:1.22.12-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /ozontz_app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o my_ozontz_app ./cmd/main.go

FROM alpine:latest

COPY --from=builder /ozontz_app/my_ozontz_app /my_ozontz_app

CMD ["/my_ozontz_app"]
