# Use the official Golang image as the base image
FROM golang:1.22-alpine

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Set the working directory to the 'cmd' directory where main.go is located
WORKDIR /app/cmd

# Build the Go app
RUN go build -o /app/main .

# Expose port 4000 to the outside world
EXPOSE 4000

# Command to run the executable
CMD ["/app/main"]
