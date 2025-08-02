# Event Bus: NATS JetStream

This document describes the NATS JetStream event bus implementation for the DBETS system.

## Overview

The NATS JetStream event bus provides the following capabilities:

- **Message Replay**: Replay events from any point in time
- **Message Filtering**: Filter events based on headers or content
- **Better Durability**: Messages are persisted to disk
- **Consumer Groups**: Multiple consumers can process the same message
- **Better Scalability**: NATS is designed for high-throughput messaging

## Architecture

```
Service A → NATS JetStream → Service B
                ↓
            Persistent Storage
                ↓
            Message Replay
                ↓
            Message Filtering
```

## Configuration

### YAML Configuration

The event bus can be configured using a YAML file:

```yaml
eventbus:
  nats:
    url: "nats://localhost:4222"
    jetstream:
      stream_name: "DBETS_EVENTS"
      subjects: ["dbets.*"]
      storage: "file"
      retention: "limits"
      max_age: "24h"
      max_msgs: 1000000
      replicas: 1
```

### Environment Variables

```bash
NATS_URL=nats://localhost:4222
```

## Usage

### Basic Usage

```go
import "github.com/SilvanoP7/dbets-shared-libs/eventbus"

// Create event bus from environment
eventBus, err := eventbus.NewEventBusFromEnv()
if err != nil {
    log.Fatal(err)
}
defer eventBus.Close()

// Publish an event
event := eventbus.NewBetPlacedEvent("bet-123", "user-456", "event-789", 100.0, 2.5)
err = eventBus.Publish("bet.placed", event)

// Subscribe to events
err = eventBus.Subscribe("bet.placed", func(event interface{}) {
    // Handle the event
    log.Printf("Received bet placed event: %+v", event)
})
```

### Advanced Usage with Filtering

```go
// Subscribe with filter
filter := func(msg *nats.Msg) bool {
    userID := msg.Header.Get("User-ID")
    return userID == "specific-user-id"
}

err = eventBus.SubscribeWithFilter("bet.placed", filter, func(event interface{}) {
    // Handle filtered event
})
```

### Event Replay

```go
// Replay events from a specific time
fromTime := time.Now().Add(-1 * time.Hour)
err = eventBus.ReplayEvents("bet.placed", fromTime, func(event interface{}) {
    // Handle replayed event
})
```

## Deployment

### Docker Compose

Add NATS service to your docker-compose.yml:

```yaml
nats:
  image: nats:2.10-alpine
  container_name: dbets-nats
  command: ["-js", "-m", "8222"]
  ports:
    - "4222:4222"
    - "8222:8222"
  volumes:
    - nats_data:/data
```

### Environment Variables

Add NATS configuration to your services:

```yaml
environment:
  NATS_URL: nats://nats:4222
```

## Monitoring

### NATS Monitoring

Access the NATS monitoring interface at `http://localhost:8222` to view:

- Stream statistics
- Consumer information
- Message rates
- Storage usage

### Stream Information

```go
// Get stream statistics
info, err := eventBus.GetStreamInfo()
if err != nil {
    log.Printf("Stream info: %+v", info)
}
```

## Testing

### Run Tests

```bash
cd dbets-repos/dbets-shared-libs
go test -v ./eventbus/...
```

### Integration Test

```bash
# Start NATS
docker run -d --name nats -p 4222:4222 -p 8222:8222 nats:2.10-alpine -js -m 8222

# Run integration tests
go test -v -tags=integration ./eventbus/...
```

## Troubleshooting

### Common Issues

1. **NATS Connection Failed**
   - Check if NATS server is running
   - Verify NATS_URL environment variable
   - Check network connectivity

2. **JetStream Not Available**
   - Ensure NATS server is started with JetStream enabled
   - Check NATS server logs for JetStream initialization

3. **Message Not Delivered**
   - Check consumer configuration
   - Verify topic names match
   - Check message acknowledgment

### Debug Mode

Enable debug logging:

```go
// Set log level
log.SetLevel(log.DebugLevel)
```

## Performance Considerations

### NATS JetStream Settings

- **Storage**: Use file storage for persistence
- **Retention**: Configure based on your requirements
- **Replicas**: Use 3+ for production high availability
- **Max Messages**: Set based on memory constraints

### Consumer Settings

- **AckWait**: Set appropriate acknowledgment timeout
- **MaxDeliver**: Configure retry limits
- **FilterDurable**: Enable for filtered consumers

## Production Deployment

### High Availability

For production, use NATS cluster:

```yaml
nats:
  image: nats:2.10-alpine
  command: ["-js", "-m", "8222", "-cluster", "nats://0.0.0.0:6222"]
  ports:
    - "4222:4222"
    - "6222:6222"
    - "8222:8222"
```

### Monitoring

Set up monitoring for:

- NATS server metrics
- JetStream stream statistics
- Consumer lag
- Message throughput

### Backup

Configure JetStream backup:

```bash
# Backup stream
nats stream backup DBETS_EVENTS backup-file
```

## Event Types

The event bus supports the following event types:

- `BetPlacedEvent` - When a bet is placed
- `BetSettledEvent` - When a bet is settled
- `CashoutRequestEvent` - When a cashout is requested
- `OddsUpdatedEvent` - When odds are updated
- `EventResultEvent` - When an event result is available
- `SettlementCompleteEvent` - When settlement is complete
- `WalletUpdatedEvent` - When wallet balance changes
- `UserCreatedEvent` - When a user is created
- `UserLoginEvent` - When a user logs in
- `UserUpdatedEvent` - When a user is updated

## Best Practices

1. **Always close connections**: Use `defer eventBus.Close()`
2. **Handle errors**: Check for errors in publish/subscribe operations
3. **Use appropriate topics**: Use descriptive topic names
4. **Monitor performance**: Use the NATS dashboard to monitor usage
5. **Configure retention**: Set appropriate retention policies for your use case
6. **Test thoroughly**: Test event publishing and consumption in your environment 