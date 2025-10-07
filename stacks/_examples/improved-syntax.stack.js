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
envs({
  DIGITALOCEAN_TOKEN: secrets('sops://secrets.enc.yaml#/digitalocean/token'),
  AWS_ACCESS_KEY_ID: secrets('sops://secrets.enc.yaml#/digitalocean/spaces_access_key'),
  AWS_SECRET_ACCESS_KEY: secrets('sops://secrets.enc.yaml#/digitalocean/spaces_secret_key'),
  CLOUDFLARE_API_TOKEN: secrets('sops://secrets.enc.yaml#/cloudflare/api_token')
})

// ============================================================================
// Feature 2: Secrets Path Shorthand
// ============================================================================

// Configure default secrets provider and path (optional, defaults shown)
secretsConfig({
  defaultProvider: 'sops',
  defaultPath: 'secrets.enc.yaml'
})

// Now you can use the shorthand secret() function

// OLD WAY (verbose):
const datadogOld = component('datadog-old', 'modules/datadog', {
  api_key: secrets('sops://secrets.enc.yaml#/datadog/api_key'),
  app_key: secrets('sops://secrets.enc.yaml#/datadog/app_key'),
  cluster_name: `${settings.common_name}-{{ .stack }}`
})

// NEW WAY (shorthand with forward slash):
const datadog = component('datadog', 'modules/datadog', {
  api_key: secret('datadog/api_key'),
  app_key: secret('datadog/app_key'),
  cluster_name: `${settings.common_name}-{{ .stack }}`
})

// ALSO WORKS (dot notation):
const argocd = component('argocd', 'modules/argocd', {
  admin_password: secret('argocd.admin_password'),
  slack_token: secret('argocd.slack_token'),
  github_client_secret: secret('argocd.github_client_secret')
})

// Still works with full path if needed
const legacy = component('legacy', 'modules/legacy', {
  token: secrets('sops://secrets.enc.yaml#/legacy/token')
})

// ============================================================================
// Feature 3: Domain Name Helpers
// ============================================================================

// OLD WAY (manual template strings):
const pgwebOld = component('pgweb-old', 'modules/pgweb', {
  domain_name: 'pgweb.{{ .stack }}.{{ .settings.domain_name }}'
  // Results in: pgweb.improved-demo.example.io
})

const argoOld = component('argo-old', 'modules/argocd', {
  domain_name: 'argo.{{ .stack }}.{{ .settings.domain_name }}'
  // Results in: argo.improved-demo.example.io
})

// NEW WAY (subdomain helper - includes stack):
const pgweb = component('pgweb', 'modules/pgweb', {
  domain_name: subdomain('pgweb')
  // Expands to: pgweb.{{ .stack }}.{{ .settings.domain_name }}
  // Results in: pgweb.improved-demo.example.io
})

const argo = component('argo', 'modules/argocd', {
  domain_name: subdomain('argo')
  // Expands to: argo.{{ .stack }}.{{ .settings.domain_name }}
  // Results in: argo.improved-demo.example.io
})

const natsNui = component('nats-nui', 'modules/nats-nui', {
  domain_name: subdomain('nui')
  // Results in: nui.improved-demo.example.io
})

// Use fqdn() for domains without stack prefix
const api = component('api', 'modules/api', {
  domain_name: fqdn('api')
  // Expands to: api.{{ .settings.domain_name }}
  // Results in: api.example.io (no stack in domain)
})

const www = component('www', 'modules/website', {
  domain_name: fqdn('www')
  // Results in: www.example.io
})

// Override stack name in subdomain if needed
const customStack = component('custom', 'modules/app', {
  domain_name: subdomain('app', { stack: 'production' })
  // Expands to: app.production.{{ .settings.domain_name }}
  // Results in: app.production.example.io (regardless of actual stack)
})

// ============================================================================
// Complete Example: All Features Together
// ============================================================================

const database = component('database', 'modules/database', {
  // Shorthand secrets with dot notation
  admin_password: secret('database.admin_password'),
  replication_password: secret('database.replication_password'),

  // Domain helper
  admin_ui_domain: subdomain('db-admin'),

  // Regular config
  storage_size: '100Gi',
  replicas: 3
})

const monitoring = component('monitoring', 'modules/monitoring', {
  // Shorthand secrets with forward slash
  slack_webhook: secret('monitoring/slack_webhook'),
  pagerduty_key: secret('monitoring/pagerduty_key'),

  // Multiple domain helpers
  grafana_domain: subdomain('grafana'),
  prometheus_domain: subdomain('prometheus'),
  alertmanager_domain: subdomain('alerts'),

  // Database reference
  database_url: database.connection_string
})

print('✅ New DSL features demonstrated successfully!')
print('  • Bulk environment variables')
print('  • Secret path shorthand (slash and dot notation)')
print('  • Domain helpers (subdomain and fqdn)')
