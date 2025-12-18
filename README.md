# qrkit

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/v/release/aidanuno/qrkit)](https://github.com/aidanuno/qrkit/releases)
[![CI](https://github.com/aidanuno/qrkit/actions/workflows/ci.yml/badge.svg)](https://github.com/aidanuno/qrkit/actions/workflows/ci.yml)
[![Docker Image](https://img.shields.io/badge/docker-ghcr.io-blue)](https://github.com/aidanuno/qrkit/pkgs/container/qrkit)

QR code generation tools for MCP-enabled LLM clients.

A Model Context Protocol (MCP) server for generating QR codes. Built with Go, the [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk), and [go-qrcode](https://github.com/yeqown/go-qrcode/v2), qrkit enables AI assistants to generate QR codes from any data.

## Features

- **QR Code Generation**: Generate QR codes from URLs, WiFi credentials, text, or any data for easy sharing
- **Dual Transport Modes**:
  - **STDIO Mode**: For local MCP client integration (default)
  - **HTTP Streamable Mode**: For remote access and web-based clients
- **Multiple Deployment Options**:
  - Pre-built executables for Linux, macOS, and Windows
  - Docker container support
- **Lightweight**: Built with Go for minimal resource usage
- **Cross-platform**: Supports Linux, macOS, and Windows

Planned features:

- [ ] Generate QR codes in batch mode
- [ ] Generate custom QR codes

**Example Prompts for Usage within MCP-enabled LLM clients**:

```
Generate a QR code for https://github.com/aidanuno/qrkit
```

```
Create a QR code for my WiFi, "MyNetwork" with password "password123"
```

## Installation

### Option 1: Install with Go

```bash
go install github.com/aidanuno/qrkit@latest
```

### Option 2: Docker

## Usage

### STDIO Transport Mode (Default)

STDIO mode is designed for local MCP clients like Claude Desktop, VSCode, or other MCP-compatible applications. Simply copy the configuration below into your MCP client's configuration file.

#### After Go installation

```json
{
  "mcpServers": {
    "qrkit": {
      "command": "qrkit"
    }
  }
}
```

#### With docker image

```json
{
  "mcpServers": {
    "qrkit": {
      "command": "docker",
      "args": ["run", "-i", "ghcr.io/aidanuno/qrkit:latest"]
    }
  }
}
```

### HTTP Streamable Transport Mode

HTTP mode enables remote access and integration with web-based MCP clients.

#### After Go installation

```bash
qrkit --port 3000
```

#### With docker image

```bash
docker run -p 3000:3000 ghcr.io/aidanuno/qrkit:latest --port 3000
```

The server will be accessible at `http://localhost:3000`.

#### HTTP Mode Configuration

For MCP clients that support HTTP Streamable transport, configure the endpoint:

```json
{
  "mcpServers": {
    "qrkit": {
      "url": "http://localhost:3000"
    }
  }
}
```

## License

MIT License - see [LICENSE.md](LICENSE.md) for details

## Contributing

Contributions are welcome!

## Links

- [GitHub Repository](https://github.com/aidanuno/qrkit)
- [Model Context Protocol](https://modelcontextprotocol.io/)
- [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk)

## Acknowledgments

Built with:

- [github.com/yeqown/go-qrcode/v2](https://github.com/yeqown/go-qrcode) - QR code generation
- [github.com/modelcontextprotocol/go-sdk](https://github.com/modelcontextprotocol/go-sdk) - MCP implementation
