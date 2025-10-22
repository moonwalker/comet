// Example: Complete secrets workflow
//
// This demonstrates the difference between:
// 1. Bootstrap secrets (one-time setup) - SOPS_AGE_KEY, provider credentials
// 2. Stack-level secrets - loaded DURING parsing, used by Terraform

// ============================================================
// Bootstrap Setup (run once):
// ============================================================
// comet bootstrap add SOPS_AGE_KEY op://ci-cd/sops-age-key/private
// comet bootstrap add DIGITALOCEAN_TOKEN op://production/digitalocean/token
//
// These are cached in .comet/bootstrap.state and auto-loaded
// ============================================================

const settings = {
  org: 'acme',
  app: 'myapp',
  env: 'production'
}

stack('complete-secrets', { settings })

metadata({
  description: 'Complete secrets workflow example',
  tags: ['example', 'secrets', 'sops', '1password']
})

backend('gcs', {
  bucket: 'my-terraform-state',
  prefix: `${settings.org}/${settings.app}/{{ .stack }}/{{ .component }}`
})

// Configure secrets defaults
secretsConfig({
  defaultProvider: 'sops',
  defaultPath: 'secrets/prod.yaml'
})

// Now that SOPS_AGE_KEY is bootstrapped, we can use sops:// references
component('database', 'modules/database', {
  // SOPS secret (requires SOPS_AGE_KEY from bootstrap)
  admin_password: secrets('sops://secrets/prod.yaml#admin_password'),

  // 1Password secret (loaded on-demand during stack parsing)
  backup_credentials: secrets('op://production/database/backup-key'),

  // Plain values work too
  database_name: `${settings.app}_${settings.env}`
})

component('app', 'modules/app', {
  // Mix and match secret sources
  api_key: secrets('sops://secrets/prod.yaml#api_key'),
  oauth_client_secret: secrets('op://production/oauth/client-secret'),

  // Reference outputs from other components
  database_host: state('database', 'host'),
  database_port: state('database', 'port')
})

// Bulk environment variables for Terraform
envs({
  TF_VAR_region: 'us-west-2',
  TF_VAR_environment: settings.env
})
