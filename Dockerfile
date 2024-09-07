# Use the official Golang image to create a build artifact.
FROM golang:1.23 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./

# Download all dependencies. Dependencies will be cached if the go.mod/sum files are not changed
# as of now we are not using go.sum file so we are using || true to ignore the error
RUN go mod download || true

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN GOOS=linux GOARCH=amd64 go build -o /sms2telegram

# Start a new stage from scratch
FROM scratch

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /sms2telegram /sms2telegram

# Command to run the executable
CMD ["/sms2telegram"]