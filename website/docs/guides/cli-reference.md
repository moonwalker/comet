---
sidebar_position: 5
---

# CLI Commands Reference

Comet provides a set of commands for managing your infrastructure. This page documents all available commands and their options.

## Global Flags

These flags are available for all commands:

- `--help` - Display help information
- `--version` - Print version information

## comet version

Print the current version of Comet.

```bash
comet version
```

**Example Output:**
```
comet version 0.1.0
```

## comet list

List available stacks or components within a stack.

### List All Stacks

```bash
comet list
```

**Example Output:**
```
Available stacks:
  - dev
  - staging
  - production
```

### List Components in a Stack

```bash
comet list <stack>
```

**Example:**
```bash
comet list dev
```

**Example Output:**
```
Components in stack 'dev':
  - vpc
  - gke
  - cloudsql
  - redis
```

## comet types

Generate TypeScript definitions for IDE support.

```bash
comet types
```

**What it does:**
- Creates `index.d.ts` in your stacks directory
- Provides autocomplete and type hints in your IDE
- Enables type checking for your stack files

**Usage in stack files:**

Add this reference to the top of your `.stack.js` files:

```javascript
/// <reference path="./index.d.ts" />
```

Or enable type checking:

```javascript
// @ts-check
/// <reference path="./index.d.ts" />
```

**Benefits:**
- ✅ Autocomplete for all Comet functions
- ✅ Inline documentation
- ✅ Type checking catches errors before running
- ✅ Better refactoring support

See the [TypeScript Support](../guides/stacks.md#typescript-support) section for more details.

## comet plan

Show what changes will be made by the current configuration without actually applying them.

### Plan All Components

```bash
comet plan <stack>
```

**Example:**
```bash
comet plan production
```

### Plan a Specific Component

```bash
comet plan <stack> <component>
```

**Example:**
```bash
comet plan production vpc
```

**Output:**
Shows Terraform plan output with additions, changes, and deletions.

## comet init

Initialize backends and providers without running plan or apply operations. This is useful for setting up the environment before querying outputs or troubleshooting initialization issues.

### Initialize All Components

```bash
comet init <stack>
```

**Example:**
```bash
comet init dev
```

### Initialize a Specific Component

```bash
comet init <stack> <component>
```

**Example:**
```bash
comet init dev vpc
```

**What it does:**
- Generates backend configuration files (`backend.tf.json`)
- Generates provider configuration files (`providers_gen.tf`)
- Downloads required provider plugins
- Configures remote state backend
- Does NOT run plan, apply, or make infrastructure changes

**Use cases:**
- Before running `comet output` on a fresh checkout
- Troubleshooting provider/backend initialization issues
- CI/CD validation pipelines that only need to query infrastructure

## comet apply

Create or update infrastructure based on your configuration.

### Apply All Components

```bash
comet apply <stack>
```

**Example:**
```bash
comet apply dev
```

:::warning
This will apply changes to all components in the stack. Make sure to run `comet plan` first to review changes.
:::

### Apply a Specific Component

```bash
comet apply <stack> <component>
```

**Example:**
```bash
comet apply dev vpc
```

**Flags:**
- `--auto-approve` - Skip interactive approval (use with caution)

## comet output

Display output values from infrastructure components.

### Show All Outputs from a Component

```bash
comet output <stack> <component>
```

**Example:**
```bash
comet output production gke
```

**Example Output:**
```json
{
  "cluster_endpoint": "https://35.123.45.67",
  "cluster_ca_certificate": "LS0tLS1CRU...",
  "cluster_name": "prod-gke-cluster"
}
```

### Show All Outputs from All Components

```bash
comet output <stack>
```

**Example:**
```bash
comet output production
```

## comet destroy

Destroy previously-created infrastructure.

### Destroy All Components

```bash
comet destroy <stack>
```

**Example:**
```bash
comet destroy dev
```

:::danger
This will destroy ALL infrastructure in the stack. This action cannot be undone. Always review with `comet plan` first.
:::

### Destroy a Specific Component

```bash
comet destroy <stack> <component>
```

**Example:**
```bash
comet destroy dev vpc
```

**Flags:**
- `--auto-approve` - Skip interactive approval (use with extreme caution)

## comet clean

Delete Terraform-related folders and files (`.terraform`, state files, etc.) for a stack or component.

### Clean All Components

```bash
comet clean <stack>
```

**Example:**
```bash
comet clean dev
```

### Clean a Specific Component

```bash
comet clean <stack> <component>
```

**Example:**
```bash
comet clean dev vpc
```

:::tip
Use this command when you want to start fresh with Terraform initialization, or to clean up after destroying infrastructure.
:::

## comet export

Export stack configuration to standalone Terraform files for inspection or manual execution.

```bash
comet export <stack> [component] -o <output-dir>
```

**Example - Export Entire Stack:**
```bash
comet export production -o ./exported/production
```

**Example - Export Single Component:**
```bash
comet export production vpc -o ./exported/production-vpc
```

**Flags:**
- `-o, --output` - Output directory for exported files (required)

This generates standard Terraform files that can be used independently of Comet:
- `backend.tf.json` - Backend configuration
- `providers_gen.tf` - Provider configurations
- `*.tfvars.json` - Variable values
- Module source files (if applicable)

## comet kube

Generate kubeconfig for Kubernetes clusters.

```bash
comet kube <stack> <component>
```

**Example:**
```bash
comet kube production gke
```

This command generates a kubeconfig file that you can use to connect to your Kubernetes cluster:

```bash
# Use the generated kubeconfig
export KUBECONFIG=~/.kube/comet-production-gke
kubectl get nodes
```

## Common Workflows

### Initial Deployment

```bash
# 1. List available stacks
comet list

# 2. Review changes
comet plan dev

# 3. Apply infrastructure
comet apply dev

# 4. Check outputs
comet output dev
```

### Read-Only Query (Fresh Checkout)

```bash
# 1. Initialize backends and providers
comet init production

# 2. Query infrastructure outputs
comet output production vpc
comet output production gke
```

### Updating Infrastructure

```bash
# 1. Review what will change
comet plan production vpc

# 2. Apply the changes
comet apply production vpc

# 3. Verify outputs
comet output production vpc
```

### Destroying Infrastructure

```bash
# 1. Review what will be destroyed
comet plan production  # Should show deletions

# 2. Destroy infrastructure
comet destroy production

# 3. Clean up Terraform files
comet clean production
```

### Working with Specific Components

```bash
# Plan a single component
comet plan dev gke

# Apply a single component
comet apply dev gke

# Get component outputs
comet output dev gke

# Destroy a single component
comet destroy dev gke
```

## Tips

### Use Plan Before Apply

Always run `plan` before `apply` to review changes:

```bash
comet plan production
# Review the output carefully
comet apply production
```

### Component Dependencies

Comet respects component dependencies within a stack. When you apply a stack, components with dependencies are applied in the correct order.

### Parallel Execution

When applying multiple independent components, Comet can execute them in parallel for faster deployments.

### State Management

Comet uses the backend configuration defined in your stack files. Make sure your backend is properly configured and accessible before running commands.

### Error Recovery

If a command fails partway through:

1. Review the error message
2. Fix the issue in your configuration
3. Run the command again - Terraform's state management will resume from where it left off

:::danger Take care

This action is dangerous

:::
