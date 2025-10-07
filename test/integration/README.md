# Integration Tests

These tests verify Comet's functionality end-to-end.

## Running Tests

```bash
cd test/integration
go test -v
```

## Test Coverage

- `basic_test.go` - Basic CLI operations (list, parse, version)
- Add more tests for plan, apply operations with mock Terraform

## Adding Tests

Create new `*_test.go` files in this directory. Tests should:

1. Build the Comet binary
2. Execute commands
3. Verify output
4. Clean up resources

## Example Test

```go
func TestNewFeature(t *testing.T) {
    // Build comet
    buildCmd := exec.Command("go", "build", "-o", "comet-test", ".")
    buildCmd.Dir = filepath.Join("..", "..")
    if err := buildCmd.Run(); err != nil {
        t.Fatalf("Failed to build comet: %v", err)
    }
    defer os.Remove(filepath.Join("..", "..", "comet-test"))

    // Run command
    cmd := exec.Command(filepath.Join("..", "..", "comet-test"), "your-command")
    cmd.Dir = filepath.Join("..", "..")
    output, err := cmd.CombinedOutput()

    // Assert
    if err != nil {
        t.Fatalf("Command failed: %v\nOutput: %s", err, output)
    }
}
```

## CI Integration

These tests can be run in CI pipelines to ensure Comet works correctly:

```yaml
# .github/workflows/test.yml
name: Test
on: [push, pull_request]

jobs:
  integration:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.23'
      - name: Run Integration Tests
        run: |
          cd test/integration
          go test -v
```
