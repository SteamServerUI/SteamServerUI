# Use Debian-based VS Code dev container image
FROM mcr.microsoft.com/vscode/devcontainers/base:bookworm

# Install Go 1.24.2
RUN wget -q https://go.dev/dl/go1.24.2.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.24.2.linux-amd64.tar.gz && \
    rm go1.24.2.linux-amd64.tar.gz
ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH=/go
RUN mkdir -p /go/bin && chown -R vscode:vscode /go

# Install Node.js 22.15.0
RUN curl -fsSL https://deb.nodesource.com/setup_22.x | bash - && \
    apt-get install -y nodejs && \
    npm install -g npm@latest

RUN dpkg --add-architecture i386 && \
    apt-get update && \
    apt-get install -y lib32gcc-s1

# Install build tools for cross-compilation
RUN apt-get update && apt-get install -y \
    git \
    curl \
    make \
    && rm -rf /var/lib/apt/lists/*

# Set up non-root user
USER vscode
WORKDIR /workspaces