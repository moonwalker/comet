# DSL Quick Reference

> **💡 Remember:** Comet's DSL is JavaScript! You can create any helper functions you want. This reference shows built-in features and common user-created patterns.

---

## Bulk Environment Variables

```javascript
// Old way (still works)
envs('VAR1', 'value1')
envs('VAR2', 'value2')

// New way (bulk)
envs({
  VAR1: 'value1',
  VAR2: 'value2',
  VAR3: 'value3'
})
```

## Secrets Path Shorthand

```javascript
// Optional configuration (defaults shown)
secretsConfig({
  defaultProvider: 'sops',
  defaultPath: 'secrets.enc.yaml'
})

// Old way (still works)
secrets('sops://secrets.enc.yaml#/path/to/secret')

// New shorthand with forward slash
secret('path/to/secret')

// New shorthand with dot notation
secret('path.to.secret')

// Full path override when needed
secrets('sops://other-file.yaml#/special/secret')
secrets('op://vault/item/field')
```

## Userland Patterns

Create your own helpers for your team's patterns:

```javascript
// Your own domain helpers
const opts = { base_domain: 'example.io' }

function subdomain(name) {
  return `${name}.{{ .stack }}.${opts.base_domain}`
}

function fqdn(name) {
  return `${name}.${opts.base_domain}`
}

// Your own component factories
function k8sApp(name, config) {
  return component(name, 'modules/k8s-app', {
    namespace: 'default',
    replicas: 2,
    ...config
  })
}

// Your own credential presets
function setupDigitalOcean() {
  envs({
    DIGITALOCEAN_TOKEN: secret('digitalocean/token'),
    AWS_ACCESS_KEY_ID: secret('digitalocean/spaces_access_key'),
    AWS_SECRET_ACCESS_KEY: secret('digitalocean/spaces_secret_key')
  })
}
```

See [Userland Patterns](userland-patterns.md) for more examples.

## Complete Example

```javascript
const opts = {
  domain: 'example.io'
}

stack('dev', { opts })

// Add stack metadata
metadata({
  description: 'Development environment',
  owner: 'dev-team',
  tags: ['dev', 'testing']
})

// Configure secrets
secretsConfig({
  defaultProvider: 'sops',
  defaultPath: 'secrets.enc.yaml'
})

// Bulk env vars with shorthand secrets
envs({
  TOKEN: secret('provider/token'),
  KEY: secret('provider.key')
})

// Define your own helpers
const subdomain = (name) => `${name}.{{ .stack }}.${opts.domain}`

// Component with helpers
component('app', 'modules/app', {
  password: secret('app/password'),
  domain_name: subdomain('app')
})
```

## Migration Checklist

- [ ] Add `secretsConfig()` at top of stack file (optional)
- [ ] Replace multiple `envs()` calls with single object
- [ ] Replace `secrets()` with `secret()` for standard paths
- [ ] Create your own domain/component helpers if needed
- [ ] Test your stack: `comet plan <stack>`

## Tips

- All features are **backward compatible**
- Migrate gradually - old and new syntax can coexist
- Use `secret()` for most cases, `secrets()` for special cases
- Build your own helpers for your team's patterns
- Keep it simple - Comet stays minimal, you add what you need

