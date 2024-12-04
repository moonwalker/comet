# Comet 

Cosmic tool for provisioning and managing infrastructure.

## Overview

Comet is a command-line interface (CLI) tool designed to streamline infrastructure provisioning and management. It provides a unified interface for handling infrastructure operations with modern tooling and practices.

## Features

- Infrastructure provisioning and management
- Terraform/OpenTofu integration
- JavaScript configuration language
- Configurable through YAML

## Installation

### Prerequisites

- Go 1.23 or later

### Building from Source

```bash
git clone https://github.com/moonwalker/comet.git
cd comet
go build
```

## Usage

```bash
comet [command] [flags]
```

For detailed command documentation, use:
```bash
comet --help
```

## Configuration

Comet can be configured using `comet.yaml` in your project directory. 

## Development

### Requirements

- Go 1.23+

### Setup

1. Clone the repository
```bash
git clone https://github.com/moonwalker/comet.git
```

2. Build the project
```bash
go build
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the terms specified in the project's license file.
