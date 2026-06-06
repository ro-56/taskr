FROM ubuntu:20.04

# COPY --from=ghcr.io/astral-sh/uv:latest /uv /uvx /bin/

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    curl \
    git \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

# Install Go
ARG GO_VERSION=1.26.4
RUN curl -fsSL "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz" \
    | tar -C /usr/local -xz
ENV PATH="/usr/local/go/bin:${PATH}"

ARG USERNAME=devuser

RUN useradd -m -s /bin/bash $USERNAME

# Persist bash history.
RUN SNIPPET="export PROMPT_COMMAND='history -a' && export HISTFILE=/commandhistory/.bash_history" \
  && mkdir /commandhistory \
  && touch /commandhistory/.bash_history \
  && chown -R $USERNAME /commandhistory

# Set `DEVCONTAINER` environment variable to help with orientation
ENV DEVCONTAINER=true

RUN mkdir -p /workspace && \
    chown -R $USERNAME /workspace

WORKDIR /workspace

USER $USERNAME

# Install Claude Code CLI
RUN curl -fsSL https://claude.ai/install.sh | bash
