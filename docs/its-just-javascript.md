# 💡 It's Just JavaScript!

## The Most Important Thing to Know About Comet

**Comet's configuration DSL is a superset of JavaScript.** This means:

✅ **You can use any JavaScript feature**
- Functions
- Variables
- Loops
- Conditionals
- Array methods (map, filter, reduce)
- String interpolation
- Destructuring
- Imports/exports
- Everything else!

✅ **You can create your own helper functions**
- Domain helpers
- Component factories
- Credential presets
- Tag templates
- Validation functions
- Anything you can imagine!

✅ **You're not limited to built-in features**
- Comet provides minimal primitives
- You build the abstractions you need
- No waiting for features to be added
- Complete control over your patterns

## Quick Example

```javascript
const opts = {
  org: 'mycompany',
  domain: 'example.io',
  region: 'us-east-1'
}

stack('production', { opts })

// ============================================================================
// CREATE YOUR OWN HELPERS - It's just JavaScript!
// ============================================================================

// Domain helper matching YOUR team's pattern
function subdomain(name) {
  return `${name}.{{ .stack }}.${opts.domain}`
}

// Component factory for YOUR infrastructure
function k8sApp(name, config = {}) {
  return component(name, 'modules/k8s-app', {
    namespace: opts.org,
    domain: subdomain(name),
    region: opts.region,
    replicas: 2,
    ...config  // Override defaults
  })
}

// Credential preset for YOUR providers
function setupAWS(role = 'default') {
  envs({
    AWS_ACCESS_KEY_ID: secret(`aws/${role}/access_key`),
    AWS_SECRET_ACCESS_KEY: secret(`aws/${role}/secret_key`),
    AWS_REGION: opts.region
  })
}

// Validation helper
function requireEnv(name) {
  const value = envs(name)
  if (!value) throw new Error(`Missing required env var: ${name}`)
  return value
}

// ============================================================================
// USE YOUR HELPERS
// ============================================================================

setupAWS('production')
requireEnv('DATABASE_URL')

const api = k8sApp('api', {
  replicas: 5,          // Override default
  database_url: secret('database/url')
})

const admin = k8sApp('admin', {
  replicas: 1
})

// Or deploy multiple apps with a loop
['service1', 'service2', 'service3'].forEach(name => {
  k8sApp(name, {
    image: `${opts.org}/${name}:latest`
  })
})

// Conditional infrastructure
if (opts.org === 'production') {
  component('backup', 'modules/backup', {
    schedule: '0 2 * * *'
  })
}
```

## What Comet Provides (Built-in)

These are **universal primitives** that work for everyone:

1. **`envs()`** - Environment variables (now accepts objects too!)
2. **`secret()`** - Secrets shorthand with configurable defaults
3. **`secrets()`** - Full secrets path (original)
4. **`component()`** - Define infrastructure components
5. **`stack()`** - Define stacks
6. **`backend()`** - Configure state backend
7. **`append()`** - Add raw Terraform/OpenTofu code
8. **Template syntax** - `{{ .stack }}`, `{{ .component }}`, etc.
9. **`state()`** - Cross-stack references (in templates)

## What You Build (Userland)

These are **your patterns** - create exactly what you need:

- 🛠️ **Domain helpers** matching your DNS structure
- 🛠️ **Component factories** for your common infrastructure
- 🛠️ **Credential presets** for your cloud providers
- 🛠️ **Tag templates** matching your tagging strategy
- 🛠️ **Validation functions** for your requirements
- 🛠️ **Utility functions** for your workflows
- 🛠️ **Data transformations** for your configs
- 🛠️ **Anything else** you can code in JavaScript!

## Why This Design?

### ✅ Benefits:
- **No waiting** - Build what you need immediately
- **No opinions** - Your patterns, not Comet's assumptions
- **Transparent** - You see exactly what your code does
- **Flexible** - Easy to modify as needs change
- **Simple core** - Comet stays minimal and maintainable
- **Full power** - All of JavaScript at your disposal

### ❌ Without this:
- Wait for features to be added
- Work around Comet's opinions
- Limited to what's built-in
- Fork and maintain your own version
- Complex feature requests
- Opinionated abstractions

## Where to Learn More

1. **[Userland Patterns](userland-patterns.md)** - Comprehensive guide with examples
2. **[DSL Improvements](dsl-improvements.md)** - New built-in features
3. **[Quick Reference](dsl-quick-reference.md)** - Syntax cheat sheet
4. **`stacks/_examples/userland-helpers.stack.js`** - Working examples

## TypeScript Support (Automatic)

Comet includes TypeScript definitions for **autocomplete and type checking** in your IDE - but you still write JavaScript!

### How It Works

TypeScript definitions can be generated on-demand using the `comet types` command. Comet embeds the definitions in the binary and writes `stacks/index.d.ts` when you run:

```bash
comet types
```

Then add this to your `.stack.js` files for IDE autocomplete:

```javascript
/// <reference path="./index.d.ts" />
```

Or enable type checking:

```javascript
// @ts-check
/// <reference path="./index.d.ts" />
```

### Benefits

- ✅ **Autocomplete** for all Comet functions
- ✅ **Inline documentation** as you type
- ✅ **Type checking** catches errors before running
- ✅ **No compilation** - still plain JavaScript
- ✅ **Better refactoring** in VS Code

### Example

```javascript
// @ts-check
/// <reference path="../types/index.d.ts" />

// Autocomplete suggests all options
secretsConfig({
  defaultProvider: 'sops',  // IDE knows valid values
  defaultPath: 'secrets.enc.yaml'
})

// JSDoc types for your helpers
/**
 * @param {string} name
 * @param {number} [replicas=2]
 */
function k8sApp(name, replicas = 2) {
  return component(name, 'modules/k8s-app', { replicas })
}

// TypeScript validates parameters
const api = k8sApp('api', 5)  // ✅ OK
// const bad = k8sApp()       // ❌ Error: missing name
```

See `stacks/_examples/with-typescript-support.stack.js` for a complete example.

## The Bottom Line

**You don't need to wait for Comet to add features. Just write JavaScript!**

If you find yourself thinking "I wish Comet had a feature for X", remember:
1. It's JavaScript - you can build it yourself
2. It'll be exactly what you need (not a compromise)
3. You can share it with your team
4. You're not dependent on Comet's release cycle

**Happy coding! 🚀**
