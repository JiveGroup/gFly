# Testing

This directory contains all test files for the application.

```bash
go test -v -coverprofile=cover.out -coverpkg=gfly/app/utils ./test/transform_utils_test.go && go tool cover -func=cover.out | grep transform
 ```

## Purpose

The test directory is used for:
- Unit tests for individual components
- Integration tests for testing component interactions
- End-to-end tests for testing complete application flows
- Test utilities and helpers

## Structure

Tests should be organized to mirror the structure of the application code they are testing:
- Unit tests for models should be in `test/models/`
- Tests for HTTP controllers should be in `test/http/controllers/`
- And so on...

## Best Practices

- Write tests for all new features and bug fixes
- Aim for high test coverage, especially for critical components
- Use descriptive test names that explain what is being tested
- Keep tests independent and isolated from each other
- Use mocks and stubs for external dependencies
- Run tests regularly during development
