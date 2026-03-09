# GitHub Repository Info CLI

A simple command-line tool written in Go that fetches and displays information about public GitHub repositories using the GitHub API.

## Features

- Displays key repository statistics (stars, forks, creation date)
- Comprehensive error handling (404, 403, rate limits)
- Fast and lightweight single-binary distribution
- Clean, formatted output

## Requirements

- [Go](https://go.dev/dl/) version 1.21 or higher

## Installation

### Option 1: Run directly with Go (for development)

```bash
git clone https://github.com/YOUR_USERNAME/YOUR_REPO.git
cd YOUR_REPO
go run . -name <owner> -repo <repository>
```

### Option 2: Build and install (for production use)
# Build the binary
go build -o github-cli

# Run the binary
./github-cli -name <owner> -repo <repository>

Windows users: Use github-cli.exe instead of ./github-cli

### Usage
# Basic usage
./github-cli -name golang -repo go

# View help
./github-cli -h