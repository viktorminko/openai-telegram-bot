FROM golang:1.19

# Install FFmpeg
RUN apt-get update && \
    apt-get install -y ffmpeg && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Copy your application code and build it
COPY . /app
WORKDIR /app
RUN go build -o myapp

# Run the application
CMD ["/app/myapp"]
