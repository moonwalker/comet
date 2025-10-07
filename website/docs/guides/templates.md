---
sidebar_position: 6
---

# Templates and Functions

Comet uses Go templates with additional functions to enable dynamic configuration. This powerful feature allows you to reference variables, use conditional logic, and access cross-stack state.

## Template Syntax

Templates use the `{{ }}` syntax:

```javascript
const vpc = component('vpc', 'modules/vpc', {
  name: 'vpc-{{ .stack }}',  // Becomes 'vpc-dev', 'vpc-prod', etc.
  region: '{{ .settings.region }}'
})
```

## Built-in Variables

### `.stack`

The current stack name:

```javascript
backend('gcs', {
  bucket: 'terraform-state',
  prefix: '{{ .stack }}/{{ .component }}'  // dev/vpc, prod/vpc, etc.
})
```

### `.component`

The current component name:

```javascript
const database = component('db', 'modules/cloudsql', {
  instance_name: '{{ .stack }}-{{ .component }}'  // dev-db, prod-db
})
```

### `.settings`

Access stack settings:

```javascript
stack('production', {
  project_id: 'my-gcp-project',
  region: 'us-central1',
  db_tier: 'db-n1-standard-2'
})

const db = component('db', 'modules/cloudsql', {
  project: '{{ .settings.project_id }}',
  region: '{{ .settings.region }}',
  tier: '{{ .settings.db_tier }}'
})
```

## Template Functions

### `state` - Cross-Stack References

Reference outputs from other stacks:

```javascript
const app = component('app', 'modules/application', {
  vpc_id: '{{ (state "infrastructure" "vpc").id }}',
  db_host: '{{ (state "data" "database").connection_name }}'
})
```

See the [Cross-Stack References](/docs/guides/cross-stack-references) page for more details.

### `secrets` - Encrypted Secrets

Access SOPS-encrypted secrets:

```javascript
const db = component('database', 'modules/cloudsql', {
  password: '{{ secrets "sops://secrets.enc.yaml#/database/password" }}'
})
```

See the [Secrets Management](/docs/guides/secrets-management) page for more details.

## Conditional Logic

Use `if`, `else`, and `end` for conditional values:

```javascript
const instance = component('vm', 'modules/compute', {
  // Different machine types per environment
  machine_type: '{{ if eq .stack "production" }}n1-standard-4{{ else }}n1-standard-1{{ end }}',
  
  // Enable features only in production
  deletion_protection: '{{ if eq .stack "production" }}true{{ else }}false{{ end }}'
})
```

### Comparison Operators

- `eq` - Equal to
- `ne` - Not equal to
- `lt` - Less than
- `le` - Less than or equal
- `gt` - Greater than
- `ge` - Greater than or equal

```javascript
const backup = component('backup', 'modules/backup', {
  // Retention days based on environment
  retention_days: '{{ if eq .stack "production" }}30{{ else if eq .stack "staging" }}7{{ else }}1{{ end }}'
})
```

## String Functions

### `printf` - String Formatting

```javascript
const bucket = component('storage', 'modules/gcs', {
  name: '{{ printf "%s-%s-data" .settings.project_id .stack }}'
  // Results in: myproject-dev-data
})
```

### `lower` / `upper` - Case Conversion

```javascript
const resource = component('resource', 'modules/generic', {
  name: '{{ .stack | lower }}',  // dev, staging, production
  label: '{{ .stack | upper }}'  // DEV, STAGING, PRODUCTION
})
```

### `replace` - String Replacement

```javascript
const name = component('service', 'modules/app', {
  // Replace underscores with hyphens
  service_name: '{{ .settings.app_name | replace "_" "-" }}'
})
```

### `trim` - Remove Whitespace

```javascript
const config = component('config', 'modules/config', {
  value: '{{ .settings.some_value | trim }}'
})
```

## List and Object Functions

### `join` - Join List Elements

```javascript
const firewall = component('firewall', 'modules/firewall', {
  // Join list of IPs
  source_ranges: '{{ .settings.allowed_ips | join "," }}'
})
```

### `split` - Split String

```javascript
// In stack settings
stack('dev', {
  regions_str: 'us-central1,us-east1,us-west1'
})

// Use in component (though JavaScript would be better for this)
const multi_region = component('app', 'modules/app', {
  regions: '{{ .settings.regions_str | split "," }}'
})
```

## Default Values

### `default` - Provide Fallback

```javascript
const vm = component('vm', 'modules/compute', {
  // Use default if not set
  machine_type: '{{ .settings.machine_type | default "n1-standard-1" }}',
  zone: '{{ .settings.zone | default "us-central1-a" }}'
})
```

## Nested Templates

Access nested settings:

```javascript
stack('production', {
  gcp: {
    project_id: 'my-project',
    region: 'us-central1'
  },
  aws: {
    region: 'us-west-2',
    account_id: '123456789'
  }
})

const gcp_resource = component('gke', 'modules/gke', {
  project: '{{ .settings.gcp.project_id }}',
  region: '{{ .settings.gcp.region }}'
})

const aws_resource = component('eks', 'modules/eks', {
  region: '{{ .settings.aws.region }}'
})
```

## Combining JavaScript and Templates

You can combine JavaScript logic with template strings:

```javascript
// JavaScript for complex logic
const environments = {
  dev: { size: 'small', replicas: 1 },
  staging: { size: 'medium', replicas: 2 },
  production: { size: 'large', replicas: 5 }
}

stack('production', {
  env_config: environments['production']
})

// Templates for dynamic values
const app = component('app', 'modules/k8s', {
  replicas: '{{ .settings.env_config.replicas }}',
  resources: {
    size: '{{ .settings.env_config.size }}'
  }
})
```

## Common Patterns

### Environment-Specific Configuration

```javascript
const db = component('database', 'modules/cloudsql', {
  tier: '{{ if eq .stack "production" }}db-n1-standard-4{{ else }}db-f1-micro{{ end }}',
  backup_enabled: '{{ if eq .stack "production" }}true{{ else }}false{{ end }}',
  high_availability: '{{ if eq .stack "production" }}REGIONAL{{ else }}ZONAL{{ end }}'
})
```

### Resource Naming

```javascript
const resources = component('app', 'modules/app', {
  // Pattern: project-environment-component
  name: '{{ printf "%s-%s-%s" .settings.project_name .stack .component }}',
  
  // Pattern: component-environment
  alt_name: '{{ .component }}-{{ .stack }}',
  
  // Pattern: ENVIRONMENT_COMPONENT
  env_var_name: '{{ printf "%s_%s" (.stack | upper) (.component | upper) }}'
})
```

### Labels and Tags

```javascript
const instance = component('vm', 'modules/compute', {
  labels: {
    environment: '{{ .stack }}',
    component: '{{ .component }}',
    managed_by: 'comet',
    project: '{{ .settings.project_name }}'
  }
})
```

## Best Practices

1. **Use JavaScript for Complex Logic** - Templates are great for simple substitutions, but use JavaScript for complex conditionals and transformations

2. **Keep Templates Readable** - Break complex templates into multiple lines or use JavaScript variables

3. **Validate Template Output** - Use `comet export` to see the generated Terraform files and verify template expansion

4. **Document Template Variables** - Add comments explaining what template variables are used and where they come from

```javascript
// Good: Clear and documented
stack('production', {
  // GCP configuration
  project_id: 'my-gcp-project',  // GCP project ID
  region: 'us-central1',          // Primary region
  
  // Database configuration
  db_tier: 'db-n1-standard-4',    // CloudSQL tier
  db_version: 'POSTGRES_14'       // PostgreSQL version
})
```
