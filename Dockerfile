# First stage: Build the Go application using the Golang image
FROM golang:1.22.1 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download the dependencies
RUN echo "Downloading dependencies..." && go mod download

# Copy the rest of the application source code
COPY . .

# Build the initial executable
RUN echo "Building StationeersServerUI..." && go build -o StationeersServerUI ./build.go

# Second stage: Bootstrap the server using a Debian slim image
FROM golang:1.22.1 AS bootstrapper

# Set the working directory inside the container
WORKDIR /app

# Install required libraries
RUN apt-get update && apt-get install -y \
    lib32gcc-s1 \
    libc6 \
    wget \
    tar \
    gzip \
    && rm -rf /var/lib/apt/lists/*

# Copy the initial executable from the builder stage
COPY --from=builder /app/StationeersServerUI /app/StationeersServerUI

# Copy the rest of the application source code
COPY --from=builder /app /app

# Download and install SteamCMD
RUN wget https://steamcdn-a.akamaihd.net/client/installer/steamcmd_linux.tar.gz && \
    tar -xvzf steamcmd_linux.tar.gz && \
    rm steamcmd_linux.tar.gz && \
    chmod +x steamcmd.sh

# Run the initial executable to build StationeersServerControl
RUN echo "Running StationeersServerUI to build StationeersServerControl..." && ./StationeersServerUI

# Verify that the resulting executable exists
RUN echo "Verifying the existence of StationeersServerControl executable:" && \
    if ls -l /app/StationeersServerControl*; then \
        echo "StationeersServerControl build successful."; \
    else \
        echo "Error: StationeersServerControl executable not found."; \
        exit 1; \
    fi

# Copy the UIMod directory
COPY --from=builder /app/UIMod /app/UIMod

# Third stage: Run the final application using the steamcmd/steamcmd image
FROM steamcmd/steamcmd:latest AS runner

# Set the working directory inside the container
WORKDIR /app

# Install required libraries
#RUN echo "Installing required libraries..." && apt-get update && apt-get install -y \
#    lib32gcc-s1 \
#    libc6 \
#    && rm -rf /var/lib/apt/lists/*

# Copy the resulting executable from the bootstrapper stage and rename it
COPY --from=bootstrapper /app/StationeersServerControl* /app/StationeersServerControl

# Verify that the executable was copied and renamed successfully
RUN echo "Verifying the copied and renamed StationeersServerControl executable:" && \
    if ls -l /app/StationeersServerControl; then \
        echo "StationeersServerControl copy and rename successful."; \
    else \
        echo "Error: StationeersServerControl executable not found after copy."; \
        exit 1; \
    fi

# Copy the UIMod directory
COPY --from=bootstrapper /app/UIMod /app/UIMod

# Expose the ports
EXPOSE 8080 27016

# Run the application and ensure proper handling of stdin, stdout, and stderr
#CMD ["sh", "-c", "/app/StationeersServerControl < /dev/stdin > /dev/stdout 2> /dev/stderr"]
#CMD ["sh", "-c", "/app/StationeersServerControl < /dev/stdin > /dev/stdout 2> /root/.local/share/Steam/logs/stderr.txt"]
# Run the application and redirect both stdout and stderr to a log file
CMD ["sh", "-c", "/app/StationeersServerControl > /app/server.log 2>&1 & tail -f /app/server.log"]
# bare naked exe run
#CMD ["/app/StationeersServerControl"]