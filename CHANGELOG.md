# Changelog

All notable changes to Comet will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.6.6] - 2025-10-23

### Added
- **Auto-detect SOPS age key path** - Bootstrap `target` is now optional for SOPS age keys. If the source name contains "sops" and "age", it automatically uses the platform-specific default path
- **Helpful SOPS error messages** - When SOPS fails to decrypt due to missing age keys, provide clear hint suggesting `comet bootstrap` or setting `SOPS_AGE_KEY`

### Changed
- Improved bootstrap configuration - SOPS age key target path is now optional and auto-detected based on platform

## [0.6.5] - 2025-10-23

### Fixed
- **Bootstrap SOPS age key path resolution on macOS** - Bootstrap now correctly saves age keys to the platform-specific path that SOPS expects:
  - macOS without `XDG_CONFIG_HOME`: `~/Library/Application Support/sops/age/keys.txt`
  - macOS with `XDG_CONFIG_HOME`: `$XDG_CONFIG_HOME/sops/age/keys.txt`
  - Linux: `~/.config/sops/age/keys.txt` (or `$XDG_CONFIG_HOME/sops/age/keys.txt`)
  - Previously, bootstrap always saved to `~/.config/sops/age/keys.txt` on all platforms, causing SOPS to fail finding keys on macOS

## [0.6.4] - 2025-01-07

### Fixed
- Custom metadata field ordering now stable and consistent across runs
- Removed extra leading spaces from custom field values in table display

## [0.6.3] - 2025-01-07

### Changed
- **Enhanced metadata display in `comet list --details`**
  - Dynamic columns: only show owner/custom columns when data exists
  - Custom fields display in definition order (not alphabetically sorted)
  - Custom fields shown one per line for better readability
  - Optimized table width for smaller screens (20-char columns with wrapping)
  - Shortened paths by removing 'stacks/' prefix
  - Row lines between stacks for improved clarity

### Fixed
- Updated example stack files to work without requiring secret files
- Removed outdated examples that referenced deprecated features
- All examples now run successfully with `comet list --details`

## [0.6.2] - 2025-10-22

### Added
- **`metadata()` function** - Add metadata to stacks for better organization
  - Set description, owner, tags, and custom fields
  - View in `comet list` with smart truncation
  - `--details` flag shows full metadata including owner
  - Example: `metadata({ description: 'Production env', owner: 'platform-team', tags: ['prod'] })`

### Changed
- **`comet list` output** - Now displays stack metadata by default
  - Shows description (truncated at 50 chars) and first 3 tags
  - Use `--details` flag for full metadata including owner
  - More informative stack listings

## [0.6.1] - 2025-10-22

### Fixed
- Bootstrap secret files now properly end with newline character (POSIX standard)
- Ensures compatibility with tools that expect newline-terminated text files
- Prevents Git warnings about missing newlines at end of file

## [0.6.0] - 2025-10-16

### Added
- **`comet bootstrap` command** - One-time setup for secrets and dependencies. Fetches secrets from 1Password/SOPS and caches them locally, making all subsequent commands fast. No more 3-5 second delays on every command!
  - `comet bootstrap` - Run bootstrap steps
  - `comet bootstrap status` - Show what's been set up
  - `comet bootstrap clear` - Reset state
  - Bootstrap configuration in `comet.yaml` with support for secret fetching, command execution, and dependency checks
  - State tracking in `.comet/bootstrap.state`
  - Idempotent by default with `--force` flag to re-run

### Changed
- **BREAKING: Removed `op://` and `sops://` support from `env` section** - The `env` section now only supports plain values for fast startup. Use `comet bootstrap` instead for secret management.
- **`env` section is now fast** - No more slow secret resolution on every command. Plain environment variables only.

### Migration Guide
If you were using `op://` or `sops://` in your `env` section:

**Before (v0.5.0):**
```yaml
env:
  SOPS_AGE_KEY: op://vault/sops-key/private  # Slow on every command
```

**After (v0.6.0):**
```yaml
bootstrap:
  - name: sops-key
    type: secret
    source: op://vault/sops-key/private
    target: ~/.config/sops/age/keys.txt
    mode: "0600"

# Then run once: comet bootstrap
# All commands are now fast!
```

