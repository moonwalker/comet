// Example: Using SOPS with AGE encryption
//
// The SOPS_AGE_KEY is bootstrapped once and cached locally,
// allowing sops:// references to work immediately.
//
// Bootstrap setup (run once):
//   comet bootstrap add SOPS_AGE_KEY op://ci-cd/sops-age-key/private
//
// The key is cached in .comet/bootstrap.state and auto-loaded
//
// NOTE: This is a syntax example only. To run it, you would need:
// - SOPS-encrypted secrets/prod.yaml file
// - Bootstrap with SOPS_AGE_KEY configured

const settings = {
  org: 'myorg',
  app: 'myapp'
}

stack('sops-demo', { settings })

metadata({
  description: 'SOPS with AGE encryption example (syntax only)',
  tags: ['example', 'sops', 'age', 'secrets']
})

backend('gcs', {
  bucket: 'my-terraform-state',
  prefix: `${settings.org}/${settings.app}/{{ .stack }}/{{ .component }}`
})

// Example SOPS configuration and component
// Uncomment when you have actual SOPS-encrypted files
//
// secretsConfig({
//   defaultProvider: 'sops',
//   defaultPath: 'secrets/prod.yaml'
// })
//
// component('app', 'modules/app', {
//   // These secrets are decrypted using the AGE key from bootstrap
//   database_password: secrets('sops://secrets/prod.yaml#database_password'),
//   api_key: secrets('sops://secrets/prod.yaml#api_key'),
//
//   // Or use different file
//   external_api_token: secrets('sops://secrets/external.yaml#api_token')
// })
