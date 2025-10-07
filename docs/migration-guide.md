# Migration Guide: Upgrading to New DSL Features

This guide helps you migrate existing Comet stack files to use the new DSL improvements.

## Overview

The new features are **100% backward compatible**. You can:
- ✅ Keep all existing code as-is
- ✅ Adopt new features incrementally
- ✅ Mix old and new syntax in the same file
- ✅ Migrate at your own pace

## Quick Wins (Start Here)

These changes provide immediate value with minimal effort:

### 1. Add Secrets Configuration (1 line)

Add this near the top of your stack file:

```javascript
secretsConfig({
  defaultProvider: 'sops',
  defaultPath: 'secrets.enc.yaml'
})
```

Now you can use `secret()` instead of `secrets()` everywhere.

### 2. Replace Domain Template Strings

**Find:** `domain_name: 'something.{{ .stack }}.{{ .settings.domain_name }}'`
**Replace:** `domain_name: subdomain('something')`

**Find:** `domain_name: 'something.{{ .settings.domain_name }}'`
**Replace:** `domain_name: fqdn('something')`

## Step-by-Step Migration

### Before (Original Syntax)

```javascript
// stacks/production.stack.js
const settings = {
  org: 'mycompany',
  common_name: 'platform',
  domain_name: 'mycompany.io'
}

stack('production', { settings })

backend('gcs', {
  bucket: 'terraform-state',
  prefix: `${settings.org}/{{ .stack }}/{{ .component }}`
})

// Environment variables
envs('DIGITALOCEAN_TOKEN', secrets('sops://secrets.enc.yaml#/digitalocean/token'))
envs('AWS_ACCESS_KEY_ID', secrets('sops://secrets.enc.yaml#/digitalocean/spaces_access_key'))
envs('AWS_SECRET_ACCESS_KEY', secrets('sops://secrets.enc.yaml#/digitalocean/spaces_secret_key'))
envs('CLOUDFLARE_API_TOKEN', secrets('sops://secrets.enc.yaml#/cloudflare/api_token'))

// Database
const database = component('database', 'modules/postgres', {
  admin_password: secrets('sops://secrets.enc.yaml#/database/admin_password'),
  replication_password: secrets('sops://secrets.enc.yaml#/database/replication_password'),
  admin_ui_domain: 'pgweb.{{ .stack }}.{{ .settings.domain_name }}',
  storage_size: '100Gi'
})

// Monitoring
const monitoring = component('monitoring', 'modules/monitoring', {
  slack_webhook: secrets('sops://secrets.enc.yaml#/monitoring/slack_webhook'),
  pagerduty_key: secrets('sops://secrets.enc.yaml#/monitoring/pagerduty_key'),
  grafana_domain: 'grafana.{{ .stack }}.{{ .settings.domain_name }}',
  prometheus_domain: 'prometheus.{{ .stack }}.{{ .settings.domain_name }}',
  database_url: database.connection_string
})

// Public API
const api = component('api', 'modules/api', {
  api_key: secrets('sops://secrets.enc.yaml#/api/key'),
  domain_name: 'api.{{ .settings.domain_name }}',
  database_url: database.connection_string
})
```

### After (New Syntax)

```javascript
// stacks/production.stack.js
const settings = {
  org: 'mycompany',
  common_name: 'platform',
  domain_name: 'mycompany.io'
}

stack('production', { settings })

backend('gcs', {
  bucket: 'terraform-state',
  prefix: `${settings.org}/{{ .stack }}/{{ .component }}`
})

// Configure secrets once
secretsConfig({
  defaultProvider: 'sops',
  defaultPath: 'secrets.enc.yaml'
})

// Bulk environment variables
envs({
  DIGITALOCEAN_TOKEN: secret('digitalocean/token'),
  AWS_ACCESS_KEY_ID: secret('digitalocean/spaces_access_key'),
  AWS_SECRET_ACCESS_KEY: secret('digitalocean/spaces_secret_key'),
  CLOUDFLARE_API_TOKEN: secret('cloudflare/api_token')
})

// Database
const database = component('database', 'modules/postgres', {
  admin_password: secret('database/admin_password'),
  replication_password: secret('database/replication_password'),
  admin_ui_domain: subdomain('pgweb'),
  storage_size: '100Gi'
})

// Monitoring
const monitoring = component('monitoring', 'modules/monitoring', {
  slack_webhook: secret('monitoring/slack_webhook'),
  pagerduty_key: secret('monitoring/pagerduty_key'),
  grafana_domain: subdomain('grafana'),
  prometheus_domain: subdomain('prometheus'),
  database_url: database.connection_string
})

// Public API
const api = component('api', 'modules/api', {
  api_key: secret('api/key'),
  domain_name: fqdn('api'),
  database_url: database.connection_string
})
```

