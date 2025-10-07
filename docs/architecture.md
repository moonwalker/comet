# Comet Architecture

## Overview

Comet is a thin wrapper around Terraform/OpenTofu that enhances infrastructure-as-code workflows through JavaScript-based configuration and automatic code generation.

## Core Components

### 1. Stack Parser (`internal/parser/js/`)

**Responsibility:** Parse JavaScript stack files and convert them into Comet's internal schema

**Key Features:**
- Uses `goja` JavaScript runtime to execute stack files
- Provides global functions: `stack()`, `component()`, `backend()`, `append()`, `secrets()`
- Handles module resolution and `require()` statements
- Builds in-memory representation of infrastructure

**Example Flow:**
```javascript
// Stack file: stacks/dev.js
const vpc = component('vpc', 'modules/vpc', {
  cidr_block: '10.0.0.0/16'
})

// Comet internally creates a Component struct:
// {
//   Name: "vpc",
//   Source: "modules/vpc",
//   Inputs: { cidr_block: "10.0.0.0/16" }
// }
```

### 2. Templater (`internal/schema/templater.go`)

**Responsibility:** Process Go templates in configuration values

**Key Features:**
- Supports all Go template functions + Sprig functions
- Special `state()` function for cross-stack references
- Template context includes stack name, component name, settings
- Handles nested template evaluation

**Template Functions:**
- `{{ .stack }}` - Current stack name
- `{{ .settings.domain }}` - Access global settings
- `{{ (state "infra" "vpc").id }}` - Cross-stack reference

### 3. Code Generator

**Responsibility:** Generate Terraform/OpenTofu configuration files

**Generated Files:**
- `backend.tf.json` - Backend configuration
- `providers_gen.tf` - Provider configurations
- `{stack}-{component}.tfvars.json` - Variable values

**Generation Process:**
```
Stack Definition → Template Processing → Code Generation → Terraform Files
```

### 4. Executor (`internal/exec/tf/`)

**Responsibility:** Execute Terraform/OpenTofu commands

**Supported Operations:**
- `Plan` - Generate execution plan
- `Apply` - Apply changes
- `Destroy` - Destroy infrastructure
- `Output` - Retrieve outputs

### 5. Secrets Manager (`internal/secrets/`)

**Responsibility:** Handle encrypted secrets via SOPS

**Features:**
- Parse `sops://` URIs
- Decrypt SOPS files
- Extract values from encrypted YAML/JSON
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
│  - plan/apply       │
└─────────────────────┘
```

## Cross-Stack References

When you use `{{ (state "stack" "component").output }}`:

1. **Detection:** Templater identifies `state()` function call
2. **State Loading:** Executor fetches remote state for referenced component
3. **Fallback Generation:** Creates `try(data.terraform_remote_state.*.outputs.*, null)`
4. **Remote State Config:** Generates `data "terraform_remote_state"` block
5. **Local Variable:** Creates local variable with safe fallback

## Component Lifecycle

```
1. Define Component (stack.js)
   ↓
2. Parse & Validate
   ↓
3. Resolve Variables & Templates
   ↓
4. Generate Terraform Files
   ↓
5. Initialize Terraform
   ↓
6. Execute Operation (plan/apply)
```

## Directory Structure (Runtime)

```
project/
├── stacks/                    # Your stack definitions
│   ├── shared.js
│   ├── dev.js
│   └── production.js
├── modules/                   # Terraform modules
│   ├── vpc/
│   └── gke/
└── .terraform/                # Generated at runtime
    └── stacks/
        └── dev/
            └── vpc/
                ├── backend.tf.json          # Generated
                ├── providers_gen.tf         # Generated
                └── dev-vpc.tfvars.json      # Generated
```

## Extension Points

Comet can be extended in several ways:

### Custom Commands
Add new CLI commands in `cmd/`

### Additional Backends
Extend backend support in executor

### New Template Functions
Add custom template functions in `internal/schema/templater.go`

### Provider Integration
Add provider-specific helpers in component definitions

## Design Principles

1. **Minimal Abstraction** - Don't hide Terraform/OpenTofu
2. **Explicit Configuration** - Make behavior clear and predictable
3. **JavaScript Flexibility** - Leverage JS for DRY config
4. **File Generation** - Generate standard Terraform files (inspectable)
5. **State Compatibility** - Work with standard Terraform state

## Security Considerations

- **Secrets:** SOPS encryption for sensitive values
- **State Access:** Remote state access requires proper backend credentials
- **Generated Files:** May contain sensitive data, add to `.gitignore`
- **Stack Files:** Can execute arbitrary JavaScript (trust your stack files)

## Performance

- **Parsing:** Fast JavaScript execution via `goja`
- **Template Processing:** Lazy evaluation where possible
- **Parallel Execution:** Components can be applied in parallel if independent
- **Caching:** Leverages Terraform's provider and module caching
