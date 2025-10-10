---
sidebar_position: 1
---

# Getting Started with Comet

Welcome to **Comet** - a cosmic tool for provisioning and managing infrastructure with JavaScript-based configuration.

## What is Comet?

Comet is a command-line interface (CLI) tool designed to streamline infrastructure provisioning and management. It provides a unified interface for handling infrastructure operations with modern tooling and practices, built on top of Terraform/OpenTofu.

## Why Comet?

Comet fills the gap between plain Terraform/OpenTofu and heavy enterprise frameworks, offering a pragmatic solution for teams that need DRY infrastructure configurations without the overhead of complex tooling.

### Key Benefits

- **🚀 JavaScript Configuration** - Define infrastructure using a familiar, powerful programming language instead of limited HCL
- **🔄 Automatic Backend Generation** - Comet generates `backend.tf.json` files automatically based on your stack configuration
- **🔗 Cross-Stack References** - Simple `state()` template function to reference outputs from other stacks
- **🔐 Built-in Secrets Management** - Native SOPS integration for encrypted secrets in your configurations
- **📦 Component Reusability** - Define components once, reuse across environments with different configurations
- **🎯 Multi-Environment Support** - Manage dev, staging, production, and any other environments from a single codebase

## Installation

### Prerequisites

- Go 1.23 or later
- Terraform or OpenTofu installed

### Building from Source

```bash
git clone https://github.com/moonwalker/comet.git
cd comet
go build
```

This will create a `comet` binary in your current directory. You can move it to a directory in your PATH for easy access.

## Quick Start

### 1. Initialize Your Project

Create a `comet.yaml` file in your project root:

```yaml
executor: tofu  # or 'terraform'
```

### 2. Create Your First Stack

Create a stack file at `stacks/dev.stack.js`:

```javascript
// Define your stack
stack('dev', {
  project_name: 'my-app',
  region: 'us-central1'
})

// Configure backend
backend('gcs', {
  bucket: 'my-terraform-state',
  prefix: 'comet/{{ .stack }}/{{ .component }}'
})

// Define a simple component
const vpc = component('vpc', 'modules/vpc', {
  cidr_block: '10.0.0.0/16',
  region: '{{ .settings.region }}'
})
```

### 3. Create a Terraform Module

Create a basic module at `modules/vpc/main.tf`:

```hcl
variable "cidr_block" {
  type = string
}

variable "region" {
  type = string
}

resource "google_compute_network" "vpc" {
  name                    = "my-vpc"
  auto_create_subnetworks = false
}

output "id" {
  value = google_compute_network.vpc.id
}
```

### 4. Run Comet Commands

```bash
# List available stacks
comet list

# Plan changes
comet plan dev vpc

# Apply changes
comet apply dev vpc

# View outputs
comet output dev vpc
```

## Next Steps

- Review [configuration options](/docs/guides/configuration) for customizing Comet's behavior
- Learn about [stack configuration](/docs/guides/stacks) and organizing your infrastructure
- Explore [components and modules](/docs/guides/components)
- Understand [cross-stack references](/docs/guides/cross-stack-references)
- Set up [secrets management](/docs/guides/secrets-management)
- Review [best practices](/docs/advanced/best-practices)

## Getting Help

For detailed command documentation, use:

```bash
comet --help
comet <command> --help
```

Visit our [GitHub repository](https://github.com/moonwalker/comet) to report issues or contribute.
