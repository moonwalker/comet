---
sidebar_position: 6
---

# Configuration

Comet can be configured using a `comet.yaml` file in your project root. This file controls various aspects of Comet's behavior.

## Basic Configuration

```yaml
# comet.yaml
stacks_dir: stacks              # Directory containing stack files
work_dir: stacks/_components    # Working directory for components
generate_backend: false         # Auto-generate backend.tf.json
log_level: INFO                 # Log verbosity (DEBUG, INFO, WARN, ERROR)
tf_command: tofu                # Use 'tofu' or 'terraform'
```

### Configuration Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `stacks_dir` | string | `stacks` | Directory containing your stack files |
| `work_dir` | string | `stacks/_components` | Working directory where Terraform files are generated |
| `generate_backend` | boolean | `false` | Auto-generate `backend.tf.json` files |
| `log_level` | string | `INFO` | Logging verbosity: DEBUG, INFO, WARN, ERROR |
| `tf_command` | string | `tofu` | Terraform executor: `tofu` or `terraform` |
| `env` | map | `{}` | Environment variables to set before commands run |

## Environment Variables

Set plain environment variables that are loaded before any Comet command runs:

```yaml
# comet.yaml
env:
  # Plain values only
  TF_LOG: DEBUG
  AWS_REGION: us-west-2
  PROJECT_ID: my-gcp-project
```

### Features

- **Plain values only**: Fast startup, no secret resolution overhead
- **Shell precedence**: Environment variables already set in your shell take precedence
- **Early loading**: Variables are set before stack parsing begins

:::info Secret Management

The `env` section only supports plain values for fast performance. For secrets, use the `bootstrap` feature below.

:::

## Bootstrap: One-Time Secret Setup

Bootstrap fetches secrets from 1Password or SOPS and caches them locally. Run it once, then all your commands are fast!

```yaml
# comet.yaml
bootstrap:
  - name: sops-age-key
    type: secret
    source: op://vault/infrastructure/sops-age-key
    target: ~/.config/sops/age/keys.txt
    mode: "0600"
```

**Usage:**

```bash
# One-time setup (takes ~4 seconds)
comet bootstrap

# Check what's been set up
comet bootstrap status

# Now all commands are fast!
comet plan dev    # 100ms instead of 4s
comet apply dev
```

**Bootstrap Commands:**

```bash
# Run all bootstrap steps
comet bootstrap

# Force re-run all steps (even if already completed)
comet bootstrap --force

# Check status of all steps
comet bootstrap status
# Output:
#   ✅ check-tools      Completed (just now)
#   ✅ sops-age-key     Completed (5m ago)
#   ⚪ gcloud-auth      Not run yet

# Clear bootstrap state (doesn't delete cached files)
comet bootstrap clear
```

### Bootstrap Step Types

#### Secret Steps

Fetch secrets from 1Password or SOPS and save to local files with proper permissions:

```yaml
bootstrap:
  # 1Password secret
  - name: sops-key
    type: secret
    source: op://vault/item/field        # 1Password reference
    target: ~/.config/sops/age/keys.txt  # Save location (~ expanded)
    mode: "0600"                         # File permissions (optional, default: 0600)
  
  # SOPS secret
  - name: api-token
    type: secret
    source: sops://secrets.enc.yaml#/api/token  # SOPS file reference
    target: ~/.secrets/api-token
    mode: "0400"                         # Read-only for extra security
    optional: true                       # Don't fail if source doesn't exist
```

**Features:**
- Automatically creates parent directories
- Supports both `op://` (1Password) and `sops://` (SOPS) sources
- Customizable file permissions (default: `0600`)
- Path expansion (`~/` → home directory)
- Optional steps won't fail the bootstrap

#### Check Steps

Verify required tools are installed before proceeding:

```yaml
bootstrap:
  - name: check-tools
    type: check
    command: op,sops,tofu  # Comma-separated binary names
```

**Features:**
- Checks if binaries exist in PATH using `exec.LookPath`
- Fast execution (< 1ms)
- Fails immediately if any binary is missing
- Best practice: Put check steps first for early failure
- Example use cases: verify CLI tools, check system requirements

#### Command Steps

Run arbitrary shell commands for authentication, initialization, or custom setup:

```yaml
bootstrap:
  - name: gcloud-auth
    type: command
    command: gcloud auth application-default login
    check: gcloud auth application-default print-access-token  # Optional: check if already done
    optional: true  # Don't fail if command errors
  
  - name: install-deps
    type: command
    command: npm install --global some-tool
```

**Features:**
- Runs commands via `sh -c` (full shell access)
- Streams stdout/stderr to console
- Optional `check` command to skip if already done
- Mark as `optional: true` for non-critical setup
- Example use cases: cloud authentication, tool installation, git config

### Bootstrap Features

**Performance & Execution:**
- **One-time cost**: Slow operations only happen during `comet bootstrap`
- **30x faster**: Commands run in ~100ms instead of 4+ seconds
- **Idempotent**: Safe to run multiple times, skips already-completed steps
- **Sequential execution**: Steps run in order for predictable setup
- **Early termination**: Fails fast on errors unless marked `optional`

**State Management:**
- **State tracking**: Tracks completion in `.comet/bootstrap.state`
- **Smart skipping**: Only re-runs steps if target files are missing
- **Force refresh**: Use `--force` to re-run all steps
- **Custom checks**: Optional `check` field to determine if step is needed

**Step Types:**
- **Secret steps**: Fetch from 1Password (`op://`) or SOPS (`sops://`) and cache to files
- **Command steps**: Run arbitrary shell commands for authentication, setup, etc.
- **Check steps**: Verify required binaries/tools are installed