### What Changed?

1. ✅ Added `secretsConfig()` - 1 line
2. ✅ Converted 4 `envs()` calls to 1 object - Saved 3 lines
3. ✅ Converted 5 `secrets()` to `secret()` - Saved ~150 characters
4. ✅ Converted 3 domain templates to helpers - Saved ~100 characters

**Total reduction: ~30% less code with better readability**

## Find & Replace Patterns

Use these regex patterns in your editor:

### Pattern 1: Environment Variables Block

**Find (regex):**
```regex
envs\('([^']+)',\s*secrets\('sops://secrets\.enc\.yaml#/([^']+)'\)\)
```

**Manual conversion to:**
```javascript
envs({
  VAR_NAME: secret('path/to/secret'),
  // ... collect all matches
})
```

### Pattern 2: Secret References

**Find:** `secrets('sops://secrets.enc.yaml#/`
**Replace:** `secret('`

**Find:** `secrets('sops://secrets.enc.yaml#/([^']+)'\)`
**Replace:** `secret('$1')`

### Pattern 3: Subdomain Pattern

**Find:** `'([^']+)\.{{ \.stack }}\.{{ \.settings\.domain_name }}'`
**Replace:** `subdomain('$1')`

### Pattern 4: FQDN Pattern

**Find:** `'([^']+)\.{{ \.settings\.domain_name }}'`
**Replace:** `fqdn('$1')`

## Incremental Migration Strategy

### Phase 1: Add Configuration (5 minutes)
1. Add `secretsConfig()` to top of stack files
2. Test: `comet plan <stack>`

### Phase 2: Convert Secrets (10 minutes)
1. Find/replace `secrets('sops://secrets.enc.yaml#/` → `secret('`
2. Review changes
3. Test: `comet plan <stack>`

### Phase 3: Convert Domains (10 minutes)
1. Find/replace domain template strings with helpers
2. Review changes
3. Test: `comet plan <stack>`

### Phase 4: Convert Env Vars (10 minutes)
1. Group related `envs()` calls into objects
2. Review changes
3. Test: `comet plan <stack>`

### Total Migration Time: ~35 minutes per stack

## Rollback Plan

If you need to revert:

```bash
# Revert using git
git checkout HEAD -- stacks/your-stack.js

# Or manually undo changes - old syntax still works!
```

## Validation Checklist

After migration:

- [ ] `comet plan <stack>` runs without errors
- [ ] Generated Terraform files look correct
- [ ] Environment variables are set correctly
- [ ] Secrets resolve properly
- [ ] Domain names match expected values
- [ ] All components are present
- [ ] Backend configuration unchanged

## Common Issues

### Issue 1: Secrets Not Resolving

**Problem:** `secret('path')` returns undefined

**Solution:** Check that `secretsConfig()` is called before using `secret()`

### Issue 2: Wrong Domain Generated

**Problem:** Domain has extra/missing stack prefix

**Solution:** 
- Use `subdomain()` for environment-specific domains (includes stack)
- Use `fqdn()` for shared/production domains (no stack)

### Issue 3: Environment Variables Not Set

**Problem:** Env vars set with `envs({})` not available

**Solution:** Ensure object values are strings, not undefined

## Tips & Tricks

### Use Dot Notation for Nested Secrets

```javascript
// Both work the same:
secret('database/admin/password')
secret('database.admin.password')

// Use whichever feels more natural!
```

### Mix Old and New Syntax

```javascript
// Totally fine to mix:
envs({
  VAR1: secret('path1'),
  VAR2: secrets('sops://other-file.yaml#/path2')
})
```

### Shared Configuration

```javascript
// Create a shared config file
// stacks/_shared.js
export const secretsDefaults = {
  defaultProvider: 'sops',
  defaultPath: 'secrets.enc.yaml'
}

// Use in each stack
import { secretsDefaults } from './_shared.js'
secretsConfig(secretsDefaults)
```

## Need Help?

- 📖 Read [DSL Improvements](dsl-improvements.md) for complete documentation
- 📝 Check [DSL Quick Reference](dsl-quick-reference.md) for syntax
- 💡 See [examples](../stacks/_examples/) for working code
- 🐛 Found an issue? Open a GitHub issue

## Feedback Welcome!

These features are new. If you have suggestions or find issues:
1. Open a GitHub issue
2. Submit a PR with improvements
3. Share your experience

Happy migrating! 🚀
