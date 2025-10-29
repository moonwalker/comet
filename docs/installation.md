# Installation Guide

## Installing Comet

### Quick Install (Recommended)

```bash
curl -fsSL https://moonwalker.github.io/comet/install.sh | sh
```

This installs Comet to `~/.local/bin/comet`.

**Add to PATH** (if needed):
```bash
# Add to your shell profile (~/.bashrc, ~/.zshrc, etc.)
export PATH="$HOME/.local/bin:$PATH"
```

### Manual Installation

Download the latest release for your platform from the [releases page](https://github.com/moonwalker/comet/releases):

```bash
# Example for macOS ARM64
curl -L https://github.com/moonwalker/comet/releases/download/v0.6.7/comet_0.6.7_darwin_arm64.tar.gz | tar xz
mv comet ~/.local/bin/
chmod +x ~/.local/bin/comet
```

### Building from Source

Requires Go 1.23 or later.

```bash
git clone https://github.com/moonwalker/comet.git
cd comet
go build
./comet version
```

## Prerequisites

Before using Comet, you need **OpenTofu** (recommended) or **Terraform** installed.

### OpenTofu (Recommended)

**macOS:**
```bash
brew install opentofu
```

**Linux/macOS - Official Installer:**
```bash
curl -fsSL https://get.opentofu.org/install-opentofu.sh -o install-opentofu.sh
chmod +x install-opentofu.sh
./install-opentofu.sh --install-method standalone
```

**Signature Verification:**

The OpenTofu installer verifies download signatures for security. Install one of:

```bash
# Option 1: cosign (recommended)
brew install cosign

# Option 2: GPG
brew install gnupg
```

Or skip verification (not recommended):
```bash
./install-opentofu.sh --install-method standalone --skip-verify
```

**Other platforms:** See [OpenTofu installation docs](https://opentofu.org/docs/intro/install/)

### Terraform (Alternative)

**macOS:**
```bash
brew install terraform
```

**Linux/Windows:** See [Terraform installation docs](https://developer.hashicorp.com/terraform/install)

### Configuring Comet

Tell Comet which tool to use:

```yaml
# comet.yaml
tf_command: tofu       # or: terraform
```

Default is `tofu`.

## Verification

```bash
# Check Comet
comet version

# Check OpenTofu/Terraform
tofu version
# or
terraform version
```

## Next Steps

- Read the [Quick Start](../README.md#usage)
- See [Best Practices](best-practices.md)
- Browse [Examples](../stacks/_examples/)
