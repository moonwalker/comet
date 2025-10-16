---
sidebar_position: 4
---

# Secrets Management

Comet provides built-in support for managing sensitive data using SOPS (Secrets OPerationS), allowing you to encrypt secrets and safely store them in version control.

## Overview

SOPS is a tool that encrypts YAML and JSON files using various key providers (age, GPG, AWS KMS, GCP KMS, Azure Key Vault). Comet integrates SOPS seamlessly, allowing you to reference encrypted secrets directly in your stack configurations.

## Setting Up SOPS

### 1. Install SOPS

```bash
# macOS
brew install sops

# Linux
wget https://github.com/mozilla/sops/releases/download/v3.8.1/sops-v3.8.1.linux
chmod +x sops-v3.8.1.linux
sudo mv sops-v3.8.1.linux /usr/local/bin/sops
```

### 2. Choose an Encryption Method

#### Using age (Recommended for simplicity)

```bash
# Install age
brew install age  # macOS
# or
apt install age   # Linux

# Generate a key pair
age-keygen -o key.txt

# Export the public key
export SOPS_AGE_RECIPIENTS=$(grep 'public key:' key.txt | cut -d ' ' -f 3)
```

#### Using GPG

```bash
# Generate a GPG key
gpg --generate-key

# Get your key fingerprint
gpg --list-keys
```

### 3. Create SOPS Configuration

Create `.sops.yaml` in your project root:

```yaml title=".sops.yaml"
creation_rules:
  - path_regex: secrets.*\.yaml$
    age: age1xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
    # Or for GPG:
    # pgp: YOUR_GPG_FINGERPRINT
```

### 4. Bootstrap SOPS Keys with Comet

If your SOPS encryption key is stored in 1Password, use Comet's bootstrap feature to set it up:

```yaml title="comet.yaml"
bootstrap:
  - name: sops-age-key
    type: secret
    source: op://vault/infrastructure/sops-age-key
    target: ~/.config/sops/age/keys.txt
    mode: "0600"
```

Then run once:

```bash
# One-time setup (fetches key from 1Password)
comet bootstrap

# Check status
comet bootstrap status

# Now all Comet commands can decrypt SOPS files automatically
comet plan dev
```

:::tip Why Bootstrap?

SOPS needs the `SOPS_AGE_KEY` or `SOPS_AGE_KEY_FILE` environment variable to decrypt files. Instead of fetching this from 1Password on every command (4s overhead), bootstrap caches it to `~/.config/sops/age/keys.txt` once, making all subsequent commands fast.

:::
    age: age1xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
    # Or for GPG:
    # pgp: YOUR_GPG_FINGERPRINT
```

## Creating Encrypted Secrets

### 1. Create a Secrets File

Create an unencrypted template first:

```yaml title="secrets.yaml"
database:
  password: changeme
  admin_user: admin
  
api_keys:
  google: your-api-key
  stripe: your-stripe-key
  
ssh_keys:
  deploy_key: |
    -----BEGIN OPENSSH PRIVATE KEY-----
    ...
    -----END OPENSSH PRIVATE KEY-----
```

### 2. Encrypt the File

```bash
# Encrypt the file (creates secrets.enc.yaml)
sops -e secrets.yaml > secrets.enc.yaml

# Or edit directly with SOPS
sops secrets.enc.yaml
```

The encrypted file will look like:

```yaml title="secrets.enc.yaml"
database:
  password: ENC[AES256_GCM,data:jfkdjfkd==,iv:xxx...]
  admin_user: ENC[AES256_GCM,data:kdfjkd==,iv:xxx...]
sops:
  kms: []
  gcp_kms: []
  age:
    - recipient: age1xxxxx...
      enc: |
        -----BEGIN AGE ENCRYPTED FILE-----
        ...
        -----END AGE ENCRYPTED FILE-----
  ...
```

## Using Secrets in Stacks

Reference encrypted secrets using the `sops://` URI scheme:

```javascript title="stacks/production.stack.js"
const database = component('cloudsql', 'modules/database', {
  // Reference individual secrets using JSON path syntax
  password: '{{ secrets "sops://secrets.enc.yaml#/database/password" }}',
  admin_user: '{{ secrets "sops://secrets.enc.yaml#/database/admin_user" }}'
})

const api = component('api-server', 'modules/application', {
  google_api_key: '{{ secrets "sops://secrets.enc.yaml#/api_keys/google" }}',
  stripe_api_key: '{{ secrets "sops://secrets.enc.yaml#/api_keys/stripe" }}'
})
```

### SOPS URI Format

```
sops://<file-path>#<json-path>
```

- `file-path` - Path to the encrypted SOPS file (relative to project root)
- `json-path` - JSON path to the specific secret (using `/` as separator)

