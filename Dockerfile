# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:latest as builder

# Add Maintainer Info
LABEL maintainer="Sheshagiri Mallipedhi <msheshagirirao@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o http-service .

######## Start a new stage from scratch #######
FROM scratch

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/http-service .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./http-service"]