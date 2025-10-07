# Changelog

All notable changes to Comet will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Comprehensive comparison table with Terragrunt, Atmos, and plain OpenTofu
- "Why Comet?" section explaining benefits and use cases
- Architecture documentation (`docs/architecture.md`)
- Best practices guide (`docs/best-practices.md`)
- `export` command for generating standalone Terraform files
- Integration tests for basic CLI operations
- Advanced examples in README
- Enhanced feature descriptions in README

### Changed
- Enhanced README with better feature descriptions and emojis
- Improved documentation structure

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
