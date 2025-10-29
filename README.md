# Comet 

Cosmic tool for provisioning and managing infrastructure.

## Overview

Comet is a command-line interface (CLI) tool designed to streamline infrastructure provisioning and management. It provides a unified interface for handling infrastructure operations with modern tooling and practices.

## Features

- **🚀 JavaScript Configuration** - Define infrastructure using a familiar, powerful programming language instead of limited HCL
- **✨ Modern DSL Enhancements** - Bulk environment variables and secrets shorthand reduce boilerplate by ~30%
- **�️ Build Your Own Abstractions** - It's JavaScript! Create your own helper functions for domains, components, tags, etc.
- **�🔄 Automatic Backend Generation** - Comet generates `backend.tf.json` files automatically based on your stack configuration
- **🔗 Cross-Stack References** - Simple `state()` template function to reference outputs from other stacks
- **🔐 Built-in Secrets Management** - Native SOPS integration for encrypted secrets in your configurations
- **📦 Component Reusability** - Define components once, reuse across environments with different configurations
- **🎯 Multi-Environment Support** - Manage dev, staging, production, and any other environments from a single codebase
- **⚡ Terraform/OpenTofu Integration** - Works seamlessly with both Terraform and OpenTofu
- **�️ Minimal & Unopinionated** - Provides universal primitives, you build team-specific patterns
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

### Quick Install

```bash
curl -fsSL https://moonwalker.github.io/comet/install.sh | sh
```

### Prerequisites

- **OpenTofu** (recommended) or **Terraform**
- See [installation guide](docs/installation.md) for setup instructions

### Building from Source

```bash
git clone https://github.com/moonwalker/comet.git
cd comet
go build
```

## 💡 It's Just JavaScript!

**Important:** Comet's DSL is a superset of JavaScript. You can create any helper functions, abstractions, or patterns you need - you're not limited to built-in features!

