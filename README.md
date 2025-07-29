# DBets Shared Libraries

This repository contains shared Go packages used across all DBets microservices.

## 📦 Packages

### Event Bus (`eventbus/`)
- **RedisEventBus**: Redis-based event bus implementation
- **EventBus Interface**: Common interface for event bus implementations
- **Event Types**: Standardized event structures for inter-service communication

### Models (`models/`)
- **User Models**: User authentication and profile data structures
- **Betting Models**: Bet, event, market, and selection data structures
- **Wallet Models**: Financial transaction and balance data structures
- **API Models**: Request/response structures for all services

### Utils (`utils/`)
- **Common Utilities**: Shared utility functions
- **Validation**: Common validation logic
- **Helpers**: Reusable helper functions

## 🚀 Usage

### Adding to a Service

1. Add the dependency to your `go.mod`:
```go
require (
    github.com/SilvanoP7/shared-libs v0.0.0-00010101000000-000000000000
    // ... other dependencies
)

replace github.com/SilvanoP7/shared-libs => ../shared-libs
```

2. Import the packages:
```go
import (
    "github.com/SilvanoP7/shared-libs/eventbus"
    "github.com/SilvanoP7/shared-libs/models"
)
```

### Event Bus Usage

```go
// Initialize event bus
eventBus, err := eventbus.NewRedisEventBus("redis://localhost:6379")
if err != nil {
    log.Fatal(err)
}

// Publish an event
err = eventBus.Publish("bet.placed", &models.BetPlacedEvent{
    BetID:   "bet-123",
    UserID:  "user-456",
    Amount:  100.0,
    EventID: "event-789",
})
```

### Models Usage

```go
// Use shared models
user := &models.User{
    ID:       "user-123",
    Username: "john_doe",
    Email:    "john@example.com",
}

bet := &models.Bet{
    ID:         "bet-123",
    UserID:     "user-456",
    EventID:    "event-789",
    Amount:     100.0,
    Odds:       2.5,
    Status:     "pending",
}
```

## 🔧 Development

### Local Development Setup

1. Clone this repository alongside your service repositories:
```bash
mkdir ~/dbets-dev
cd ~/dbets-dev
git clone https://github.com/SilvanoP7/shared-libs.git
```

2. In each service repository, add the replace directive:
```bash
echo "replace github.com/SilvanoP7/shared-libs => ../shared-libs" >> go.mod
go mod tidy
```

### Adding New Shared Code

1. Create new packages in the appropriate directory
2. Update this README with documentation
3. Test with all dependent services
4. Commit and push changes

### Versioning

This library uses semantic versioning. When making breaking changes:

1. Update the version in `go.mod`
2. Update all dependent services to use the new version
3. Document breaking changes in the changelog

## 📋 Dependencies

- `github.com/go-redis/redis/v8` - Redis client for event bus
- `github.com/google/uuid` - UUID generation
- `gorm.io/gorm` - Database ORM

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test with dependent services
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License.
