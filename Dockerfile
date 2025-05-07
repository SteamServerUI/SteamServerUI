# Stage 1: Build the Go application
FROM golang:1.24-bullseye AS builder

# Set working directory for build
WORKDIR /build

# Copy Go module files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code to build directory
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

# Stage 2: Create a minimal runtime image
FROM debian:12-slim AS runner

# Define a non-root user and group ID
ARG APP_UID=1000
ARG APP_GID=1000

# Create a non-root user and group
RUN groupadd --gid ${APP_GID} ssui && \
    useradd --uid ${APP_UID} --gid ${APP_GID} --shell /bin/bash --create-home ssui

# Install only the essential runtime dependencies
RUN dpkg --add-architecture i386 \
    && apt-get update -y \
    && apt-get install -y --no-install-recommends ca-certificates locales lib32gcc-s1 \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Create a clean application directory and required subdirectories
WORKDIR /app
RUN mkdir -p /app/saves /app/UIMod/config

# Create entrypoint script with strict permission checking
RUN echo '#!/bin/bash\n\
# Ensure directories exist\n\
mkdir -p /app/saves /app/UIMod/config\n\
\n\
# Check permissions and exit if incorrect\n\
if [ ! -w "/app/saves" ]; then\n\
  echo "ERROR: Cannot write to /app/saves directory"\n\
  echo "===================================================="\n\
  echo "PERMISSION ERROR: Volume mounts have incorrect ownership"\n\
  echo "\n"\n\
  echo "Please run the following commands on your host system:"\n\
  echo "  mkdir -p ./saves ./UIMod/config"\n\
  echo "  sudo chown -R 1000:1000 ./saves ./UIMod/config"\n\
  echo "  chmod -R 755 ./saves ./UIMod/config"\n\
  echo "\n"\n\
  echo "Then restart the container."\n\
  echo "===================================================="\n\
  exit 1\n\
fi\n\
\n\
if [ ! -w "/app/UIMod/config" ]; then\n\
  echo "ERROR: Cannot write to /app/UIMod/config directory"\n\
  echo "===================================================="\n\
  echo "PERMISSION ERROR: Volume mounts have incorrect ownership"\n\
  echo "\n"\n\
  echo "Please run the following commands on your host system:"\n\
  echo "  mkdir -p ./saves ./UIMod/config"\n\
  echo "  sudo chown -R 1000:1000 ./saves ./UIMod/config"\n\
  echo "  chmod -R 755 ./saves ./UIMod/config"\n\
  echo "\n"\n\
  echo "Then restart the container."\n\
  echo "===================================================="\n\
  exit 1\n\
fi\n\
\n\
# If we get here, permissions are correct\n\
# Execute the main application\n\
exec /app/SSUI.x86_64 "$@"' > /app/entrypoint.sh

# Make entrypoint script executable
RUN chmod +x /app/entrypoint.sh

# Copy ONLY the built binary from the builder stage
COPY --from=builder --chown=ssui:ssui /build/build/SSUI*.x86_64 /app/SSUI.x86_64

# Set ownership recursively for the entire /app directory
RUN chown -R ssui:ssui /app
RUN chmod -R 755 /app
RUN chmod +x /app/SSUI.x86_64

# Expose the necessary ports
EXPOSE 8443 27016 27015

# Switch to the non-root user
USER ssui

# Set the entrypoint
ENTRYPOINT ["/app/entrypoint.sh"]

# No default arguments
CMD []