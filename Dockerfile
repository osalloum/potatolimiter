# Use a base image that includes SSL libraries
FROM golang:1.20-alpine as build

# Install OpenSSL and necessary libraries
RUN apk --no-cache add openssl

# Copy the Go application source code into the container
COPY . /src

# Set the working directory
WORKDIR /src

# Build the Go application
RUN go build -o /main .

# Create a minimal image with the built Go binary and SSL support
FROM alpine:latest

# Install OpenSSL to support SSL
RUN apk --no-cache add openssl ca-certificates

# Copy the binary from the build stage
COPY --from=build /main /main

# Set the SSL certificates directory and ensure they're updated
RUN update-ca-certificates

# Expose the required port (if needed, adjust to your app's listening port)
EXPOSE 443

# Run the application
ENTRYPOINT ["/main"]