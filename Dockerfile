# Build stage
FROM golang:1.19 AS build

# Install FFmpeg
RUN apt-get update && \
    apt-get install -y ffmpeg && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Copy and build the Go application
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o myapp .

# Final stage
FROM alpine:3.14
RUN apk add --no-cache ffmpeg

# Copy the binary from the build stage
COPY --from=build /app/myapp /usr/local/bin/myapp

# Run the application
CMD ["myapp"]
