---
sidebar_position: 1
---

# Best Practices

This guide covers recommended practices for organizing and managing infrastructure with Comet.

## Project Structure

### Recommended Layout

```
infrastructure/
├── comet.yaml                 # Comet configuration
├── .sops.yaml                 # SOPS configuration
├── secrets.enc.yaml           # Encrypted secrets
├── .gitignore                 # Version control exclusions
├── README.md                  # Project documentation
├── stacks/                    # Stack definitions
│   ├── shared.js              # Shared settings
│   ├── dev.stack.js
│   ├── staging.stack.js
│   └── production.stack.js
└── modules/                   # Terraform modules
    ├── vpc/
    │   ├── main.tf
    │   ├── variables.tf
    │   └── outputs.tf
    ├── gke/
    ├── database/
    └── storage/
```

## Stack Organization

### Use Shared Configuration

**DO:** Create a shared settings file
```javascript title="stacks/shared.js"
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
```javascript title="stacks/production.js"
const { settings } = require('./shared.js')

stack('production', { 
  ...settings,
  environment: 'production'
})
```

**DON'T:** Duplicate configuration across stacks
```javascript
// ❌ Bad: Duplicated in each stack
stack('production', {
  project_name: 'myapp',  // Repeated
  domain: 'example.com'   // Repeated
})
```

### Organize Stacks by Lifecycle

Separate stacks based on deployment frequency and dependencies:

```
stacks/
├── foundation.stack.js      # Rarely changes (GCP project, IAM)
├── networking.stack.js      # Infrastructure (VPC, subnets)
├── kubernetes.stack.js      # Cluster configuration
├── data.stack.js           # Databases and storage
└── applications.stack.js    # Frequently updated apps
```

## Component Design

### Keep Components Focused

**DO:** Small, single-purpose components
```javascript
// ✅ Good: Each component has one responsibility
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
component('vpc-main', 'modules/vpc', ...)
component('gke-primary-cluster', 'modules/gke', ...)
component('cloudsql-users-db', 'modules/database', ...)
```

**DON'T:** Cryptic abbreviations
```javascript
// ❌ Bad: Unclear names
component('v1', 'modules/vpc', ...)
component('k', 'modules/gke', ...)
```

## Secrets Management

### Always Encrypt Secrets

**DO:** Use SOPS for all sensitive data
```bash
# Create encrypted secrets file
sops secrets.enc.yaml
```

```javascript
const db = component('database', 'modules/db', {
  password: '{{ secrets "sops://secrets.enc.yaml#/database/password" }}'
})
```

**DON'T:** Hardcode sensitive values
```javascript
// ❌ Bad: Plaintext secrets
const db = component('database', 'modules/db', {
  password: 'super-secret-password'  // Never do this!
})
```

### Organize Secrets by Path

```yaml title="secrets.enc.yaml"
database:
  dev:
    password: "dev-password"
  prod:
    password: "prod-password"
    
api_keys:
  stripe: "sk_live_..."
  sendgrid: "SG...."
```

Use template variables for environment-specific secrets:

```javascript
password: '{{ secrets "sops://secrets.enc.yaml#/database/{{ .stack }}/password" }}'
```

## Version Control

### .gitignore Configuration

```gitignore title=".gitignore"
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

# Secrets (keep only encrypted)
secrets.yaml
!secrets.enc.yaml

# SOPS keys (never commit)
key.txt
*.age
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

## Testing Strategy

### Test in Lower Environments First

```bash
# Development → Staging → Production
comet plan dev
comet apply dev

comet plan staging
comet apply staging

comet plan production
comet apply production
```

### Always Plan Before Apply

```bash
# ✅ Good: Review changes first
comet plan production vpc
# Review output carefully
comet apply production vpc

# ❌ Bad: Direct apply without review
comet apply production vpc
```

## Cross-Stack References

### Minimize Dependencies

**DO:** Only reference what you need
```javascript
const app = component('app', 'modules/app', {
  vpc_id: '{{ (state "infra" "vpc").id }}'
})
```

**DON'T:** Create circular dependencies
```javascript
// ❌ Bad: Circular dependencies
// stack1 → stack2 → stack3 → stack1 (circular!)
```

### Document Dependencies

```javascript
/**
 * Requires: 'infrastructure' stack must be applied first
 */
const app = component('app', 'modules/app', {
  // Reference to infrastructure stack VPC
  vpc_id: '{{ (state "infrastructure" "vpc").id }}'
})
```

## Backend Configuration

### Use Consistent Naming

```javascript
backend('gcs', {
  bucket: 'my-terraform-state',
  prefix: '{{ .settings.project_name }}/{{ .stack }}/{{ .component }}'
})
```

This creates organized state paths:
```
myapp/dev/vpc/
myapp/dev/gke/
myapp/production/vpc/
myapp/production/gke/
```

## Performance

### Component Sizing

- Keep components small and focused (< 100 resources)
- Split large infrastructures across multiple components
- Use separate components for frequently changing resources

### Parallel Execution

Independent components can be applied in parallel:

```bash
# Components with no dependencies can run simultaneously
comet apply dev storage &
comet apply dev monitoring &
wait
```

## Documentation

### Document Your Stacks

```javascript
/**
 * Production Environment Stack
 * 
 * Provisions production infrastructure including:
 * - VPC with public and private subnets
 * - GKE cluster with autoscaling
 * - Cloud SQL with high availability
 * 
 * Prerequisites:
 * - GCP project must exist
 * - SOPS_AGE_KEY environment variable must be set
 * 
 * Usage:
 *   comet plan production
 *   comet apply production
 */

const { settings } = require('./shared.js')
stack('production', { ...settings })
```

### Maintain README Documentation

Include in your project README:
- Quick start guide
- Environment descriptions
- Deployment procedures
- Troubleshooting tips
- Team contacts

## CI/CD Integration

### Example GitHub Actions Workflow

```yaml title=".github/workflows/terraform-plan.yml"
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
          
      - name: Decrypt secrets
        env:
          SOPS_AGE_KEY: ${{ secrets.SOPS_AGE_KEY }}
        run: |
          echo "$SOPS_AGE_KEY" > key.txt
          
      - name: Plan Dev
        run: comet plan dev
```

## Troubleshooting

### Common Issues

**Backend configuration changed:**
```bash
rm -rf .terraform
comet apply <stack> <component>
```

**Cross-stack reference returns null:**
```bash
# Ensure referenced stack is applied first
comet apply infrastructure vpc
comet apply application webapp
```

**SOPS decryption fails:**
```bash
export SOPS_AGE_KEY_FILE=/path/to/key.txt
comet apply dev
```
