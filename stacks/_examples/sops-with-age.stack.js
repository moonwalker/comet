// Example: Using SOPS with AGE encryption
//
// The SOPS_AGE_KEY is pre-loaded from comet.yaml before this stack is parsed,
// allowing sops:// references to work immediately.
//
// In comet.yaml:
//   env:
//     SOPS_AGE_KEY: op://ci-cd/sops-age-key/private

stack({
  name: 'sops-example',
  backend: {
    type: 'gcs',
    bucket: 'my-terraform-state',
    prefix: 'sops-example'
  }
})

// SOPS_AGE_KEY is already available, so this works:
component('app', {
  source: './modules/app',
  vars: {
    // These secrets are decrypted using the AGE key set in comet.yaml
    database_password: secret('sops://secrets/prod.yaml#database_password'),
    api_key: secret('sops://secrets/prod.yaml#api_key')
  }
})
