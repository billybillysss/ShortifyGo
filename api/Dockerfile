# Stage 1: Build the binary
FROM golang:alpine AS builder

# Create build directory
RUN mkdir /build

# Add source code to the build directory
ADD . /build/
WORKDIR /build

# Build the Go binary
RUN go build -o main .

# Change the ownership to appuser and make the binary executable
RUN chmod +x /build/main

# Stage 2: Create a minimal image with the binary
FROM alpine

# Create an application user
RUN adduser -S -D -H -h /app appuser

# Use the created user
USER appuser

# Set the working directory
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /build/main /app/
# COPY .env /app/
# Expose the application port
EXPOSE 7001

# Set the default command
CMD [ "./main" ]