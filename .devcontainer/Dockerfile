FROM mcr.microsoft.com/devcontainers/go:1-1.23-bullseye

# Set working directory (match with devcontainer.json workspaceFolder)
WORKDIR /workspace

# Install MySQL client + basic tools
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
      default-mysql-client \
      curl \
      tar \
      sudo && \
    rm -rf /var/lib/apt/lists/*

# Install golang-migrate
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz \
      -o /tmp/migrate.tar.gz \
  && tar -xzf /tmp/migrate.tar.gz -C /usr/local/bin \
  && rm /tmp/migrate.tar.gz
