FROM golang:1.22-alpine3.19

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first to leverage caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project directory into the container
COPY . ./

# Build the Go application
RUN go build -o /anne-hub

# Expose the application port
EXPOSE 1323

# Command to run the application
CMD ["/anne-hub"]
