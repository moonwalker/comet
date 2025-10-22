// Example: Create your own custom helper functions in JavaScript
// This shows how to implement domain helpers (or any other patterns)
// without needing them built into Comet

const opts = {
  org: 'myorg',
  common_name: 'myapp',
  base_domain: 'example.io'  // Your domain structure, your choice
}

stack('demo', { opts })

metadata({
  description: 'Userland patterns: custom helpers and domain functions',
  tags: ['example', 'patterns', 'userland']
})

backend('gcs', {
  bucket: 'my-state',
  prefix: `${opts.org}/{{ .stack }}/{{ .component }}`
})

// ============================================================================
// USERLAND HELPERS - Define your own patterns!
// ============================================================================

// Example 1: Your own domain helpers (if you want them)
function subdomain(name) {
  // Your team's specific pattern
  return `${name}.{{ .stack }}.${opts.base_domain}`
}

function fqdn(name) {
  return `${name}.${opts.base_domain}`
}

// Example 2: Or use template strings directly (even simpler)
const domainFor = (name) => `${name}.{{ .stack }}.${opts.base_domain}`

// Example 3: Different domain pattern - preview branches
function previewDomain(name, branch) {
  return `${name}-${branch}.preview.${opts.base_domain}`
}

// Example 4: Component factory pattern
function k8sApp(name, config) {
  return component(name, 'modules/k8s-app', {
    namespace: 'default',
    domain: domainFor(name),
    ...config
  })
}

// ============================================================================
// USE YOUR HELPERS
// ============================================================================

// Use your custom helpers
const admin = k8sApp('admin', {
  replicas: 2,
  image: 'admin:latest'
})

const api = component('api', 'modules/api', {
  domain_name: fqdn('api'),  // api.example.io
  database_url: 'postgres://...'
})

const webapp = component('webapp', 'modules/frontend', {
  domain_name: subdomain('app'),  // app.demo.example.io
  api_url: api.endpoint
})

const preview = component('preview', 'modules/frontend', {
  domain_name: previewDomain('webapp', 'feature-123'),  // webapp-feature-123.preview.example.io
  api_url: api.endpoint
})

// ============================================================================
// OR JUST USE TEMPLATE STRINGS DIRECTLY
// ============================================================================

const monitoring = component('grafana', 'modules/grafana', {
  // No helper needed - just do it directly
  domain_name: `grafana.{{ .stack }}.${opts.base_domain}`,
  replicas: 1
})

print('✅ Custom userland helpers work perfectly!')
print('   No need to build every pattern into Comet.')
print('   JavaScript is flexible enough to create your own abstractions.')
