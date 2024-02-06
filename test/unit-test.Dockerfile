# ================================
# Testing image
# ================================
FROM golang:1.16-alpine

# Install OS updates and, if needed, sqlite3
RUN apk update \
    && apk upgrade \
    && apk add --no-cache git build-base

# Set up a build area
WORKDIR /go/src/go-docker-test.to

# First just resolve dependencies.
# This creates a cached layer that can be reused
# as long as your go.mod files does not change.
COPY go.mod go.sum ./
RUN go mod download

# Copy entire repo into container
COPY . .

# Start the package tests when the image runs
CMD ["go", "test", "-v", "./..."]
