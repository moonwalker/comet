# DSL Improvements Documentation

This document describes the enhanced DSL features added to Comet to reduce boilerplate and improve developer experience.

> **💡 Important:** Comet's DSL is a **superset of JavaScript**. The features below are built-in conveniences, but you can create **any helper functions you want** for your team's specific patterns. See [Userland Patterns](#3-userland-patterns) below.

## Table of Contents

1. [Bulk Environment Variables](#1-bulk-environment-variables)
2. [Secrets Path Shorthand](#2-secrets-path-shorthand)
3. [Userland Patterns](#3-userland-patterns)

---

## 1. Bulk Environment Variables

### Problem
Setting multiple environment variables required repetitive function calls:

```javascript
envs('DIGITALOCEAN_TOKEN', secrets('sops://secrets.enc.yaml#/digitalocean/token'))
envs('AWS_ACCESS_KEY_ID', secrets('sops://secrets.enc.yaml#/digitalocean/spaces_access_key'))
envs('AWS_SECRET_ACCESS_KEY', secrets('sops://secrets.enc.yaml#/digitalocean/spaces_secret_key'))
envs('CLOUDFLARE_API_TOKEN', secrets('sops://secrets.enc.yaml#/cloudflare/api_token'))
```

### Solution
The `envs()` function now accepts an object/map to set multiple variables at once:

```javascript
envs({
  DIGITALOCEAN_TOKEN: secrets('sops://secrets.enc.yaml#/digitalocean/token'),
  AWS_ACCESS_KEY_ID: secrets('sops://secrets.enc.yaml#/digitalocean/spaces_access_key'),
  AWS_SECRET_ACCESS_KEY: secrets('sops://secrets.enc.yaml#/digitalocean/spaces_secret_key'),
  CLOUDFLARE_API_TOKEN: secrets('sops://secrets.enc.yaml#/cloudflare/api_token')
})
```

### Backward Compatibility
The original syntax still works:

```javascript
// Single key-value pair
envs('MY_VAR', 'my_value')

// Get environment variable
const token = envs('MY_VAR')
```

### Benefits
- **40% less code** - 4 function calls reduced to 1
- **Better grouping** - Related credentials are visually grouped
- **More readable** - Natural JavaScript object syntax

---

## 2. Secrets Path Shorthand

### Problem
Every secret reference required verbose, repetitive syntax:

```javascript
api_key: secrets('sops://secrets.enc.yaml#/datadog/api_key'),
app_key: secrets('sops://secrets.enc.yaml#/datadog/app_key'),
slack_token: secrets('sops://secrets.enc.yaml#/argocd/slack_token')
```

### Solution
New `secret()` shorthand function with configurable defaults:

#### Configuration (Optional)
```javascript
secretsConfig({
  defaultProvider: 'sops',      // Default: 'sops'
  defaultPath: 'secrets.enc.yaml'  // Default: 'secrets.enc.yaml'
})
```

#### Usage with Forward Slash
```javascript
api_key: secret('datadog/api_key'),
app_key: secret('datadog/app_key'),
slack_token: secret('argocd/slack_token')
```

#### Usage with Dot Notation
```javascript
api_key: secret('datadog.api_key'),
app_key: secret('datadog.app_key'),
slack_token: secret('argocd.slack_token')
```

Both notations are equivalent - dots are automatically converted to slashes.

### Full Path Override
If you need to use a different provider or path for a specific secret, use the full syntax:

```javascript
// Full syntax still works
token: secrets('sops://other-secrets.yaml#/special/token')

// Or from 1Password
password: secrets('op://vault/item/field')
```

### Backward Compatibility
The original `secrets()` function remains unchanged and fully functional.

### Benefits
- **50% reduction in characters** - Much more concise
- **Better readability** - Less visual noise
- **Easier refactoring** - Change default path in one place
- **Dot notation** - More natural for hierarchical secrets

### Examples

```javascript
// Configure once per stack
secretsConfig({
  defaultProvider: 'sops',
  defaultPath: 'secrets.enc.yaml'
})

// Database component
component('database', 'modules/postgres', {
  // OLD: secrets('sops://secrets.enc.yaml#/database/admin_password')
  admin_password: secret('database/admin_password'),
  
  // OLD: secrets('sops://secrets.enc.yaml#/database/replication_password')
  replication_password: secret('database.replication_password'),
  
  // Full syntax when needed
  backup_key: secrets('sops://backup-secrets.yaml#/db/backup_key')
})
```

---

## 3. Userland Patterns

### Philosophy

Instead of building every pattern into Comet, you can create your own helper functions in JavaScript. This keeps Comet minimal while giving you maximum flexibility.

**Why not build-in domain helpers, component groups, tag templates, etc.?**
- ❌ Too opinionated - not everyone uses the same patterns
- ❌ Creates maintenance burden - supporting edge cases
- ✅ JavaScript is powerful enough - you can build what you need
- ✅ Transparent - you see exactly what your helpers do
- ✅ Flexible - easy to modify for your specific needs

### Examples

**Your own domain helpers:**
```javascript
const opts = { base_domain: 'example.io' }

// Define your team's pattern
function subdomain(name) {
  return `${name}.{{ .stack }}.${opts.base_domain}`
}

function fqdn(name) {
  return `${name}.${opts.base_domain}`
}

// Use them
component('admin', 'modules/app', {
  domain_name: subdomain('admin')  // admin.dev.example.io
})

component('api', 'modules/api', {
  domain_name: fqdn('api')  // api.example.io
})
```

**Component factories:**
```javascript
function k8sApp(name, config) {
  return component(name, 'modules/k8s-app', {
    namespace: 'default',
    replicas: 2,
    ...config
  })
}

const api = k8sApp('api', {
  replicas: 3,
  database_url: db.connection_string
})
```

**Credential presets:**
```javascript
function setupDigitalOcean() {
  envs({
    DIGITALOCEAN_TOKEN: secret('digitalocean/token'),
    AWS_ACCESS_KEY_ID: secret('digitalocean/spaces_access_key'),
    AWS_SECRET_ACCESS_KEY: secret('digitalocean/spaces_secret_key')
  })
}

setupDigitalOcean()
```

See [Userland Patterns](userland-patterns.md) for comprehensive examples and best practices.

---

## Complete Example

Putting it all together:

```javascript
const settings = {
  org: 'mycompany',
  common_name: 'platform',
  domain_name: 'mycompany.io'
}

stack('staging', { settings })

backend('gcs', {
  bucket: 'terraform-state',
  prefix: `${settings.org}/{{ .stack }}/{{ .component }}`
})

// Configure secrets defaults
secretsConfig({
  defaultProvider: 'sops',
  defaultPath: 'secrets.enc.yaml'
})

// Bulk environment variables
envs({
  DIGITALOCEAN_TOKEN: secret('digitalocean/token'),
  AWS_ACCESS_KEY_ID: secret('aws/access_key'),
  AWS_SECRET_ACCESS_KEY: secret('aws/secret_key'),
  CLOUDFLARE_API_TOKEN: secret('cloudflare/api_token')
})

// Database with shorthand secrets
const database = component('database', 'modules/postgres', {
  admin_password: secret('database.admin_password'),
  admin_ui_domain: subdomain('pgweb'),
  replicas: 3
})

// Monitoring stack
component('monitoring', 'modules/monitoring', {
  slack_webhook: secret('monitoring/slack_webhook'),
  grafana_domain: subdomain('grafana'),
  prometheus_domain: subdomain('prometheus'),
  database_url: database.connection_string
})

// Public API (no stack in domain)
component('api', 'modules/api', {
  api_key: secret('api.key'),
  domain_name: fqdn('api'),
  database_url: database.connection_string
})
```

**Result:** ~40% less code with the same functionality!

---

## Migration Guide

### From Old Syntax

**Before:**
```javascript
envs('TOKEN', secrets('sops://secrets.enc.yaml#/my/token'))
envs('KEY', secrets('sops://secrets.enc.yaml#/my/key'))

component('app', 'modules/app', {
  password: secrets('sops://secrets.enc.yaml#/app/password'),
  domain_name: 'app.{{ .stack }}.{{ .settings.domain_name }}'
})
```

**After:**
```javascript
const opts = { base_domain: 'example.io' }

secretsConfig({
  defaultProvider: 'sops',
  defaultPath: 'secrets.enc.yaml'
})

envs({
  TOKEN: secret('my/token'),
  KEY: secret('my/key')
})

// Create your own domain helper if you want it
const subdomain = (name) => `${name}.{{ .stack }}.${opts.base_domain}`

component('app', 'modules/app', {
  password: secret('app/password'),
  domain_name: subdomain('app')
})
```

### Gradual Migration
All new features are **backward compatible**. You can:
1. Keep existing code as-is
2. Use new features in new components
3. Gradually refactor existing components

No breaking changes!

---

## Summary

| Feature | Code Reduction | Key Benefit |
|---------|---------------|-------------|
| Bulk Environment Variables | ~75% | Better grouping of related config |
| Secrets Path Shorthand | ~50% | Less visual noise, easier refactoring |
| Userland Patterns | N/A | Maximum flexibility, no opinions |

**Overall: ~30-40% less boilerplate** while maintaining clarity and flexibility.
