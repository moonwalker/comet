# Changelog

All notable changes to Comet will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
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

[Unreleased]: https://github.com/moonwalker/comet/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/moonwalker/comet/releases/tag/v0.1.0
