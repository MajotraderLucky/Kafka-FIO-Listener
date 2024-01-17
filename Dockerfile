# Use the official golang image as a parent image
FROM golang:1.21

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed 
RUN go mod download 

# Copy the source from the current directory to the Working Directory inside the container 
COPY . .

# Build the Go app
RUN go build -o main .

# This container exposes port 8080 to the outside world
EXPOSE 8000

# Run the binary program produced by go install
CMD ["./main"]
