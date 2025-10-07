# Userland Patterns

> **💡 Key Concept:** Comet's configuration DSL is a **superset of JavaScript**. This means you have the full power of JavaScript to create any helper functions, abstractions, or patterns you need. **You are not limited to built-in features!**

## Philosophy

Comet provides **minimal, universal primitives** and lets you build your own abstractions in JavaScript. This keeps the core simple while giving you maximum flexibility.

**Built into Comet:**
- ✅ Bulk environment variables (`envs({})`)
- ✅ Secrets management with shorthand (`secret()`)
- ✅ Template system (`{{ .stack }}`, `{{ .component }}`)
- ✅ Component system
- ✅ Cross-stack references

**You build yourself:**
- 🛠️ Domain helpers (your domain pattern ≠ everyone's pattern)
- 🛠️ Component factories (your infrastructure patterns)
- 🛠️ Credential presets (your provider setup)
- 🛠️ Tag templates (your tagging strategy)

## Why This Approach?

### ❌ Don't build into Comet:
```javascript
// Too opinionated - assumes everyone uses this pattern
subdomain('pgweb')  // → 'pgweb.{{ .stack }}.{{ .settings.domain_name }}'
```

### ✅ Do build yourself:
```javascript
// Your team's pattern, defined once, used everywhere
function subdomain(name) {
  return `${name}.{{ .stack }}.${opts.base_domain}`
}
```

**Benefits:**
- **No assumptions** - Works for everyone's patterns
- **Transparent** - You see exactly what it does
- **Flexible** - Easy to modify for your needs
- **Simple core** - Comet stays minimal and maintainable

## Common Patterns

### 1. Domain Helpers

Different teams structure domains differently:

```javascript
// Pattern 1: Stack-based subdomains
const opts = { base_domain: 'example.io' }
const subdomain = (name) => `${name}.{{ .stack }}.${opts.base_domain}`
// Usage: subdomain('api') → 'api.dev.example.io'

// Pattern 2: Environment-based
const env = 'staging'
const domain = (name) => `${name}-${env}.example.io`
// Usage: domain('api') → 'api-staging.example.io'

// Pattern 3: Preview branches
const previewDomain = (name, branch) => `${name}-${branch}.preview.example.io`
// Usage: previewDomain('app', 'feat-123') → 'app-feat-123.preview.example.io'

// Pattern 4: Different per service
const adminDomain = (name) => `${name}-admin.internal.example.io`
const publicDomain = (name) => `${name}.example.io`
```

### 2. Component Factories

Create reusable component builders:

```javascript
// Kubernetes app with sensible defaults
function k8sApp(name, config) {
  return component(name, 'modules/k8s-app', {
    namespace: 'default',
    replicas: 2,
    domain: `${name}.{{ .stack }}.${opts.domain}`,
    ...config  // Override defaults
  })
}

// Usage
const api = k8sApp('api', {
  replicas: 3,  // Override
  database_url: db.connection_string
})

// Database with monitoring
function database(name, size) {
  const db = component(name, 'modules/postgres', {
    storage_size: size,
    admin_password: secret(`${name}/admin_password`)
  })
  
  component(`${name}-pgweb`, 'modules/pgweb', {
    database_url: db.connection_string,
    domain: `${name}-admin.{{ .stack }}.${opts.domain}`
  })
  
  return db
}

// Usage
const mainDb = database('main-db', '100Gi')
```

### 3. Credential Presets

Your provider setup, your way:

```javascript
// Digital Ocean credentials
function setupDigitalOcean() {
  envs({
    DIGITALOCEAN_TOKEN: secret('digitalocean/token'),
    AWS_ACCESS_KEY_ID: secret('digitalocean/spaces_access_key'),
    AWS_SECRET_ACCESS_KEY: secret('digitalocean/spaces_secret_key')
  })
}

// GCP credentials
function setupGCP(project) {
  envs({
    GOOGLE_PROJECT: project,
    GOOGLE_CREDENTIALS: secret('gcp/credentials')
  })
}

// AWS with multiple roles
function setupAWS(role) {
  const basePath = `aws/${role}`
  envs({
    AWS_ACCESS_KEY_ID: secret(`${basePath}/access_key`),
    AWS_SECRET_ACCESS_KEY: secret(`${basePath}/secret_key`),
    AWS_REGION: 'us-east-1'
  })
}

// Usage
setupDigitalOcean()
setupGCP('my-project-123')
setupAWS('developer')
```

### 4. Tag Templates

Your tagging strategy:

```javascript
// Standard tags for all resources
function commonTags(additional = {}) {
  return {
    environment: '{{ .stack }}',
    managed_by: 'comet',
    team: opts.team_name,
    cost_center: opts.cost_center,
    ...additional
  }
}

// Usage
component('app', 'modules/app', {
  tags: commonTags({ service: 'api', tier: 'backend' })
})
```

### 5. Multi-Component Stacks

Deploy related components together:

```javascript
// Data stack: database + admin UI + backups
function dataStack(name) {
  const db = component(`${name}-db`, 'modules/postgres', {
    admin_password: secret(`${name}/db_password`)
  })
  
  component(`${name}-pgweb`, 'modules/pgweb', {
    database_url: db.connection_string,
    domain: `${name}-admin.{{ .stack }}.${opts.domain}`
  })
  
  component(`${name}-backup`, 'modules/backup', {
    database_url: db.connection_string,
    s3_bucket: `${opts.org}-${name}-backups`
  })
  
  return db
}

// Monitoring stack: Prometheus + Grafana + Alertmanager
function monitoringStack() {
  const prometheus = component('prometheus', 'modules/prometheus', {
    domain: `metrics.{{ .stack }}.${opts.domain}`
  })
  
  component('grafana', 'modules/grafana', {
    domain: `grafana.{{ .stack }}.${opts.domain}`,
    datasource_url: prometheus.url
  })
  
  component('alertmanager', 'modules/alertmanager', {
    slack_webhook: secret('monitoring/slack_webhook')
  })
}
```

## Organizing Shared Helpers

### Option 1: Shared JS File

```javascript
// stacks/_shared/helpers.js
export function subdomain(name) {
  return `${name}.{{ .stack }}.example.io`
}

export function k8sApp(name, config) {
  return component(name, 'modules/k8s-app', {
    namespace: 'default',
    ...config
  })
}

// stacks/dev.stack.js
import { subdomain, k8sApp } from './_shared/helpers.js'

const api = k8sApp('api', {
  domain: subdomain('api')
})
```

### Option 2: Inline in Each Stack

```javascript
// stacks/dev.stack.js
const opts = { domain: 'example.io' }

// Define helpers at top
const subdomain = (name) => `${name}.{{ .stack }}.${opts.domain}`
const k8sApp = (name, cfg) => component(name, 'modules/k8s-app', cfg)

// Use below
const api = k8sApp('api', { domain: subdomain('api') })
```

### Option 3: Generated from Config

```javascript
// stacks/config.js (imported by all stacks)
export const config = {
  domain: 'example.io',
  org: 'myorg'
}

export const helpers = {
  subdomain: (name) => `${name}.{{ .stack }}.${config.domain}`,
  orgBucket: (name) => `${config.org}-${name}`
}
```

## Best Practices

### ✅ Do:
- Create helpers for **your** repeated patterns
- Keep helpers simple and transparent
- Define helpers per-project or per-stack
- Use JavaScript's full power (loops, conditionals, etc.)

### ❌ Don't:
- Try to make helpers "universal" for all use cases
- Hide important configuration details
- Create deep abstraction hierarchies
- Reinvent frameworks inside Comet

## Examples

See complete examples in:
- `stacks/_examples/userland-helpers.stack.js`
- `stacks/_examples/component-factories.stack.js`

## Summary

**Comet provides minimal primitives. You build the abstractions you need.**

This keeps:
- Comet simple and maintainable
- Your code transparent and flexible
- Your team's patterns explicit
- Everyone happy 🎉
