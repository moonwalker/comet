---
sidebar_position: 2
---

# Architecture

Understanding Comet's architecture helps you make the most of its features and troubleshoot issues effectively.

## Overview

Comet is a thin wrapper around Terraform/OpenTofu that enhances infrastructure-as-code workflows through JavaScript-based configuration and automatic code generation.

## Core Components

### 1. JavaScript Parser

**Responsibility:** Parse JavaScript stack files and convert them into Comet's internal schema

**Key Features:**
- Uses `goja` JavaScript runtime to execute stack files
- Provides global functions: `stack()`, `component()`, `backend()`, `append()`, `secrets()`
- Handles module resolution and `require()` statements
- Builds in-memory representation of infrastructure

**Example Flow:**
```javascript
// Stack file input
const vpc = component('vpc', 'modules/vpc', {
  cidr_block: '10.0.0.0/16'
})

// Comet creates internal Component struct:
// {
//   Name: "vpc",
//   Source: "modules/vpc",
//   Inputs: { cidr_block: "10.0.0.0/16" }
// }
```

### 2. Template Engine

**Responsibility:** Process Go templates in configuration values

**Key Features:**
- Supports Go template functions + Sprig library functions
- Special `state()` function for cross-stack references
- Special `secrets()` function for SOPS integration
- Template context includes stack name, component name, settings
- Handles nested template evaluation

**Template Functions:**
- `{{ .stack }}` - Current stack name
- `{{ .component }}` - Current component name
- `{{ .settings.KEY }}` - Access stack settings
- `{{ (state "stack" "component").output }}` - Cross-stack reference
- `{{ secrets "sops://..." }}` - Access encrypted secrets

### 3. Code Generator

**Responsibility:** Generate Terraform/OpenTofu configuration files

**Generated Files:**
- `backend.tf.json` - Backend configuration
- `providers_gen.tf` - Provider configurations
- `{stack}-{component}.tfvars.json` - Variable values
- Remote state data sources (for cross-stack refs)

**Generation Process:**
```
Stack Definition → Template Processing → Code Generation → Terraform Files
```

### 4. Terraform/OpenTofu Executor

**Responsibility:** Execute infrastructure operations

**Supported Operations:**
- `init` - Initialize Terraform/OpenTofu
- `plan` - Generate execution plan
- `apply` - Apply changes
- `destroy` - Destroy infrastructure
- `output` - Retrieve outputs

### 5. Secrets Manager

**Responsibility:** Handle encrypted secrets via SOPS

**Features:**
- Parse `sops://` URIs
- Decrypt SOPS files on-demand
- Extract values using JSON path syntax
- Support for various SOPS backends (age, GPG, cloud KMS)

## Data Flow

```
┌─────────────────────┐
│  Stack Files (.js)  │
│  - dev.js           │
│  - prod.js          │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│  JavaScript Parser  │
│  - Execute JS       │
│  - Build Schema     │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│  Template Engine    │
│  - Process {{ }}    │
│  - Cross-stack refs │
│  - Decrypt secrets  │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│  Code Generator     │
│  - backend.tf.json  │
│  - providers_gen.tf │
│  - *.tfvars.json    │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│ Terraform/OpenTofu  │
│  - init             │
│  - plan/apply       │
└─────────────────────┘
```

## Cross-Stack Reference Flow

When you use `{{ (state "infrastructure" "vpc").id }}`:

1. **Detection:** Template engine identifies `state()` function call
2. **State Configuration:** Reads backend config from infrastructure stack
3. **Remote State Block:** Generates `data "terraform_remote_state"` block
4. **Fallback Generation:** Creates safe fallback with `try()` function
5. **Variable Creation:** Creates local variable with the reference

Generated Terraform:
```hcl
data "terraform_remote_state" "infrastructure_vpc" {
  backend = "gcs"
  config = {
    bucket = "my-terraform-state"
    prefix = "infrastructure/vpc"
  }
}

locals {
  vpc_id = try(
    data.terraform_remote_state.infrastructure_vpc.outputs.id,
    null
  )
}
```

## Component Lifecycle

