// Example stack demonstrating the new DSL improvements:
// 1. Bulk Environment Variables
// 2. Secrets Path Shorthand
// 3. Domain Name Helpers

const settings = {
  org: 'myorg',
  common_name: 'myapp',
  domain_name: 'example.io'
}

stack('improved-demo', { settings })

metadata({
  description: 'DSL improvements: bulk envs, secret shortcuts, domain helpers',
  tags: ['example', 'dsl', 'features']
})

backend('gcs', {
  bucket: 'my-tf-state',
  prefix: `${settings.org}-${settings.common_name}/stacks/{{ .stack }}/{{ .component }}`
})

// ============================================================================
// Feature 1: Bulk Environment Variables
// ============================================================================

// OLD WAY (verbose):
// envs('DIGITALOCEAN_TOKEN', secrets('sops://secrets.enc.yaml#/digitalocean/token'))
// envs('AWS_ACCESS_KEY_ID', secrets('sops://secrets.enc.yaml#/digitalocean/spaces_access_key'))
// envs('AWS_SECRET_ACCESS_KEY', secrets('sops://secrets.enc.yaml#/digitalocean/spaces_secret_key'))
// envs('CLOUDFLARE_API_TOKEN', secrets('sops://secrets.enc.yaml#/cloudflare/api_token'))

// NEW WAY (bulk object syntax):
// NOTE: Commented out - requires actual secrets file
// envs({
//   DIGITALOCEAN_TOKEN: secrets('sops://secrets.enc.yaml#/digitalocean/token'),
//   AWS_ACCESS_KEY_ID: secrets('sops://secrets.enc.yaml#/digitalocean/spaces_access_key'),
//   AWS_SECRET_ACCESS_KEY: secrets('sops://secrets.enc.yaml#/digitalocean/spaces_secret_key'),
//   CLOUDFLARE_API_TOKEN: secrets('sops://secrets.enc.yaml#/cloudflare/api_token')
// })

// Example with plain values
envs({
  EXAMPLE_VAR_1: 'value1',
  EXAMPLE_VAR_2: 'value2',
  REGION: 'us-west-2'
})

// ============================================================================
// Feature 2: Secrets Path Shorthand
// ============================================================================

// Configure default secrets provider and path (optional, defaults shown)
// NOTE: Commented out - for demonstration only
// secretsConfig({
//   defaultProvider: 'sops',
//   defaultPath: 'secrets.enc.yaml'
// })

// Now you can use the shorthand secret() function

// OLD WAY (verbose):
// const datadogOld = component('datadog-old', 'modules/datadog', {
//   api_key: secrets('sops://secrets.enc.yaml#/datadog/api_key'),
//   app_key: secrets('sops://secrets.enc.yaml#/datadog/app_key'),
//   cluster_name: `${settings.common_name}-{{ .stack }}`
// })

// NEW WAY (shorthand with forward slash):
// const datadog = component('datadog', 'modules/datadog', {
//   api_key: secret('datadog/api_key'),
//   app_key: secret('datadog/app_key'),
//   cluster_name: `${settings.common_name}-{{ .stack }}`
// })

// ALSO WORKS (dot notation):
// const argocd = component('argocd', 'modules/argocd', {
//   admin_password: secret('argocd.admin_password'),
//   slack_token: secret('argocd.slack_token'),
//   github_client_secret: secret('argocd.github_client_secret')
// })

// Still works with full path if needed
// const legacy = component('legacy', 'modules/legacy', {
//   token: secrets('sops://secrets.enc.yaml#/legacy/token')
// })

// ============================================================================
// Feature 3: Template Strings for Domains
// ============================================================================

// Use template strings with placeholders for dynamic values:
const pgweb = component('pgweb', 'modules/pgweb', {
  domain_name: 'pgweb.{{ .stack }}.{{ .settings.domain_name }}'
  // Results in: pgweb.improved-demo.example.io
})

const argo = component('argo', 'modules/argocd', {
  domain_name: 'argo.{{ .stack }}.{{ .settings.domain_name }}'
  // Results in: argo.improved-demo.example.io
})

// For domains without stack prefix:
const api = component('api', 'modules/api', {
  domain_name: 'api.{{ .settings.domain_name }}'
  // Results in: api.example.io (no stack in domain)
})

// You can also create your own helper functions (see userland-helpers.stack.js)

// ============================================================================
// Complete Example: All Features Together
// ============================================================================

const database = component('database', 'modules/database', {
  // Template for domain
  admin_ui_domain: 'db-admin.{{ .stack }}.{{ .settings.domain_name }}',

  // Regular config
  storage_size: '100Gi',
  replicas: 3
})

const monitoring = component('monitoring', 'modules/monitoring', {
  // Multiple template domains
  grafana_domain: 'grafana.{{ .stack }}.{{ .settings.domain_name }}',
  prometheus_domain: 'prometheus.{{ .stack }}.{{ .settings.domain_name }}',
  alertmanager_domain: 'alerts.{{ .stack }}.{{ .settings.domain_name }}',

  // Database reference
  database_url: database.connection_string
})

print('✅ New DSL features demonstrated successfully!')
print('  • Bulk environment variables')
print('  • Secret path shorthand (slash and dot notation)')
print('  • Domain helpers (subdomain and fqdn)')
