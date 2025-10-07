---
sidebar_position: 3
---

# Cross-Stack References

Cross-stack references allow components in one stack to access outputs from components in another stack. This is essential for sharing infrastructure across environments and maintaining clean separation of concerns.

## Overview

When you need to reference infrastructure from another stack, Comet provides the `state()` template function that automatically handles remote state data sources and dependencies.

## Syntax

```javascript
'{{ (state "STACK_NAME" "COMPONENT_NAME").OUTPUT_PROPERTY }}'
```

**Parameters:**
- `STACK_NAME` - The name of the stack containing the component
- `COMPONENT_NAME` - The name of the component
- `OUTPUT_PROPERTY` - The output value you want to reference

## Basic Example

### Infrastructure Stack

First, create a base infrastructure stack:

```javascript title="stacks/infrastructure.stack.js"
stack('infrastructure', {
  project_id: 'my-gcp-project',
  region: 'us-central1'
})

backend('gcs', {
  bucket: 'my-terraform-state',
  prefix: 'comet/{{ .stack }}/{{ .component }}'
})

const project = component('project', 'modules/gcp-project', {
  name: 'my-app',
  project_id: '{{ .settings.project_id }}'
})

const vpc = component('vpc', 'modules/vpc', {
  project_id: project.id,
  cidr_block: '10.0.0.0/16'
})

const gke = component('gke', 'modules/gke', {
  project_id: project.id,
  network: vpc.id,
  cluster_name: 'main-cluster'
})
```

### Application Stack

Now reference the infrastructure from an application stack:

```javascript title="stacks/application.stack.js"
stack('application', {
  app_name: 'my-webapp'
})

backend('gcs', {
  bucket: 'my-terraform-state',
  prefix: 'comet/{{ .stack }}/{{ .component }}'
})

const webapp = component('webapp', 'modules/k8s-deployment', {
  // Reference outputs from the infrastructure stack
  project_id: '{{ (state "infrastructure" "project").id }}',
  cluster_endpoint: '{{ (state "infrastructure" "gke").endpoint }}',
  cluster_ca_cert: '{{ (state "infrastructure" "gke").ca_certificate }}',
  vpc_id: '{{ (state "infrastructure" "vpc").id }}'
})
```

## How It Works

When Comet encounters a `state()` reference, it automatically:

1. **Generates a remote state data source** for the referenced component
2. **Creates safe fallbacks** using Terraform's `try()` function
3. **Configures the backend** to match the referenced stack's backend

For the example above, Comet generates:

```hcl
data "terraform_remote_state" "infrastructure_gke" {
  backend = "gcs"
  config = {
    bucket = "my-terraform-state"
    prefix = "comet/infrastructure/gke"
  }
}

locals {
  cluster_endpoint = try(
    data.terraform_remote_state.infrastructure_gke.outputs.endpoint,
    null
  )
}
```

## Multiple References

You can reference multiple components from multiple stacks:

```javascript
const app = component('app', 'modules/application', {
  // Reference from infrastructure stack
  network: '{{ (state "infrastructure" "vpc").id }}',
  cluster: '{{ (state "infrastructure" "gke").name }}',
  
  // Reference from data stack
  database_host: '{{ (state "data" "cloudsql").connection_name }}',
  redis_host: '{{ (state "data" "redis").host }}',
  
  // Reference from monitoring stack
  prometheus_url: '{{ (state "monitoring" "prometheus").url }}'
})
```

## Environment-Specific References

Reference different stacks based on the environment:

```javascript title="stacks/app-dev.stack.js"
const app = component('app', 'modules/application', {
  // Reference dev infrastructure
  vpc_id: '{{ (state "infrastructure-dev" "vpc").id }}'
})
```

```javascript title="stacks/app-prod.stack.js"
const app = component('app', 'modules/application', {
  // Reference production infrastructure
  vpc_id: '{{ (state "infrastructure-prod" "vpc").id }}'
})
```

## Best Practices

### 1. Separate Infrastructure Layers

Organize stacks by lifecycle and responsibility:

```
stacks/
├── foundation.stack.js      # GCP project, IAM, etc.
├── networking.stack.js      # VPCs, subnets, firewall rules
├── kubernetes.stack.js      # GKE clusters
├── data.stack.js           # Databases, storage
└── applications.stack.js    # Applications and services
```

### 2. Minimize Cross-Stack Dependencies

Only reference what you actually need:

```javascript
// ✅ Good: Only reference specific outputs
const app = component('app', 'modules/app', {
  cluster_endpoint: '{{ (state "infra" "gke").endpoint }}'
})

// ❌ Bad: Don't create unnecessary dependencies
const app = component('app', 'modules/app', {
  // Don't reference things you don't use
  vpc_id: '{{ (state "infra" "vpc").id }}',
  random_output: '{{ (state "infra" "something").value }}'
})
```

### 3. Use Clear Naming

Make stack and component names descriptive:

```javascript
// ✅ Good: Clear what's being referenced
cluster: '{{ (state "infrastructure-prod" "gke-primary").endpoint }}'

// ❌ Bad: Unclear names
cluster: '{{ (state "infra" "k1").e }}'
```

### 4. Deploy in Order

Ensure referenced stacks are deployed before dependent stacks:

```bash
# Deploy in dependency order
comet apply foundation
comet apply networking
comet apply kubernetes
comet apply applications
```

## Common Patterns

### Shared VPC Pattern

```javascript title="stacks/shared-vpc.stack.js"
const vpc = component('shared-vpc', 'modules/vpc', {
  cidr_block: '10.0.0.0/16'
})
```

```javascript title="stacks/service-a.stack.js"
const service_a = component('api', 'modules/service', {
  vpc_id: '{{ (state "shared-vpc" "shared-vpc").id }}'
})
```

```javascript title="stacks/service-b.stack.js"
const service_b = component('worker', 'modules/service', {
  vpc_id: '{{ (state "shared-vpc" "shared-vpc").id }}'
})
```

### Multi-Environment Pattern

```javascript title="stacks/infra-prod.stack.js"
const gke = component('gke', 'modules/gke', {
  cluster_name: 'prod-cluster'
})
```

```javascript title="stacks/app-prod.stack.js"
const app = component('app', 'modules/app', {
  cluster: '{{ (state "infra-prod" "gke").endpoint }}'
})
```

## Limitations

- Referenced components must already exist and have been successfully applied
- Both stacks must use compatible backend configurations
- Changes to referenced outputs require re-planning dependent stacks
- Circular dependencies between stacks are not supported
