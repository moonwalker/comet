// Example: Complete secrets workflow
//
// This demonstrates the difference between:
// 1. Pre-loaded env vars (comet.yaml) - needed BEFORE parsing
// 2. Stack-level secrets - loaded DURING parsing, used by Terraform

// ============================================================
// In comet.yaml:
// ============================================================
// env:
//   # SOPS AGE key must be available before stack parsing
//   SOPS_AGE_KEY: op://ci-cd/sops-age-key/private
//
//   # Any other early-stage environment variables
//   TF_LOG: DEBUG
// ============================================================

stack({
  name: 'complete-secrets-example',
  backend: {
    type: 'gcs',
    bucket: 'my-terraform-state',
    prefix: 'complete-example'
  }
})

// Now that SOPS_AGE_KEY is set, we can use sops:// references
component('database', {
  source: './modules/database',
  vars: {
    // SOPS secret (requires SOPS_AGE_KEY from comet.yaml)
    admin_password: secret('sops://secrets/db.yaml#admin_password'),

    // 1Password secret (loaded on-demand during stack parsing)
    backup_credentials: secret('op://production/database/backup-key'),

    // Plain values work too
    database_name: 'myapp_production'
  }
})

component('app', {
  source: './modules/app',
  vars: {
    // Mix and match secret sources
    api_key: secret('sops://secrets/app.yaml#api_key'),
    oauth_client_secret: secret('op://production/oauth/client-secret'),

    // Reference outputs from other components
    database_host: state('database', 'host'),
    database_port: state('database', 'port')
  },
  envs: {
    // Environment variables for Terraform execution
    TF_VAR_region: 'us-west-2'
  }
})
