// Example: Using SOPS with AGE encryption
//
// The SOPS_AGE_KEY is bootstrapped once and cached locally,
// allowing sops:// references to work immediately.
//
// Bootstrap setup (run once):
//   comet bootstrap add SOPS_AGE_KEY op://ci-cd/sops-age-key/private
//
// The key is cached in .comet/bootstrap.state and auto-loaded

const settings = {
  org: 'myorg',
  app: 'myapp'
}

stack('sops-demo', { settings })

metadata({
  description: 'SOPS with AGE encryption example',
  tags: ['example', 'sops', 'age', 'secrets']
})

backend('gcs', {
  bucket: 'my-terraform-state',
  prefix: `${settings.org}/${settings.app}/{{ .stack }}/{{ .component }}`
})

// Configure SOPS as default provider
secretsConfig({
  defaultProvider: 'sops',
  defaultPath: 'secrets/prod.yaml'
})

// SOPS_AGE_KEY is already available from bootstrap, so this works:
component('app', 'modules/app', {
  // These secrets are decrypted using the AGE key from bootstrap
  database_password: secrets('sops://secrets/prod.yaml#database_password'),
  api_key: secrets('sops://secrets/prod.yaml#api_key'),
  
  // Or use different file
  external_api_token: secrets('sops://secrets/external.yaml#api_token')
})
