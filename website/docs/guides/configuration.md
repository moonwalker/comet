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

Pre-load environment variables before any Comet command runs. This is useful for:
- Setting secrets needed during stack parsing (like `SOPS_AGE_KEY`)
- Configuring Terraform behavior (like `TF_LOG`)
- Setting cloud provider credentials

```yaml
# comet.yaml
env:
  # Plain values - fast and simple
  TF_LOG: DEBUG
  AWS_REGION: us-west-2
  PROJECT_ID: my-gcp-project
```

### Features

- **Shell precedence**: Environment variables already set in your shell take precedence over config values
- **Secret resolution**: Supports `op://` (1Password) and `sops://` references
- **Early loading**: Variables are set before stack parsing begins

### Secret References

You can reference secrets from 1Password or SOPS directly in your config:

```yaml
env:
  # 1Password reference
  SOPS_AGE_KEY: op://vault/sops-age-key/private
  
  # SOPS reference
  API_TOKEN: sops://secrets.enc.yaml#/api/token
```

:::warning Performance Impact

Secret references (`op://`, `sops://`) are resolved on **EVERY** Comet command, including fast operations like `comet list`. This can add **3-5 seconds** per secret due to CLI overhead from tools like the 1Password CLI.

**Recommended approach for frequently-used secrets:**

```bash
# Set in your shell once (one-time cost)
export SOPS_AGE_KEY=$(op read "op://vault/sops-age-key/private")

# Or add to your shell config (~/.bashrc, ~/.zshrc, ~/.config/fish/config.fish)
```

**Use secret references in `comet.yaml` only when:**
- Secrets change frequently
- Shell setup is inconvenient (CI/CD environments)
- The 3-5 second overhead is acceptable for your workflow

:::

### Example: CI/CD Environment

In CI/CD environments where secrets are provided by the platform, config-based loading might be more convenient despite the overhead:

```yaml
# comet.yaml - CI/CD focused
env:
  # GitHub Actions provides these
  SOPS_AGE_KEY: op://ci-cd/sops-age-key/private
  
  # Plain values from environment
  TF_LOG: ${{ env.TF_LOG }}
```

### Example: Local Development

For local development, prefer shell environment variables for speed:

```yaml
# comet.yaml - Local development
env:
  # Only plain values, no secret resolution
  TF_LOG: DEBUG
  AWS_REGION: us-west-2
```

```bash
# ~/.bashrc, ~/.zshrc, or ~/.config/fish/config.fish
export SOPS_AGE_KEY=$(op read "op://vault/sops-age-key/private")
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

### Keep Secrets Out of Config

For frequently-accessed secrets, set them in your shell instead of `comet.yaml`:

```bash
# ✅ Good - one-time cost
export SOPS_AGE_KEY=$(op read "op://vault/key/private")

# ❌ Avoid - 4s penalty on every command
# comet.yaml with: SOPS_AGE_KEY: op://vault/key/private
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