## [0.5.0] - 2025-10-10

### Added
- **Debug logging** - Added detailed debug logs for performance profiling of stack parsing, esbuild bundling, and secret resolution. Enable with `log_level: debug` in config or `COMET_LOG_LEVEL=debug` environment variable.
- **Configuration documentation** - New comprehensive configuration guide in website docs covering all options, environment variables, and performance considerations.
- **`comet types` command** - Generate TypeScript definitions for IDE support on-demand

### Fixed
- Skip parsing TypeScript definition files (`.d.ts`) to prevent parse errors

### Changed
- **Performance warning for config-based secrets** - Added warning when using `op://` or `sops://` references in `comet.yaml` env section, as these are resolved on every command and can add 3-5 seconds. Documentation now recommends setting frequently-used secrets in shell environment instead.
- TypeScript definitions are now opt-in via `comet types` instead of auto-generated

### Added
- **Config-based environment variables** - Pre-load environment variables from `comet.yaml` before any command runs. Perfect for setting `SOPS_AGE_KEY` and other secrets needed during stack parsing. Supports secret resolution via `op://` and `sops://` prefixes. Shell environment variables take precedence. ⚠️ **Note:** Secret resolution can be slow (3-5s per secret with 1Password CLI); consider setting in shell for frequently-used values.
- **`comet init` command** - Initialize backends and providers without running plan/apply operations. Useful for read-only operations like `comet output` or troubleshooting provider/backend initialization issues.
- **DSL Improvements** - Two core enhancements to reduce boilerplate by ~30%:
  - Bulk environment variables: `envs({})` accepts objects to set multiple vars at once
  - Secrets path shorthand: New `secret()` function with configurable defaults and dot notation support
- **"It's Just JavaScript!" philosophy** - Emphasized that users can create any helper functions they need
- **AGENTS.md** - Guidelines for AI agents working on the codebase
- Comprehensive comparison table with Terragrunt, Atmos, and plain OpenTofu
- "Why Comet?" section explaining benefits and use cases
- Architecture documentation (`docs/architecture.md`)
- Best practices guide (`docs/best-practices.md`)
- DSL improvements documentation (`docs/dsl-improvements.md`)
- DSL quick reference guide (`docs/dsl-quick-reference.md`)
- **Userland patterns guide (`docs/userland-patterns.md`)** - Comprehensive guide on building your own abstractions
- **"It's Just JavaScript!" guide (`docs/its-just-javascript.md`)** - Prominent documentation emphasizing extensibility
- Example stacks demonstrating new features and patterns
- `export` command for generating standalone Terraform files
- Integration tests for basic CLI operations
- Advanced examples in README
- Enhanced feature descriptions in README

### Changed
- Enhanced README with better feature descriptions and emojis
- **Emphasized JavaScript extensibility** throughout documentation
- Improved documentation structure
- `envs()` function now accepts both old syntax (key, value) and new object syntax for backward compatibility

### Fixed
- (List any bugs fixed in future releases)

## [0.1.0] - 2024-01-01

### Added
- Initial release
- JavaScript-based stack configuration
- Automatic backend generation
- Cross-stack references via `state()` function
- SOPS secrets integration
- Support for Terraform and OpenTofu
- CLI commands: plan, apply, destroy, list, output, clean

[Unreleased]: https://github.com/moonwalker/comet/compare/v0.6.6...HEAD
[0.6.6]: https://github.com/moonwalker/comet/releases/tag/v0.6.6
[0.6.5]: https://github.com/moonwalker/comet/releases/tag/v0.6.5
[0.6.4]: https://github.com/moonwalker/comet/releases/tag/v0.6.4
[0.6.3]: https://github.com/moonwalker/comet/releases/tag/v0.6.3
[0.6.2]: https://github.com/moonwalker/comet/releases/tag/v0.6.2
[0.6.1]: https://github.com/moonwalker/comet/releases/tag/v0.6.1
[0.6.0]: https://github.com/moonwalker/comet/releases/tag/v0.6.0
[0.5.0]: https://github.com/moonwalker/comet/releases/tag/v0.5.0
[0.1.0]: https://github.com/moonwalker/comet/releases/tag/v0.1.0
