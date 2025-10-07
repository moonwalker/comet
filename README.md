# Comet 

Cosmic tool for provisioning and managing infrastructure.

## Overview

Comet is a command-line interface (CLI) tool designed to streamline infrastructure provisioning and management. It provides a unified interface for handling infrastructure operations with modern tooling and practices.

## Features

- **🚀 JavaScript Configuration** - Define infrastructure using a familiar, powerful programming language instead of limited HCL
- **🔄 Automatic Backend Generation** - Comet generates `backend.tf.json` files automatically based on your stack configuration
- **🔗 Cross-Stack References** - Simple `state()` template function to reference outputs from other stacks
- **🔐 Built-in Secrets Management** - Native SOPS integration for encrypted secrets in your configurations
- **📦 Component Reusability** - Define components once, reuse across environments with different configurations
- **🎯 Multi-Environment Support** - Manage dev, staging, production, and any other environments from a single codebase
- **⚡ Terraform/OpenTofu Integration** - Works seamlessly with both Terraform and OpenTofu
- **🛠 Minimal Abstraction** - Thin wrapper that doesn't hide the underlying IaC tool
- **📝 YAML Configuration** - Simple `comet.yaml` for project-level settings
- **🌟 Template Support** - Use Go templates in component configurations for dynamic values

## Why Comet?

Comet fills the gap between plain Terraform/OpenTofu and heavy enterprise frameworks, offering a pragmatic solution for teams that need DRY infrastructure configurations without the overhead of complex tooling.

### The Challenge

While Terraform/OpenTofu are powerful, they have limitations when managing multi-environment infrastructures:

- ❌ **Backend configuration cannot be dynamic** - You must use partial configuration or wrappers
- ❌ **No native multi-environment patterns** - Workspaces aren't suitable for production isolation
- ❌ **Verbose .tfvars management** - Maintaining separate variable files for each environment
- ❌ **Complex cross-stack references** - Manual remote state data source configuration
- ❌ **Limited DRY capabilities** - HCL's declarative nature makes abstraction difficult

### How Comet Solves This

Comet provides:

- ✅ **JavaScript-based configuration** - Leverage a familiar, powerful language for infrastructure config
- ✅ **Automatic backend generation** - No more manual backend.tf files
- ✅ **Built-in cross-stack references** - Simple `state()` function for referencing other stacks
- ✅ **SOPS integration** - Native encrypted secrets support
- ✅ **Clean component reuse** - Share modules across environments with ease
- ✅ **Minimal abstraction** - Thin wrapper that doesn't hide Terraform/OpenTofu

### Comparison with Alternatives

| Feature | **Comet** | **Terragrunt** | **Atmos** | **Plain OpenTofu** |
|---------|-----------|----------------|-----------|-------------------|
| **Config Language** | JavaScript ✨ | HCL + YAML | YAML 📄 | HCL |
| **Learning Curve** | Moderate | Moderate | **Steep** | Low |
| **Backend Config** | ✅ Auto-generated | ✅ Native | ✅ Native | ❌ Manual |
| **Cross-Stack Refs** | ✅ `state()` function | ✅ Dependencies | ✅ Remote state | ⚠️ Manual setup |
| **Module Reuse** | ✅ JavaScript logic | ✅ Dependencies | ✅ Imports/Mixins | ⚠️ Copy-paste |
| **Secrets Management** | ✅ SOPS built-in | ❌ Bring your own | ❌ Bring your own | ❌ Manual |
| **Templating** | ✅ JS template literals | ⚠️ Functions | ⚠️ Go templates | ❌ Limited |
| **Community Size** | Small 🐭 | Large 🐘 | Medium 🐈 | Huge 🦕 |
| **Maturity** | Young | Very Mature | Mature | Stable |
| **Opinionation** | Low | Medium | **Very High** | Minimal |
| **Enterprise Features** | ❌ | ✅ | ✅✅✅ | ❌ |
| **Vendor Lock-in** | None | None | Cloud Posse | None |
| **Ideal For** | Small-Medium teams | Most teams | Large enterprises | Simple setups |

### When to Choose Comet

**Choose Comet if:**
- ✅ You have **< 50 components** across multiple environments
- ✅ Your team **prefers JavaScript** over YAML/HCL
- ✅ You want **minimal abstraction** and transparency
- ✅ You value **simplicity over extensive features**
- ✅ You need **built-in secrets management** (SOPS)
- ✅ You're comfortable maintaining a custom tool

**Consider alternatives if:**
- ⚠️ You need **enterprise governance** features (policy enforcement, compliance)
- ⚠️ You have **100+ components** across multiple orgs/regions
- ⚠️ You want the **most battle-tested** solution (Terragrunt)
- ⚠️ You need **Cloud Posse's reference architectures** (Atmos)

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

### `comet export`

**Description:** Export stack to standalone Terraform files.

**Usage:**
```
comet export <stack> [component] -o <output-dir>
```

## Advanced Examples

### Multi-Environment Setup

```javascript
// stacks/shared.js
const settings = {
  project_name: 'myapp',
  domain: 'example.com',
  gcp_project: 'my-gcp-project'
}

module.exports = { settings }
```

```javascript
// stacks/production.js
const { settings } = require('./shared.js')

stack('production', { settings })

backend('gcs', {
  bucket: 'my-terraform-state',
  prefix: `${settings.project_name}/{{ .stack }}/{{ .component }}`
})

const vpc = component('vpc', 'modules/vpc', {
  cidr_block: '10.0.0.0/16',
  region: 'us-central1'
})

const gke = component('gke', 'modules/gke', {
  network: vpc.id,
  cluster_name: `${settings.project_name}-prod`
})
```

### Cross-Stack References

Reference outputs from other stacks:

```javascript
// stacks/application.js
const webapp = component('webapp', 'modules/k8s-app', {
  // Reference infrastructure stack outputs
  cluster_endpoint: '{{ (state "infrastructure" "gke").endpoint }}',
  vpc_id: '{{ (state "infrastructure" "vpc").id }}'
})
```

### Dynamic Provider Configuration

```javascript
// Automatically generate provider configurations
append('providers', [
  `provider "google" {`,
  `  project = "${settings.gcp_project}"`,
  `  region  = "us-central1"`,
  `}`
])
```

### Secrets Management

```javascript
// Using SOPS for encrypted secrets
const db = component('database', 'modules/cloudsql', {
  password: secrets('sops://secrets.enc.yaml#/database/password'),
  admin_user: secrets('sops://secrets.enc.yaml#/database/admin_user')
})
```

For more examples, see the [docs](https://github.com/moonwalker/comet/tree/main/docs) directory.

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
