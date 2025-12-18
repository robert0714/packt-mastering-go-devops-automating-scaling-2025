# First stage: Build the application
FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o mycli
# Second stage: Create a smaller image for the final executable
FROM alpine
WORKDIR /root/
COPY --from=builder /app/mycli .
ENTRYPOINT ["./mycli"]