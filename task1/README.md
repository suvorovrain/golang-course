# CLI tool for getting information about the GitHub repository

A simple command-line tool written in **Go** that fetches and displays key information about a public GitHub repository using the official GitHub API.

---

## Features

- Fetches repository metadata by owner and name.
- Displays:
    - Repository Name
    - Description
    - Star Count
    - Fork Count
    - Creation Date

## Prerequisites

- [Go](https://go.dev/dl/) (version 1.20 or later recommended) installed on your machine.

## Installation & Usage

You don't need to install the binary globally. You can run the tool directly from the source code.

### 1. Clone the repository

```bash
git clone <your-repository-url>
cd <your-repository-directory>
```

### 2. Run the tool

```bash
go run main.go -owner <OWNER> -repo <REPO_NAME>
```
