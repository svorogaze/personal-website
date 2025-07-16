FROM golang:1.23.11-alpine3.22 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /myapp

FROM alpine:3.22
WORKDIR /
COPY --from=builder /myapp /myapp
EXPOSE 8081
ENTRYPOINT ["/myapp"]

