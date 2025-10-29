# Comet 

Minimal infrastructure orchestration tool for Terraform/OpenTofu with JavaScript configuration.

## Features

- **🚀 JavaScript Configuration** - Use a real programming language instead of HCL
- **🔄 Automatic Backend Generation** - No more manual backend.tf files
- **🔗 Cross-Stack References** - Simple `state()` function for referencing outputs
- **🔐 Built-in Secrets** - Native SOPS integration
- **📦 Multi-Environment** - DRY configurations across dev/staging/prod
- **⚡ Minimal & Unopinionated** - Thin wrapper, you build your own patterns

**It's Just JavaScript!** Create helper functions, abstractions, and patterns for your team's needs.

## Installation

```bash
curl -fsSL https://moonwalker.github.io/comet/install.sh | sh
```

**Prerequisites:** [OpenTofu](https://opentofu.org) or [Terraform](https://www.terraform.io)

See [installation guide](docs/installation.md) for details.

## Quick Start

**1. Create a stack file:**

```javascript
// stacks/dev.stack.js
stack('dev', { 
  project: 'myapp',
  region: 'us-central1' 
})

backend('gcs', {
  bucket: 'my-terraform-state',
  prefix: '{{ .stack }}/{{ .component }}'
})

component('vpc', 'modules/vpc', {
  cidr: '10.0.0.0/16'
})
```

**2. Run commands:**

```bash
comet plan dev      # Show changes
comet apply dev     # Create infrastructure
comet output dev    # Show outputs
comet destroy dev   # Tear down
```

## Documentation

- [Installation Guide](docs/installation.md)
- [Best Practices](docs/best-practices.md)
- [DSL Quick Reference](docs/dsl-quick-reference.md)
- [Cross-Stack References](docs/cross-stack-references.md)
- [It's Just JavaScript!](docs/its-just-javascript.md)
- [Examples](stacks/_examples/)

## Key Concepts

**Stacks** - Environments (dev, staging, prod)
```javascript
stack('production', { settings })
```

**Components** - Terraform modules
```javascript
const vpc = component('vpc', 'modules/vpc', { cidr: '10.0.0.0/16' })
```

**Cross-Stack References** - Reference other stacks
```javascript
cluster_id: '{{ (state "infrastructure" "gke").cluster_id }}'
```

**Secrets** - SOPS integration
```javascript
password: secret('database/password')
```

**Your Own Helpers** - It's JavaScript!
```javascript
function k8sApp(name, config) {
  return component(name, 'modules/k8s-app', {
    replicas: 2,
    ...config
  })
}
```

See [full documentation](https://github.com/moonwalker/comet/tree/main/docs) for details. 

## License

This project is licensed under the terms specified in the project's license file.
