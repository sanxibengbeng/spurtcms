# SpurtCMS Testing Guide

This document provides guidelines for writing and running tests for the SpurtCMS project.

## Table of Contents

1. [Introduction](#introduction)
2. [Testing Structure](#testing-structure)
3. [Running Tests](#running-tests)
4. [Writing Tests](#writing-tests)
5. [Mocking Dependencies](#mocking-dependencies)
6. [Test Coverage](#test-coverage)
7. [Continuous Integration](#continuous-integration)

## Introduction

SpurtCMS uses Go's built-in testing framework for unit and integration testing. The goal is to maintain high test coverage to ensure code quality and prevent regressions.

## Testing Structure

Tests are organized alongside the code they test, following Go's convention:

```
package/
  ├── file.go
  └── file_test.go
```

We have several types of tests:

- **Unit tests**: Test individual functions and methods in isolation
- **Integration tests**: Test interactions between components
- **End-to-end tests**: Test complete workflows from the user's perspective

## Running Tests

### Running All Tests

To run all tests in the project:

```bash
go test ./...
```

### Running Tests for a Specific Package

```bash
go test ./package/...
```

### Running a Specific Test

```bash
go test -run TestFunctionName ./package/...
```

### Running Tests with Coverage

```bash
go test -cover ./...
```

For detailed coverage information:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Writing Tests

### Test File Naming

Test files should be named with the `_test.go` suffix and placed in the same package as the code they test.

### Test Function Naming

Test functions should be named `TestXxx` where `Xxx` is the name of the function or feature being tested.

### Table-Driven Tests

Use table-driven tests to test multiple scenarios with the same test logic:

```go
func TestSomething(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
        wantErr  bool
    }{
        {
            name:     "valid input",
            input:    "valid",
            expected: "result",
            wantErr:  false,
        },
        {
            name:     "invalid input",
            input:    "invalid",
            expected: "",
            wantErr:  true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := Something(tt.input)
            
            if (err != nil) != tt.wantErr {
                t.Errorf("Something() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            
            if result != tt.expected {
                t.Errorf("Something() = %v, want %v", result, tt.expected)
            }
        })
    }
}
```

### Test Helpers

Create helper functions for common test setup and teardown operations:

```go
func setupTest(t *testing.T) func() {
    // Setup code
    
    return func() {
        // Teardown code
    }
}

func TestSomething(t *testing.T) {
    cleanup := setupTest(t)
    defer cleanup()
    
    // Test code
}
```

## Mocking Dependencies

### Database Mocking

Use `github.com/DATA-DOG/go-sqlmock` to mock database interactions:

```go
func TestDatabaseFunction(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Failed to create mock: %v", err)
    }
    defer db.Close()
    
    // Set up expectations
    mock.ExpectQuery("SELECT (.+) FROM users").
        WithArgs(1).
        WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
            AddRow(1, "John"))
    
    // Call the function that uses the database
    user, err := GetUser(db, 1)
    
    // Assert results
    if err != nil {
        t.Errorf("Unexpected error: %v", err)
    }
    
    if user.Name != "John" {
        t.Errorf("Expected name 'John', got '%s'", user.Name)
    }
    
    // Verify all expectations were met
    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("Unfulfilled expectations: %v", err)
    }
}
```

### HTTP Mocking

Use `net/http/httptest` to mock HTTP requests:

```go
func TestHTTPHandler(t *testing.T) {
    // Create a request
    req, err := http.NewRequest("GET", "/path", nil)
    if err != nil {
        t.Fatal(err)
    }
    
    // Create a response recorder
    rr := httptest.NewRecorder()
    
    // Create the handler
    handler := http.HandlerFunc(YourHandler)
    
    // Serve the request
    handler.ServeHTTP(rr, req)
    
    // Check the status code
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("Handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }
    
    // Check the response body
    expected := `{"message":"success"}`
    if rr.Body.String() != expected {
        t.Errorf("Handler returned unexpected body: got %v want %v",
            rr.Body.String(), expected)
    }
}
```

### File System Mocking

Use `github.com/spf13/afero` for file system operations that need to be tested:

```go
func TestFileOperations(t *testing.T) {
    // Create a mock file system
    fs := afero.NewMemMapFs()
    
    // Write a test file
    afero.WriteFile(fs, "/test.txt", []byte("test content"), 0644)
    
    // Call the function that uses the file system
    content, err := ReadFile(fs, "/test.txt")
    
    // Assert results
    if err != nil {
        t.Errorf("Unexpected error: %v", err)
    }
    
    if string(content) != "test content" {
        t.Errorf("Expected content 'test content', got '%s'", content)
    }
}
```

## Test Coverage

Aim for at least 80% test coverage for all packages. Critical components should have higher coverage.

## Continuous Integration

Tests are automatically run in the CI pipeline for every pull request and merge to the main branch. PRs with failing tests will not be merged.

### CI Test Workflow

1. Checkout code
2. Set up Go environment
3. Install dependencies
4. Run tests with coverage
5. Upload coverage report
6. Fail if coverage is below threshold

## Best Practices

1. **Test behavior, not implementation**: Focus on what the code does, not how it does it.
2. **Keep tests simple**: Each test should verify one specific behavior.
3. **Use descriptive test names**: The test name should describe what is being tested.
4. **Test edge cases**: Include tests for boundary conditions and error cases.
5. **Avoid test interdependence**: Tests should not depend on the state from other tests.
6. **Clean up after tests**: Tests should not leave behind artifacts that could affect other tests.
7. **Use assertions judiciously**: Prefer standard Go error reporting over assertion libraries.
8. **Test public APIs**: Focus on testing the public interface of packages.
9. **Mock external dependencies**: Use mocks for databases, file systems, and external services.
10. **Keep tests fast**: Slow tests discourage frequent testing.
