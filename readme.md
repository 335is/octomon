# Octopus Monitor

 Octopus Monitor periodically executes health checks on an Octopus Deploy server.

## Requirements

### Go 1.15.7

[Download Go](https://golang.org/dl/)

## How to Use

Requires two environment variables pointing to the Octopus Deploy server.

```bash
export OCTOMON_OCTOPUS_ADDRESS=https://demo.octopus.com
export OCTOMON_OCTOPUS_APIKEY=API-GUEST
export OCTOMON_HEALTHCHECK_INTERVAL=10s
go run main.go
```

## Example Octopus Deploy Server

[Demo Octopus Deploy](https://demo.octopus.com)

username: guest

password: guest

API key: API-GUEST

## API Documentation

[Swagger UI](https://demo.octopus.com/swaggerui/index.html)

[Octopus REST API](https://octopus.com/docs/api-and-integration/api)

[Octopus Deploy API](https://github.com/OctopusDeploy/OctopusDeploy-Api/wiki)

## Example API Calls

Hit these URLs in your browser to test access to the Octopus Deploy server.

[Server Information](https://demo.octopus.com/api?apikey=API-GUEST)

[Get Projects](https://demo.octopus.com/api/projects?apikey=API-GUEST)

[Get Environments](https://demo.octopus.com/api/environments?apikey=API-GUEST)
