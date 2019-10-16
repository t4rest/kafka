FROM golang:1.13-alpine AS builder
WORKDIR /build
COPY go.* ./
RUN go mod download
COPY . .
RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o go-pub


FROM alpine:latest
COPY --from=builder /build/go-pub /
EXPOSE 3000
CMD ["/go-pub"]