# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Lightning is a custom RPC (Remote Procedure Call) framework written in Go that implements a private binary protocol for efficient data serialization. The project demonstrates how to build an RPC system without relying on reflection or external serialization frameworks like Protocol Buffers.

## Key Components

### Core Architecture
- **Server** (`server/server.go`): RPC server implementation that handles connections, request parsing, and response generation
- **Client** (`client/client.go`): Client-side implementation for making RPC calls  
- **Codec** (`codec/`): Pluggable encoding/decoding interface with multiple implementations:
  - `gob.go`: Go's built-in gob encoding
  - `json_codec.go`: JSON-based encoding
  - `lightning_codec.go`: Custom Lightning protocol codec
- **Serialization** (`serialization/`): Lightning protocol implementation for binary serialization
- **Main** (`main/main.go`): Example server and client usage

### Lightning Protocol Details
- Custom binary serialization supporting `int` and `string` types only
- Uses magic numbers for packet boundaries (start: `[33, 44, 55, 12]`, end: `[23, 84, 55, 29]`)
- Big-endian byte ordering for consistent cross-platform compatibility
- Packet format: `[MAGIC_START][param_count][param_types][param_lengths][param_data][MAGIC_END]`

## Common Development Tasks

### Running the Application
```bash
go run main/main.go
```

### Running Tests
```bash
go test ./...                    # Run all tests
go test ./serialization          # Test Lightning protocol serialization
go test -v ./serialization       # Verbose test output
```

### Building
```bash
go build ./main                  # Build the main application
go build .                       # Build all packages
```

### Code Quality
```bash
go fmt ./...                     # Format code
go vet ./...                     # Static analysis
go mod tidy                      # Clean up dependencies
```

## Protocol Implementation Notes

The Lightning protocol currently supports only `int64` and `string` types but can be extended. The serialization uses reflection to marshal/unmarshal struct fields, making it suitable for complex nested objects containing these supported types.

Key constants:
- Magic number: `0x3bef5c`
- Supported codecs: `application/gob`, `application/json`, `application/lightning`
- Default timeout: 8 seconds (hardcoded in server response handling)

## Important Notes
The Lightning protocol currently supports only `int64` and `string` types but can be extended. The serialization uses reflection to marshal/serialize struct fields, making it suitable for complex nested objects containing these supported types.