**Configuration:**
- **Optional steps**: Mark steps as `optional: true` to continue on failure
- **File permissions**: Customize with `mode` field (default: `0600`)
- **Path expansion**: Supports `~/` and environment variables in paths
- **Multiple sources**: Mix 1Password, SOPS, and commands in one workflow

### Example: Complete Bootstrap Setup

```yaml
# comet.yaml
bootstrap:
  # 1. Check required tools (fast, fails early)
  - name: check-tools
    type: check
    command: op,sops,tofu
  
  # 2. Fetch SOPS key from 1Password
  - name: sops-key
    type: secret
    source: op://vault/infrastructure/sops-age-key
    target: ~/.config/sops/age/keys.txt
    mode: "0600"
  
  # 3. Authenticate with cloud provider (optional)
  - name: gcloud-auth
    type: command
    command: gcloud auth application-default login
    check: gcloud auth application-default print-access-token
    optional: true
  
  # 4. Fetch additional secrets
  - name: api-credentials
    type: secret
    source: op://vault/api/credentials
    target: ~/.secrets/api.json
    mode: "0600"
```

### Advanced Bootstrap Examples

**Multi-cloud authentication:**
```yaml
bootstrap:
  # AWS authentication
  - name: aws-sso
    type: command
    command: aws sso login --profile production
    check: aws sts get-caller-identity --profile production
    optional: true
  
  # GCP authentication
  - name: gcp-auth
    type: command
    command: gcloud auth application-default login
    check: gcloud auth application-default print-access-token
    optional: true
  
  # Azure authentication
  - name: azure-login
    type: command
    command: az login
    optional: true
```

**Multiple environment secrets:**
```yaml
bootstrap:
  # Development SOPS key
  - name: sops-dev
    type: secret
    source: op://vault/sops-dev/key
    target: ~/.config/sops/dev-key.txt
  
  # Production SOPS key (different key!)
  - name: sops-prod
    type: secret
    source: op://vault/sops-prod/key
    target: ~/.config/sops/prod-key.txt
    mode: "0400"  # Read-only for production
```

**Developer machine setup:**
```yaml
bootstrap:
  # Check all required tools
  - name: check-dev-tools
    type: check
    command: git,op,sops,tofu,gcloud,kubectl,helm
  
  # Fetch all secrets
  - name: sops-key
    type: secret
    source: op://vault/sops/key
    target: ~/.config/sops/age/keys.txt
  
  - name: github-token
    type: secret
    source: op://vault/github/pat
    target: ~/.config/gh/token
    mode: "0600"
  
  # Configure git
  - name: git-config
    type: command
    command: |
      git config --global user.email "dev@company.com"
      git config --global user.name "Developer"
    optional: true
  
  # Kubernetes setup
  - name: k8s-config
    type: command
    command: gcloud container clusters get-credentials prod-cluster --region us-central1
    optional: true
```

### Migration from v0.5.0

If you were using `op://` or `sops://` in the `env` section (removed in v0.6.0), migrate to `bootstrap`:

```yaml
# OLD (v0.5.0) - Slow on every command
env:
  SOPS_AGE_KEY: op://vault/key/private

# NEW (v0.6.0) - Fast after bootstrap
bootstrap:
  - name: sops-key
    type: secret
    source: op://vault/key/private
    target: ~/.config/sops/age/keys.txt
    mode: "0600"
```

## Command-Line Flags

Configuration can also be overridden via command-line flags:

```bash
# Override config file location
comet --config=custom-config.yaml list

# Override stacks directory
comet --dir=infrastructure/stacks list

# Override log level
COMET_LOG_LEVEL=debug comet list
```

## Environment Variable Precedence

When the same variable is defined in multiple places, Comet uses this precedence order (highest to lowest):

1. **Shell environment** - Variables already set in your shell
2. **Config file** - Variables defined in `comet.yaml`
3. **Default values** - Comet's built-in defaults

**Example:**

```yaml
# comet.yaml
env:
  AWS_REGION: us-west-2
```

```bash
# Shell takes precedence
export AWS_REGION=eu-west-1

comet apply dev  # Uses eu-west-1, not us-west-2
```

## Best Practices

### Use Bootstrap for Secrets

For secrets needed by Comet or Terraform, use the bootstrap feature:

```yaml
# comet.yaml
bootstrap:
  - name: sops-key
    type: secret
    source: op://vault/key/private
    target: ~/.config/sops/age/keys.txt
    mode: "0600"
```

```bash
# One-time setup
comet bootstrap

# Now all commands are fast
comet plan dev  # 100ms, not 4s
```

### Use Version Control

Commit `comet.yaml` to version control for team consistency:

```yaml
# comet.yaml - safe to commit
stacks_dir: infrastructure/stacks
work_dir: infrastructure/_components
generate_backend: true
tf_command: tofu

env:
  # Never commit actual secrets!
  # These are just variable names/references
  AWS_REGION: us-west-2
  TF_LOG: INFO
```

### Project-Specific Settings

Use different configurations for different environments:

```bash
# Development
comet --config=comet.dev.yaml apply dev

# Production (stricter settings)
comet --config=comet.prod.yaml apply production
```

## Debug Logging

Enable debug logging to troubleshoot performance or parsing issues:

```yaml
# comet.yaml
log_level: debug
```

Or via environment variable:

```bash
COMET_LOG_LEVEL=debug comet list
```

Debug logs show:
- Stack parsing duration
- esbuild bundling time
- Secret resolution time
- Component registration details

## Related Documentation

- [Secrets Management](/docs/guides/secrets-management) - Detailed guide on working with secrets
- [CLI Reference](/docs/guides/cli-reference) - Complete command-line reference
- [Best Practices](/docs/advanced/best-practices) - Recommended patterns and workflows
