# DBets Shared Libraries

This repository contains shared libraries used across all DBets microservices.

## Structure

- `eventbus/` - Redis-based event bus implementation
- `models/` - Shared data models and structs

## Usage

Add this as a dependency in your service:

```go
go get github.com/dbets/shared-libs
```

Then import the packages:

```go
import (
    "github.com/dbets/shared-libs/eventbus"
    "github.com/dbets/shared-libs/models"
)
```

## Development

To update shared libraries:

1. Make changes to the appropriate package
2. Commit and push changes
3. Update the dependency in all services: `go get -u github.com/dbets/shared-libs`
