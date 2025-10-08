# Comet Best Practices

## Project Structure

### Recommended Layout

```
infrastructure/
├── comet.yaml                 # Comet configuration
├── secrets.yaml               # Unencrypted secrets template
├── secrets.enc.yaml           # SOPS-encrypted secrets
├── .gitignore                 # Ignore generated files
├── stacks/                    # Stack definitions
│   ├── shared.js              # Shared settings
│   ├── dev.js
│   ├── staging.js
│   └── production.js
└── modules/                   # Terraform modules
    ├── vpc/
    │   ├── main.tf
    │   ├── variables.tf
    │   └── outputs.tf
    ├── kubernetes/
    └── database/
```

### Stack Organization

**DO:** Use a shared configuration file
```javascript
// stacks/shared.js
const settings = {
  project_name: 'myapp',
  domain: 'example.com',
  regions: {
    primary: 'us-central1',
    secondary: 'us-east1'
  }
}

module.exports = { settings }
```

**DO:** Import shared settings in each stack
```javascript
// stacks/production.js
const { settings } = require('./shared.js')

stack('production', { settings })
```

**DON'T:** Duplicate configuration across stacks
```javascript
// ❌ Bad: Duplicated in each stack
stack('production', {
  project_name: 'myapp',  // Repeated
  domain: 'example.com'   // Repeated
})
```

## Component Design

### Keep Components Focused

**DO:** Small, single-purpose components
```javascript
// Good: Each component has one responsibility
const vpc = component('vpc', 'modules/vpc', {
  cidr_block: '10.0.0.0/16'
})

const gke = component('gke', 'modules/gke', {
  network: vpc.id,
  cluster_name: 'my-cluster'
})
```

**DON'T:** Mega-components that do everything
```javascript
// ❌ Bad: One component doing too much
const infrastructure = component('all', 'modules/everything', {
  // Too many responsibilities
})
```

### Use Consistent Naming

**DO:** Clear, descriptive component names
```javascript
component('vpc-prod', 'modules/vpc', ...)
component('gke-main-cluster', 'modules/gke', ...)
component('cloudsql-primary', 'modules/database', ...)
```

**DON'T:** Cryptic abbreviations
```javascript
// ❌ Bad: Unclear names
component('v1', 'modules/vpc', ...)
component('k', 'modules/gke', ...)
```

## Cross-Stack References

### Minimize Dependencies

**DO:** Only reference what you need
```javascript
// Good: Specific references
const app = component('app', 'modules/app', {
  vpc_id: '{{ (state "infra" "vpc").id }}'
})
```

**DON'T:** Create unnecessary coupling
```javascript
// ❌ Bad: Circular or complex dependencies
// stack1 → stack2 → stack3 → stack1 (circular!)
```

### Document Dependencies

**DO:** Comment cross-stack references
```javascript
// Requires: 'infra' stack must be applied first
const app = component('app', 'modules/app', {
  // Reference to infrastructure stack VPC
  vpc_id: '{{ (state "infra" "vpc").id }}'
})
```

## Secrets Management

### Always Use SOPS

**DO:** Encrypt all secrets
```bash
# Create encrypted secrets file
sops secrets.enc.yaml
```

```javascript
// Reference encrypted secrets
const db = component('database', 'modules/db', {
  password: secrets('sops://secrets.enc.yaml#/database/password')
})
```

**DON'T:** Hardcode sensitive values
```javascript
// ❌ Bad: Plaintext secrets
const db = component('database', 'modules/db', {
  password: 'super-secret-password'  // Never do this!
})
```

### SOPS AGE Key Setup

SOPS requires the `SOPS_AGE_KEY` to be set **before** stack parsing begins. Use `comet.yaml` to pre-load it:

**DO:** Configure in comet.yaml (recommended)
```yaml
# comet.yaml
env:
  # Load SOPS AGE key from 1Password before stack parsing
  SOPS_AGE_KEY: op://ci-cd/sops-age-key/private
```

**Alternative:** Export in shell
```bash
# In your shell profile or CI/CD
export SOPS_AGE_KEY="AGE-SECRET-KEY-1..."
```

**Why use comet.yaml?**
- ✅ Automatic - no manual shell setup needed
- ✅ Team-consistent - everyone uses the same config
- ✅ CI/CD friendly - works seamlessly in pipelines
- ✅ Secure - secrets fetched from 1Password on-demand

### Use Path-Based Organization

```yaml
# secrets.enc.yaml
database:
  dev:
    password: "dev-password"
  prod:
    password: "prod-password"
github:
  token: "ghp_..."
cloudflare:
  api_token: "..."
```

```javascript
// Reference with clear paths
password: secrets('sops://secrets.enc.yaml#/database/{{ .stack }}/password')
```