```
1. Define Component (stack.js)
   ↓
2. Parse & Validate
   ↓
3. Resolve Variables & Templates
   ↓
4. Decrypt Secrets (if any)
   ↓
5. Generate Terraform Files
   ↓
6. Initialize Terraform
   ↓
7. Execute Operation (plan/apply/destroy)
```

## Directory Structure (Runtime)

```
project/
├── stacks/                    # Your stack definitions
│   ├── shared.js
│   ├── dev.stack.js
│   └── production.stack.js
│
├── modules/                   # Terraform modules
│   ├── vpc/
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   └── outputs.tf
│   └── gke/
│
└── .terraform/                # Generated at runtime (gitignored)
    └── stacks/
        └── dev/
            └── vpc/
                ├── .terraform/           # Terraform cache
                ├── backend.tf.json       # Generated
                ├── providers_gen.tf      # Generated
                └── dev-vpc.tfvars.json   # Generated
```

## Design Principles

1. **Minimal Abstraction** - Don't hide Terraform/OpenTofu; generate standard files
2. **Explicit Configuration** - Make behavior clear and predictable
3. **JavaScript Flexibility** - Leverage JS for DRY configuration and logic
4. **File Generation** - Generate inspectable Terraform files
5. **State Compatibility** - Work with standard Terraform state

## Security Considerations

### Secrets Handling

- Secrets are decrypted in-memory during execution
- Never written to disk in plaintext
- SOPS integration requires proper key management
- Secrets can use various backends (age, GPG, cloud KMS)

### State Access

- Remote state access requires proper backend credentials
- Backend configuration uses the same credentials as Terraform
- State files may contain sensitive data

### Generated Files

- Generated files may contain sensitive values
- Always add `.terraform/` to `.gitignore`
- Consider generated files as temporary/ephemeral

### JavaScript Execution

- Stack files execute arbitrary JavaScript code
- Only run trusted stack files
- JavaScript runs in a sandboxed `goja` environment

## Performance Characteristics

### Parsing

- Fast JavaScript execution via `goja` VM
- Caching of parsed stack definitions
- Lazy evaluation where possible

### Template Processing

- Templates processed once per component
- Cross-stack state fetched on-demand
- SOPS files decrypted once per execution

### Parallel Execution

- Independent components can be applied in parallel
- Dependency graph determines execution order
- Terraform's parallel execution is preserved

### Caching

- Leverages Terraform's provider caching
- Module caching follows Terraform behavior
- State is always fetched fresh

## Extension Points

Comet can be extended in several ways:

### Custom Commands

Add new CLI commands by creating files in `cmd/`:

```go
// cmd/custom.go
package cmd

import "github.com/spf13/cobra"

var customCmd = &cobra.Command{
  Use:   "custom",
  Short: "Custom command",
  Run: func(cmd *cobra.Command, args []string) {
    // Implementation
  },
}
```

### Additional Backends

Extend backend support in the executor by adding new backend types.

### Custom Template Functions

Add functions in `internal/schema/templater.go`:

```go
funcMap["customFunc"] = func(arg string) string {
  // Custom logic
  return result
}
```

### Provider Integration

Add provider-specific helpers in component definitions or create custom modules.

## Comparison with Raw Terraform

| Aspect | Comet | Raw Terraform |
|--------|-------|---------------|
| **Configuration** | JavaScript | HCL |
| **Backend Setup** | Auto-generated | Manual configuration |
| **Cross-Stack Refs** | Built-in `state()` | Manual data sources |
| **Secrets** | SOPS integration | Manual external tools |
| **Multi-Environment** | Stack files | Workspaces or copy-paste |
| **State Files** | Per-component | Per workspace |
| **Complexity** | Higher learning curve | Simpler, more direct |
| **Flexibility** | JavaScript logic | HCL limitations |

## Debugging

### View Generated Files

```bash
# Export to see generated Terraform
comet export dev vpc -o ./exported

# Files will be in ./exported/
```

### Enable Verbose Logging

Set environment variables:
```bash
export TF_LOG=DEBUG
comet plan dev vpc
```

### Inspect State

```bash
# View outputs
comet output dev vpc

# Use Terraform directly on generated files
cd .terraform/stacks/dev/vpc
terraform state list
```
