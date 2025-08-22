# Cross-Stack References

Cross-stack references allow you to share resources and outputs between different stacks in Comet. This enables you to separate concerns, share infrastructure across environments, and maintain clean dependencies between different parts of your infrastructure.

## Overview

When components in one stack need to reference outputs from components in another stack, Comet automatically handles the complexity of setting up remote state data sources and managing dependencies.

## Syntax

Use the `state` template function to reference components from other stacks:

```javascript
{{ (state "STACK_NAME" "COMPONENT_NAME").OUTPUT_PROPERTY }}
```

Where:
- `STACK_NAME`: The name of the target stack
- `COMPONENT_NAME`: The name of the component within that stack
- `OUTPUT_PROPERTY`: The specific output property you want to reference

## Example

### Base Infrastructure Stack (`infra.stack.js`)

```javascript
const project = component('project', 'modules/gcp-project', {
  name: 'my-app-prod'
})

const gke = component('gke', 'modules/gke-cluster', {
  project_id: project.id,
  cluster_name: 'prod-cluster'
})
```

### Application Stack (`app.stack.js`)

```javascript
const webapp = component('webapp', 'modules/k8s-deployment', {
  // Reference outputs from the infra stack
  project_id: '{{ (state "infra" "project").id }}',
  cluster_endpoint: '{{ (state "infra" "gke").endpoint }}',
  cluster_ca_cert: '{{ (state "infra" "gke").ca_certificate }}'
})
```

## How It Works

When Comet encounters cross-stack references, it automatically:

1. **Generates remote state data sources** for the referenced components
2. **Creates safe fallbacks** using Terraform's `try()` function
3. **Configures the backend** using the same backend configuration with adjusted paths

For the example above, Comet would generate:

```hcl
data "terraform_remote_state" "project" {
  backend = "gcs"
  config = {
    bucket = "my-tf-state-bucket"
    prefix = "stacks/infra/project"
  }
}

locals {
  project_id = try(
    data.terraform_remote_state.project.outputs.id,
    null
  )
}
```

## Best Practices

1. **Separate Concerns**: Use different stacks for different layers (infrastructure, applications, monitoring)
2. **Minimize Dependencies**: Only reference what you actually need to avoid tight coupling
3. **Use Descriptive Names**: Make stack and component names clear and consistent
4. **Plan Deployment Order**: Ensure referenced stacks are deployed before dependent stacks

## Limitations

- Referenced components must exist and have been successfully deployed
- Cross-stack references create implicit dependencies that must be managed during deployment
- Backend configuration must be compatible between referencing and referenced stacks