## Backend Configuration

### Use Consistent Naming

**DO:** Include stack and component in state path
```javascript
backend('gcs', {
  bucket: 'my-terraform-state',
  prefix: '${project}/{{ .stack }}/{{ .component }}/terraform.tfstate'
})
```

**DON'T:** Use flat state structure
```javascript
// ❌ Bad: No organization
backend('gcs', {
  bucket: 'my-terraform-state',
  prefix: 'state'  // All state in one place
})
```

## Version Control

### .gitignore

Always ignore generated files:

```gitignore
# Terraform
.terraform/
*.tfstate
*.tfstate.*
*.tfvars
.terraform.lock.hcl

# Comet generated
**/backend.tf.json
**/providers_gen.tf
**/*-*.tfvars.json
**/*.planfile

# Secrets (keep only encrypted)
secrets.yaml
!secrets.enc.yaml
```

### What to Commit

**DO commit:**
- ✅ Stack definitions (`.js` files)
- ✅ Module source code
- ✅ Encrypted secrets (`.enc.yaml`)
- ✅ `comet.yaml` configuration
- ✅ Documentation

**DON'T commit:**
- ❌ Generated Terraform files
- ❌ `.terraform/` directories
- ❌ State files
- ❌ Unencrypted secrets
- ❌ Plan files

## Testing

### Test Stack Syntax

```bash
# Check if stacks parse correctly
comet list
comet list dev
```

### Dry Run Before Apply

```bash
# Always plan first
comet plan production vpc

# Review the plan carefully
comet apply production vpc
```

### Test in Lower Environments First

```bash
# Development → Staging → Production
comet apply dev
comet apply staging
comet apply production
```

## Team Workflows

### Code Review

**Include in PRs:**
- Stack definition changes
- Module changes
- Documentation updates
- Test results from `comet plan`

### CI/CD Integration

```yaml
# Example GitHub Actions workflow
name: Terraform Plan

on:
  pull_request:
    paths:
      - 'stacks/**'
      - 'modules/**'

jobs:
  plan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Comet
        run: |
          # Install comet
          
      - name: Plan Dev
        run: comet plan dev
        
      - name: Plan Staging
        run: comet plan staging
```

### Environment Promotion

```bash
# 1. Test in dev
comet apply dev

# 2. Verify in staging
comet apply staging

# 3. Plan production
comet plan production

# 4. Apply production (after approval)
comet apply production
```

## Performance Tips

### Parallel Execution

When components are independent:

```bash
# Apply multiple components in parallel
comet apply dev vpc &
comet apply dev database &
wait
```

### Component Sizing

- Keep components small and focused
- Avoid massive components with hundreds of resources
- Split large infrastructures across multiple components

### State Backend Optimization

- Use remote backends (GCS, S3) for team collaboration
- Enable state locking to prevent conflicts
- Use separate state files per component (Comet does this automatically)

## Troubleshooting

### Common Issues

**Issue:** "Backend configuration has changed"
```bash
# Solution: Re-initialize Terraform
rm -rf .terraform
comet init <stack> <component>
```

**Issue:** "Required plugins are not installed"
```bash
# Solution: Initialize providers and backends
comet init <stack> <component>
```

**Issue:** Cross-stack reference returns null
```bash
# Solution: Ensure referenced stack is applied first
comet apply infra vpc
comet init app webapp   # Initialize to query outputs
comet apply app webapp  # Now vpc outputs are available
```

**Issue:** SOPS decryption fails
```bash
# Solution 1: Set SOPS_AGE_KEY in comet.yaml (recommended)
# comet.yaml:
# env:
#   SOPS_AGE_KEY: op://ci-cd/sops-age-key/private

# Solution 2: Export in shell
export SOPS_AGE_KEY="your-age-key"
comet apply dev
```

## Upgrade Strategy

When upgrading Comet:

1. **Test in dev environment first**
2. **Review changelog for breaking changes**
3. **Update one environment at a time**
4. **Keep Comet version consistent across team**

```bash
# Check version
comet version

# Upgrade gradually
# dev → staging → production
```

## Documentation

### Document Your Stacks

Add comments to stack files:

```javascript
/**
 * Production Environment Stack
 * 
 * This stack provisions the production infrastructure including:
 * - VPC with public and private subnets
 * - GKE cluster with autoscaling
 * - Cloud SQL database with failover
 * 
 * Prerequisites:
 * - GCP project must exist
 * - SOPS_AGE_KEY must be set
 * 
 * Usage:
 *   comet plan production
 *   comet apply production
 */

const { settings } = require('./shared.js')
stack('production', { settings })

// ... component definitions
```

### Maintain a README

Keep a project-level README with:
- Quick start guide
- Environment descriptions
- Deployment procedures
- Troubleshooting tips
