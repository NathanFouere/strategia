# strategia

The goal of this project is to create a minimalist clone of the openfront game.

## Prerequisites

You only need nix installed

## Installation

### Client

```bash
nix develop
cd client
bun install
```

### Server

```bash
nix develop
cd server
go mod tidy
```

## Run

### Client

```bash
nix develop
cd client
bun dev
```

### Server

```bash
nix develop
cd server
go run cmd/main.go
```
