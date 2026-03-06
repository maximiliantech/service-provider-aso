#!/usr/bin/env bash
set -euo pipefail

echo "==> Initializing git submodules..."
git submodule update --init --recursive

echo "==> Fixing Go module cache permissions..."
sudo mkdir -p /go/pkg/mod
sudo chown -R "$(id -u):$(id -g)" /go/pkg/mod

echo "==> Downloading Go module dependencies..."
go mod download
(cd api && go mod download)

echo "==> Dev container setup complete!"
