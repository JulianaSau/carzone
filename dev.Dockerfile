FROM golang:latest

# Set the Current Working Directory inside the container
WORKDIR /app

RUN mkdir "/build"

ENV PROJECT_DIR=/app \
    GO111MODULE=on \
    CGO_ENABLED=0

# Copy the source from the current directory to the Working Directory inside the container
ADD . .

RUN go get github.com/githubnemo/CompileDaemon
RUN go install github.com/githubnemo/CompileDaemon

# Build the Go app
ENTRYPOINT CompileDaemon -build="go build -o /build/app" -command="/build/app"
