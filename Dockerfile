# Stage 1: Extract version from Go source file
FROM alpine:latest AS extractor

# Set working directory
WORKDIR /app

# Copy only the specific Go file needed for version extraction
COPY src/config/config.go .

# Extract the version string and write it to a file
# Use sed to find the line, capture the version between quotes, and print only the captured part.
RUN sed -n 's/^\s*Version\s*=\s*"\([^"]*\)".*/\1/p' config.go > version.txt

# --- End of Extractor Stage ---

# --- Start of Runner Stage ---
FROM debian:12-slim AS runner

# Define a non-root user and group ID
ARG APP_UID=1000
ARG APP_GID=1000

# Create a non-root user and group first
RUN groupadd --gid ${APP_GID} stationeers && \
    useradd --uid ${APP_UID} --gid ${APP_GID} --shell /bin/bash --create-home stationeers

# Set the working directory inside the container
WORKDIR /app

# Define an argument for the release version (replace 'latest' or specific tag as needed)
ARG RELEASE_TAG=latest
# Define the GitHub repository
ARG GITHUB_REPO=JacksonTheMaster/StationeersServerUI
# Define the base name of the asset (without version/arch)
ARG BASE_ASSET_NAME=StationeersServerControl
# Define the architecture suffix for the asset
ARG ASSET_ARCH=x86_64

# Add dependencies and tools needed
RUN dpkg --add-architecture i386 \
    && apt-get update -y \
    && apt-get install -y --no-install-recommends ca-certificates locales lib32gcc-s1 file curl \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Copy the extracted version file from the 'extractor' stage
# Set ownership to the non-root user
COPY --from=extractor --chown=stationeers:stationeers /app/version.txt /app/version.txt

# Attempt to copy a pre-built binary matching the pattern from the build context's ./build directory
# Use --chown for proper ownership.
# The target name is set directly to the final desired name.
COPY --chown=stationeers:stationeers ./build/StationeersServerControl*.x86_64 /app/StationeersServerControl.x86_64

# Download if necessary, make executable, and verify
# Run these steps as root initially for permissions to install/download
RUN \
    # Check if the binary was successfully copied from ./build in the previous step
    if [ -f "/app/StationeersServerControl.x86_64" ]; then \
        echo "Using pre-built binary found in ./build/"; \
    else \
        # If not found locally, proceed to download from GitHub
        echo "No pre-built binary found in ./build/. Downloading from GitHub..."; \
        echo "Reading version from extracted file..." && \
        VERSION=$(cat /app/version.txt) && \
        if [ -z "${VERSION}" ]; then \
            echo "Error: Could not read version from /app/version.txt." >&2; \
            exit 1; \
        fi && \
        # Construct the asset name dynamically using ARGs and the extracted version
        DYNAMIC_ASSET_NAME="${BASE_ASSET_NAME}v${VERSION}.${ASSET_ARCH}" && \
        echo "Constructed asset name: ${DYNAMIC_ASSET_NAME}" && \
        # Proceed with download using the dynamic name
        echo "Downloading release ${RELEASE_TAG} from ${GITHUB_REPO}, asset ${DYNAMIC_ASSET_NAME}..." && \
        curl --fail --silent --show-error -L -o /app/StationeersServerControl.x86_64 \
            "https://github.com/${GITHUB_REPO}/releases/download/${RELEASE_TAG}/${DYNAMIC_ASSET_NAME}" || \
        # Handle download failure
        ( echo "Error: Failed to download asset '${DYNAMIC_ASSET_NAME}' from release '${RELEASE_TAG}'. Check GITHUB_REPO, RELEASE_TAG, and ensure the asset name format matches the release." >&2; exit 1 ); \
    fi && \
    \
    # Make the binary executable (whether copied or downloaded)
    echo "Making binary executable..." && \
    chmod +x /app/StationeersServerControl.x86_64 && \
    \
    # Verify that the executable exists and is executable
    echo "Verifying the StationeersServerControl executable..." && \
    if [ -f "/app/StationeersServerControl.x86_64" ] && [ -x "/app/StationeersServerControl.x86_64" ]; then \
        echo "Verification successful: /app/StationeersServerControl.x86_64 exists and is executable."; \
        echo "File details:"; \
        ls -l /app/StationeersServerControl.x86_64; \
        file /app/StationeersServerControl.x86_64; \
    else \
        echo "Error: Verification failed. /app/StationeersServerControl.x86_64 not found or not executable." >&2; \
        exit 1; \
    fi && \
    # Ensure the final binary is owned by the non-root user
    chown stationeers:stationeers /app/StationeersServerControl.x86_64

# COPY ./LICENSE /app/LICENSE # Keep commented unless needed

# Copy the UIMod folder into the application directory, owned by the non-root user
COPY --chown=stationeers:stationeers ./UIMod /app/UIMod

# Expose the ports (doesn't require root)
EXPOSE 8443 27016 27015

# Switch to the non-root user
USER stationeers

# Set the entrypoint to the application using the consistent name
ENTRYPOINT ["/app/StationeersServerControl.x86_64"]

# Provide default arguments to the entrypoint
CMD []