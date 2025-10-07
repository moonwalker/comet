---
sidebar_position: 2
---

# Components

Components are individual pieces of infrastructure within a stack. Each component maps to a Terraform module and its configuration.

## What is a Component?

A component represents a single infrastructure resource or logical group of resources, such as:
- A VPC network
- A Kubernetes cluster
- A database instance
- A storage bucket
- A load balancer

## Defining Components

Use the `component()` function to define a component:

```javascript
const componentName = component('name', 'module-path', {
  // input variables
})
```

**Parameters:**
- `name` - Unique identifier for the component within the stack
- `module-path` - Path to the Terraform module (relative or absolute)
- `inputs` - Object containing variable values for the module

## Basic Example

```javascript title="stacks/dev.stack.js"
const vpc = component('vpc', 'modules/vpc', {
  cidr_block: '10.0.0.0/16',
  region: 'us-central1',
  enable_flow_logs: true
})
```

This creates a component named `vpc` that uses the module at `modules/vpc/` with the specified input variables.

## Component Dependencies

Components can reference each other within the same stack:

```javascript
const vpc = component('vpc', 'modules/vpc', {
  cidr_block: '10.0.0.0/16'
})

const subnet = component('subnet', 'modules/subnet', {
  vpc_id: vpc.id,  // Reference the VPC's output
  cidr_block: '10.0.1.0/24'
})

const gke = component('gke', 'modules/gke', {
  network: vpc.id,
  subnetwork: subnet.id,
  cluster_name: 'my-cluster'
})
```

When you reference `vpc.id`, Comet automatically creates a dependency, ensuring the VPC is created before the subnet and GKE cluster.

## Using Templates in Components

Template variables and functions can be used in component inputs:

```javascript
const database = component('db', 'modules/cloudsql', {
  name: 'db-{{ .stack }}',  // Becomes 'db-dev', 'db-prod', etc.
  region: '{{ .settings.region }}',
  tier: '{{ .settings.db_tier }}',
  
  // Use conditional logic with templates
  deletion_protection: '{{ if eq .stack "production" }}true{{ else }}false{{ end }}'
})
```

## Component Operations

### Plan Changes

Preview what changes will be made:

```bash
comet plan <stack> <component>
```

Example:
```bash
comet plan dev vpc
```

### Apply Changes

Create or update infrastructure:

```bash
comet apply <stack> <component>
```

Example:
```bash
comet apply dev vpc
```

### View Outputs

Display component outputs:

```bash
comet output <stack> <component>
```

Example:
```bash
comet output dev vpc
```

### Destroy Infrastructure

Remove infrastructure:

```bash
comet destroy <stack> <component>
```

Example:
```bash
comet destroy dev vpc
```

## Working with All Components

You can operate on all components in a stack by omitting the component name:

```bash
# Plan all components
comet plan dev

# Apply all components
comet apply dev

# Destroy all components
comet destroy dev
```

## Component Naming Best Practices

**DO:** Use descriptive, unique names
```javascript
component('vpc-main', 'modules/vpc', ...)
component('gke-primary-cluster', 'modules/gke', ...)
component('cloudsql-users-db', 'modules/database', ...)
```

**DON'T:** Use generic or cryptic names
```javascript
component('v1', 'modules/vpc', ...)      // ❌ Too generic
component('k', 'modules/gke', ...)        // ❌ Unclear
component('thing', 'modules/database', ...) // ❌ Not descriptive
```

## Module Structure

Components use standard Terraform modules. A basic module structure:

```
modules/vpc/
├── main.tf        # Main resource definitions
├── variables.tf   # Input variables
└── outputs.tf     # Output values
```

Example `modules/vpc/variables.tf`:
```hcl
variable "cidr_block" {
  type        = string
  description = "CIDR block for the VPC"
}

variable "region" {
  type        = string
  description = "Region for the VPC"
}
```

Example `modules/vpc/outputs.tf`:
```hcl
output "id" {
  value       = google_compute_network.vpc.id
  description = "VPC ID"
}

output "self_link" {
  value       = google_compute_network.vpc.self_link
  description = "VPC self link"
}
```

These outputs become available for other components to reference via `vpc.id` or `vpc.self_link`.
