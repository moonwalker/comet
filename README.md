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

```
comet [command] [flags]
```

For detailed command documentation, use:
```
comet --help
```

## Commands

### `comet version`

**Description:** Print the version.

**Usage:**
```
comet version
```

### `comet plan`

**Description:** Show changes required by the current configuration.

**Usage:**
```
comet plan <stack> [component]
```

### `comet output`

**Description:** Show output values from components.

**Usage:**
```
comet output <stack> [component]
```

### `comet list`

**Description:** List stacks or components.

**Usage:**
```
comet list [stack]
```

### `comet destroy`

**Description:** Destroy previously-created infrastructure.

**Usage:**
```
comet destroy <stack> [component]
```

### `comet clean`

**Description:** Delete Terraform-related folders and files.

**Usage:**
```
comet clean <stack> [component]
```

### `comet apply`

**Description:** Create or update infrastructure.

**Usage:**
```
comet apply <stack> [component]
```

## Configuration

Comet can be configured using `comet.yaml` in your project directory. 

## Development

### Requirements

- Go 1.23+

### Setup

1. Clone the repository
```
git clone https://github.com/moonwalker/comet.git
```
2. Build the project
```
go build
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the terms specified in the project's license file.
