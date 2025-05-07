# Stage 1: Build the Go application
FROM golang:1.24-bullseye AS builder

# Set working directory
WORKDIR /app

# Copy Go module files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Install Node.js and npm for frontend build
RUN curl -fsSL https://deb.nodesource.com/setup_18.x | bash - && \
    apt-get install -y nodejs

# Build frontend and Go application
RUN go mod tidy && \
    cd ssui-interfacev2 && \
    npm install && \
    npm install @sveltejs/vite-plugin-svelte && \
    cd .. && \
    go run ./build/build.go

# Stage 2: Create the runtime image
FROM debian:12-slim AS runner

# Define a non-root user and group ID
ARG APP_UID=1000
ARG APP_GID=1000

# Create a non-root user and group
RUN groupadd --gid ${APP_GID} ssui && \
    useradd --uid ${APP_UID} --gid ${APP_GID} --shell /bin/bash --create-home ssui

# Set the working directory
WORKDIR /app

# Install runtime dependencies
RUN dpkg --add-architecture i386 \
    && apt-get update -y \
    && apt-get install -y --no-install-recommends ca-certificates locales lib32gcc-s1 \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Copy the built binary from the builder stage
COPY --from=builder --chown=ssui:ssui /app/build/SSUI*.x86_64 /app/SSUI.x86_64

# Make the binary executable
RUN chmod +x /app/SSUI.x86_64

# Copy UIMod folder
COPY --chown=ssui:ssui ./UIMod /app/UIMod

# Set ownership
RUN chown -R ssui:ssui /app/

# Expose the necessary ports
EXPOSE 8443 27016 27015

# Switch to the non-root user
USER ssui

# Set the entrypoint
ENTRYPOINT ["/app/SSUI.x86_64"]

# Provide default arguments to the entrypoint
CMD []