## Examples

### Database Credentials

```javascript
const db = component('database', 'modules/cloudsql', {
  instance_name: 'prod-db',
  database_version: 'POSTGRES_14',
  
  // Encrypted secrets
  root_password: '{{ secrets "sops://secrets.enc.yaml#/database/root_password" }}',
  app_db_password: '{{ secrets "sops://secrets.enc.yaml#/database/app_password" }}'
})
```

### API Keys

```javascript
const app = component('application', 'modules/k8s-app', {
  env_vars: {
    DATABASE_URL: '{{ (state "data" "db").connection_string }}',
    
    // Encrypted API keys
    STRIPE_SECRET_KEY: '{{ secrets "sops://secrets.enc.yaml#/api_keys/stripe" }}',
    SENDGRID_API_KEY: '{{ secrets "sops://secrets.enc.yaml#/api_keys/sendgrid" }}',
    JWT_SECRET: '{{ secrets "sops://secrets.enc.yaml#/app/jwt_secret" }}'
  }
})
```

### SSH Keys

```javascript
const deployment = component('deploy-key', 'modules/secret', {
  name: 'github-deploy-key',
  data: {
    ssh_key: '{{ secrets "sops://secrets.enc.yaml#/ssh_keys/deploy_key" }}'
  }
})
```

## Environment-Specific Secrets

Maintain separate secrets files for each environment:

```
secrets/
├── dev.enc.yaml
├── staging.enc.yaml
└── production.enc.yaml
```

Reference the appropriate file in each stack:

```javascript title="stacks/production.stack.js"
const db = component('db', 'modules/database', {
  password: '{{ secrets "sops://secrets/production.enc.yaml#/database/password" }}'
})
```

```javascript title="stacks/dev.stack.js"
const db = component('db', 'modules/database', {
  password: '{{ secrets "sops://secrets/dev.enc.yaml#/database/password" }}'
})
```

## Best Practices

### 1. Never Commit Unencrypted Secrets

Add to `.gitignore`:

```gitignore title=".gitignore"
# Unencrypted secrets
secrets.yaml
secrets-*.yaml
!secrets*.enc.yaml  # Allow encrypted files
```

### 2. Use Descriptive Paths

```yaml
# ✅ Good: Clear hierarchy
database:
  production:
    password: secret123
    user: produser
  staging:
    password: secret456
    user: staginguser

# ❌ Bad: Flat and unclear
db_prod_pass: secret123
db_prod_user: produser
```

### 3. Rotate Secrets Regularly

```bash
# Edit encrypted file
sops secrets.enc.yaml

# Update the secret values, then apply
comet apply production database
```

### 4. Separate Keys by Environment

Different encryption keys for different environments:

```yaml title=".sops.yaml"
creation_rules:
  - path_regex: secrets/dev\.enc\.yaml$
    age: age1dev...
  
  - path_regex: secrets/production\.enc\.yaml$
    age: age1prod...
```

### 5. Use Key Management Services for Production

For production environments, use cloud KMS:

```yaml title=".sops.yaml"
creation_rules:
  - path_regex: secrets/production\.enc\.yaml$
    gcp_kms: projects/my-project/locations/global/keyRings/sops/cryptoKeys/sops-key
```

## Decrypting Secrets

To view or edit encrypted secrets:

```bash
# Decrypt to stdout
sops -d secrets.enc.yaml

# Edit in place (decrypts, opens editor, encrypts on save)
sops secrets.enc.yaml
```

## Integration with CI/CD

### Option 1: Direct Environment Variables

Provide the decryption key directly to your CI/CD pipeline:

```bash
# For age
export SOPS_AGE_KEY_FILE=/path/to/key.txt

# For GPG
gpg --import private-key.asc

# Then run Comet commands
comet apply production
```

### Option 2: Using Bootstrap

If your CI/CD has 1Password CLI available, use bootstrap for consistency:

```yaml
# comet.yaml
bootstrap:
  - name: sops-key
    type: secret
    source: op://vault/sops-key/private
    target: ~/.config/sops/age/keys.txt
    mode: "0600"
```

```bash
# In CI/CD pipeline
comet bootstrap
comet apply production
```

## Troubleshooting

### "Failed to get data key" Error

Make sure your encryption key is available:

```bash
# For age
export SOPS_AGE_KEY_FILE=/path/to/key.txt

# For GPG
gpg --list-secret-keys  # Verify key is imported
```

### Secret Not Found

Verify the JSON path:

```bash
# View decrypted structure
sops -d secrets.enc.yaml

# Check the path matches your reference
# If the structure is:
# database:
#   password: xxx
# Then use: #/database/password
```
