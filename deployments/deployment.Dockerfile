# ================================
# Build image
# ================================
FROM golang:1.21-alpine as build

# Install OS updates and, if needed, sqlite3
RUN apk update \
    && apk upgrade \
    && apk add --no-cache git

# Set up a build area
WORKDIR /go/src/go-docker-dev.to

# First just resolve dependencies.
# This creates a cached layer that can be reused
# as long as your go.mod files does not change.
COPY go.mod go.sum ./
RUN go mod download

# Copy entire repo into container
COPY . .

# Build everything, with optimizations
RUN go build cmd/main/main.go

# Switch to the staging area
WORKDIR /staging

# Copy main executable to staging area
RUN cp /go/src/go-docker-dev.to/main ./

# Copy any resouces from the public directory and views directory if the directories exist
# Ensure that by default, neither the directory nor any of its contents are writable.
RUN [ -d /go/src/go-docker-dev.to/assets ] && { mv /go/src/go-docker-dev.to/assets ./assets && chmod -R a-w ./assets; } || true

# ================================
# Run image
# ================================
FROM alpine

# Make sure all system packages are up to date.
# RUN apk add

# Create a user and group with /app as its home directory
RUN addgroup -S syndya && adduser -S syndya -G syndya -h /app

# Switch to the new home directory
WORKDIR /app

# Copy built executable and any staged resources from builder
COPY --from=build --chown=syndya:syndya /staging /app

# Ensure all further commands run as the app user
USER syndya:syndya

# Let Docker bind to port 8080
EXPOSE 8080

ENV MATCHFINDER_LUASCRIPT=assets/matchup.lua

# Start the service when the image is run, default to listening on 8080 in production environment
ENTRYPOINT ["./main"]
