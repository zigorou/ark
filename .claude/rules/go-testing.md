---
paths:
  - "**/*_test.go"
---

# Go Testing Conventions

## Standards

- Use **table-driven tests** with `t.Run(name, ...)` for multiple cases.
- Test file: `foo_test.go` alongside `foo.go` in the same package.
- External test package (`package foo_test`) for black-box testing exported APIs.
- Use `t.Helper()` in all test helper functions.
- Use `t.Parallel()` where tests are independent.
- Run with `-race` flag to detect data races.
- No global mutable state in tests; use `t.TempDir()` for temp files.

## Coverage Goals

- Core logic packages (`internal/*`): aim for ≥ 80% coverage.
- CLI command packages (`cmd/*`): integration-style tests preferred over unit tests.

## Test Naming

```go
func TestFunctionName(t *testing.T) { ... }         // basic
func TestFunctionName_scenario(t *testing.T) { ... } // scenario variant
```