👉 **[Read: It's Just JavaScript!](docs/its-just-javascript.md)** - Learn how to build your own domain helpers, component factories, and more.

## Usage

```
comet [command] [flags]
```

### Global Flags

- `--config` - Specify config file (default: `comet.yaml`)
- `--dir` - Specify stacks directory (default: `stacks`)

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

### `comet init`

**Description:** Initialize backends and providers without running plan/apply operations.

**Usage:**
```
comet init <stack> [component]
```

Useful for:
- Preparing infrastructure for read-only operations (like `comet output`)
- Troubleshooting provider/backend initialization issues
- CI/CD validation pipelines

### `comet output`

**Description:** Show output values from components.

**Usage:**
```
comet output <stack> [component]
```

### `comet list`

**Description:** List stacks or components with metadata.

**Usage:**
```bash
# List all stacks with basic metadata
comet list

# List all stacks with full metadata details
comet list --details

# List components in a specific stack
comet list <stack>
```

**Flags:**
- `--details, -d` - Show full metadata including owner

**Output:**
```
# Default view
+-------+-------------------------------------------+------------------------+---------------------+
| STACK |                DESCRIPTION                |          TAGS          |        PATH         |
+-------+-------------------------------------------+------------------------+---------------------+
| dev   | Development environment for Comet testing | dev, testing, non-prod | stacks/dev.stack.js |
+-------+-------------------------------------------+------------------------+---------------------+

# With --details flag
+-------+-------------------------------------------+---------------+------------------------+---------------------+
| STACK |                DESCRIPTION                |     OWNER     |          TAGS          |        PATH         |
+-------+-------------------------------------------+---------------+------------------------+---------------------+
| dev   | Development environment for Comet testing | platform-team | dev, testing, non-prod | stacks/dev.stack.js |
+-------+-------------------------------------------+---------------+------------------------+---------------------+
```

### `comet types`

**Description:** Generate TypeScript definitions for IDE support.

**Usage:**
```
comet types
```

Generates `index.d.ts` in your stacks directory, providing autocomplete and type hints when editing stack files in your IDE. Add `/// <reference path="./index.d.ts" />` to your stack files to enable IDE support.

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

### Stack Metadata

Add metadata to your stacks for better organization and documentation:

```javascript
// stacks/production.js
stack('production', { settings })

metadata({
  description: 'Production environment with HA configuration',
  owner: 'platform-team',
  tags: ['prod', 'critical', 'us-west'],
  custom: {
    slack_channel: '#prod-alerts',
    oncall: 'team-platform'
  }
})

backend('gcs', {
  bucket: 'prod-terraform-state',
  prefix: 'infra/{{ .stack }}/{{ .component }}'
})
```

**View metadata:**
```bash
# Show basic metadata
comet list

# Show full metadata including owner
comet list --details
```

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

Comet provides built-in SOPS integration with convenient shorthand syntax:

```javascript
// Configure secrets defaults (optional)
secretsConfig({
  defaultProvider: 'sops',
  defaultPath: 'secrets.enc.yaml'
})

// Old verbose syntax (still works)
const db1 = component('database', 'modules/cloudsql', {
  password: secrets('sops://secrets.enc.yaml#/database/password'),
  admin_user: secrets('sops://secrets.enc.yaml#/database/admin_user')
})

// New shorthand - much cleaner!
const db2 = component('database', 'modules/cloudsql', {
  password: secret('database/password'),  // or secret('database.password')
  admin_user: secret('database.admin_user')
})
```

### Modern DSL Enhancements

Comet includes productivity features that reduce boilerplate by ~30%:

**Bulk Environment Variables:**
```javascript
// Set multiple environment variables at once
envs({
  AWS_ACCESS_KEY_ID: secret('aws/access_key'),
  AWS_SECRET_ACCESS_KEY: secret('aws/secret_key'),
  CLOUDFLARE_API_TOKEN: secret('cloudflare/token')
})
```

**Build Your Own Abstractions:**

Since Comet's DSL is JavaScript, you can create any helper functions you need for your team's patterns:

```javascript
const opts = { domain: 'example.io' }

// Your own domain helpers
const subdomain = (name) => `${name}.{{ .stack }}.${opts.domain}`
const fqdn = (name) => `${name}.${opts.domain}`

// Your own component factories
function k8sApp(name, config) {
  return component(name, 'modules/k8s-app', {
    namespace: 'default',
    replicas: 2,
    domain: subdomain(name),
    ...config
  })
}

// Your own credential presets
function setupAWS() {
  envs({
    AWS_ACCESS_KEY_ID: secret('aws/access_key'),
    AWS_SECRET_ACCESS_KEY: secret('aws/secret_key'),
    AWS_REGION: 'us-east-1'
  })
}

// Use them!
setupAWS()
const api = k8sApp('api', { 
  replicas: 3,
  database_url: db.connection_string 
})
```

**Comet stays minimal and unopinionated. You build exactly what you need.**

See [DSL Improvements](docs/dsl-improvements.md) and [Userland Patterns](docs/userland-patterns.md) for complete documentation.

For more examples, see the [docs](https://github.com/moonwalker/comet/tree/main/docs) directory.

## Configuration

Comet can be configured using `comet.yaml` in your project directory.

### Basic Configuration

```yaml
# comet.yaml
stacks_dir: stacks              # Directory containing stack files
work_dir: stacks/_components    # Working directory for components
generate_backend: false         # Auto-generate backend.tf.json
log_level: INFO                 # Log verbosity
tf_command: tofu                # Use 'tofu' or 'terraform'
```

### Environment Variables

Set plain environment variables that are loaded before commands run:

```yaml
# comet.yaml
env:
  TF_LOG: DEBUG
  AWS_REGION: us-west-2
  PROJECT_ID: my-project
```

**Features:**
- Plain values only (fast startup)
- Shell environment variables take precedence
- Loaded before stack parsing begins

**Note:** The `env` section only supports plain values. For secrets, use the `bootstrap` feature below.

### Bootstrap: One-Time Secret Setup

Bootstrap fetches secrets from 1Password/SOPS and caches them locally. Run once, then all commands are fast!

```yaml
# comet.yaml
bootstrap:
  - name: sops-age-key
    type: secret
    source: op://vault/infrastructure/sops-age-key  # Your 1Password path
    target: ~/.config/sops/age/keys.txt
    mode: "0600"
```

**Usage:**

```bash
# One-time setup (takes ~4s)
comet bootstrap

# Check status
comet bootstrap status

# Now all commands are fast!
comet plan dev    # 100ms instead of 4s!
comet apply dev
```

**Features:**
- ✅ **One-time cost** - Slow fetch only happens once
- ✅ **Fast commands** - All subsequent commands use cached secrets
- ✅ **Idempotent** - Safe to run multiple times
- ✅ **State tracking** - Knows what's been set up
- ✅ **Standard location** - Uses SOPS's default key path

**Bootstrap Types:**

```yaml
bootstrap:
  # Fetch and cache secrets
  - name: sops-key
    type: secret
    source: op://vault/item/field
    target: ~/.config/sops/age/keys.txt
    mode: "0600"
  
  # Check required tools exist
  - name: check-tools
    type: check
    command: op,sops,tofu  # Comma-separated
  
  # Run custom commands
  - name: gcloud-auth
    type: command
    command: gcloud auth application-default login
    optional: true
```

**Migrating from v0.5.0:**

If you were using `op://` or `sops://` in the `env` section, migrate to `bootstrap`:

```yaml
# OLD (v0.5.0) - Slow on every command
env:
  SOPS_AGE_KEY: op://vault/key/private

# NEW (v0.6.0) - Fast after bootstrap
bootstrap:
  - name: sops-key
    type: secret
    source: op://vault/key/private
    target: ~/.config/sops/age/keys.txt
    mode: "0600"
``` 

**Recommended approach for frequently-used secrets:**

```bash
# Set in your shell (one-time cost)
export SOPS_AGE_KEY=$(op read "op://ci-cd/sops-age-key/private")

# Or add to ~/.bashrc, ~/.zshrc, ~/.config/fish/config.fish, etc.
```

Only use secret references in `comet.yaml` for:
- Secrets that change frequently
- CI/CD environments where shell setup is inconvenient
- Situations where the ~3-5s overhead is acceptable

See [Best Practices](docs/best-practices.md) for more configuration examples. 

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
