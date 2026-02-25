# healthcheckd

Multi-Service Health Aggregator - monitors multiple service endpoints for health status.

## Purpose

Aggregate health checks across multiple services and endpoints. Generates Grafana dashboard configuration.

## Installation

```bash
go build -o healthcheckd ./cmd/healthcheckd
```

## Usage

```bash
healthcheckd <service1> <service2> ...
```

Format: `name=url[method]`

### Examples

```bash
# Health check single endpoint
healthcheckd api=http://localhost:8080/health

# Health check multiple services
healthcheckd api=http://localhost:8080/health web=http://localhost:3000

# Health check with custom method
healthcheckd service=http://localhost:9000/status[GET]
```

## Output

```
=== SERVICE HEALTH CHECK ===

api                    UP (15ms)
web                    UP (8ms)
database               DOWN (5000ms)

Summary: 2 UP, 1 DOWN

=== GRAFANA DASHBOARD CONFIG ===
{
  "dashboard": {
    "title": "Service Health Dashboard",
    "panels": [...]
  }
}
```

## Dependencies

- Go 1.21+
- github.com/fatih/color

## Build and Run

```bash
# Build
go build -o healthcheckd ./cmd/healthcheckd

# Run
go run ./cmd/healthcheckd api=http://localhost:8080/health db=http://localhost:5432/health
```

## License

MIT