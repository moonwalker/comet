---
sidebar_position: 1
---

# Stacks

Stacks are the core organizational unit in Comet. A stack represents a complete environment (like dev, staging, or production) and contains one or more components.

## What is a Stack?

A stack is a JavaScript file that defines:
- The stack name and settings
- Backend configuration for Terraform state
- One or more infrastructure components
- Shared variables and configuration

## Creating a Stack

Create a stack file in the `stacks/` directory with the `.stack.js` or `.js` extension:

```javascript title="stacks/dev.stack.js"
// Define the stack with settings
stack('dev', {
  project_name: 'myapp',
  environment: 'development',
  region: 'us-central1'
})

// Configure the backend
backend('gcs', {
  bucket: 'my-terraform-state-bucket',
  prefix: 'comet/{{ .stack }}/{{ .component }}'
})

// Define components
const vpc = component('vpc', 'modules/vpc', {
  cidr_block: '10.0.0.0/16'
})
```

## Stack Settings

Settings are custom key-value pairs accessible throughout your stack using templates:

```javascript
stack('production', {
  project_name: 'myapp',
  domain: 'example.com',
  region: 'us-central1',
  db_tier: 'db-n1-standard-2'
})

// Use settings in components
const database = component('db', 'modules/cloudsql', {
  tier: '{{ .settings.db_tier }}',
  region: '{{ .settings.region }}'
})
```

## Shared Configuration

Create reusable configuration using JavaScript modules:

```javascript title="stacks/shared.js"
const settings = {
  project_name: 'myapp',
  domain: 'example.com',
  regions: {
    primary: 'us-central1',
    secondary: 'us-east1'
  }
}

module.exports = { settings }
```

Then import in your stacks:

```javascript title="stacks/production.js"
const { settings } = require('./shared.js')

stack('production', { 
  ...settings,
  environment: 'production'
})
```

## Backend Configuration

The `backend()` function configures where Terraform stores state:

### Google Cloud Storage

```javascript
backend('gcs', {
  bucket: 'my-terraform-state',
  prefix: 'comet/{{ .stack }}/{{ .component }}'
})
```

### AWS S3

```javascript
backend('s3', {
  bucket: 'my-terraform-state',
  key: 'comet/{{ .stack }}/{{ .component }}/terraform.tfstate',
  region: 'us-west-2'
})
```

### Local (for testing)

```javascript
backend('local', {
  path: '.terraform/{{ .stack }}/{{ .component }}/terraform.tfstate'
})
```

## Template Variables

Use template variables in your stack configuration:

- `{{ .stack }}` - Current stack name
- `{{ .component }}` - Current component name  
- `{{ .settings.KEY }}` - Access stack settings

```javascript
backend('gcs', {
  bucket: 'terraform-state-{{ .settings.project_name }}',
  prefix: '{{ .stack }}/{{ .component }}'
})
```

## Listing Stacks

View all available stacks:

```bash
comet list
```

View components in a specific stack:

```bash
comet list dev
```

## Multiple Environments

Create separate stack files for each environment:

```
stacks/
├── shared.js           # Shared configuration
├── dev.stack.js        # Development environment
├── staging.stack.js    # Staging environment
└── production.stack.js # Production environment
```

Each can have environment-specific configuration while sharing common settings.
