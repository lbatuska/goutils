# Golang Utility Tools

This repository contains various utilities designed to aid Golang development. These tools are intended to improve productivity, simplify common tasks, and enhance testing and error handling.

## Included Tools (Non-Exhaustive List)

Assertions `IsNillable`, `NotNil`, `NilPtr`, `True`, `Equal`, for runtime assertion

Logging: `ConsoleLoggerImpl`, `FileLoggerImpl`, buffered logging using channels

Type: `Optional`, `Result`, as an alternative for "if err nil" error handling

Testing: `AssertEqual`, `AssertTrue`, `AssertNil`, to help quickly write tests

## Versioning

### latest:

Rolling version, pointing to the latest commit, no guarantees other than the code will compile.

### For major versions:

Anything outside the unstable directory is guaranteed to have a stable api and whatever cases and edge cases are already covered by tests are guaranteed to not change. Other behavior not covered by tests have no guarantees. (you can add any new test to add a new guarante for the future)
