# Use official golang image as a parent image
FROM golang:1.21

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed 
RUN go mod download 

# Copy the source from the current directory to the Working Directory inside the container 
COPY . .

# Build the Go app
RUN go build -o api2kafka .

# This container exposes port 8080 to the outside world
EXPOSE 8086

# Run the binary program produced by go install
#CMD ["./kafka_listener"]
CMD sh -c "go run api2kafka/api2kafka.go"
