# Octopus Monitor

 Octopus Monitor periodically executes health checks on an Octopus Deploy server.

## Requirements

### Go 1.11.5

[Download Go](https://golang.org/dl/)

### Dep 0.4.1

[Dep Releases](https://github.com/golang/dep/releases/)

## How to Use

Requires two environment variables pointing to the Octopus Deploy server.

```bash
export OCTOMON_OCTOPUS_ADDRESS=https://demo.octopusdeploy.com
export OCTOMON_OCTOPUS_APIKEY=API-GUEST
go run main.go
```

## Example Octopus Deploy Server

[Demo Octopus Deploy](https://demo.octopusdeploy.com)
username: guest
password: guest
API key: API-GUEST

## API Documentation

[Swagger UI](http://demo.octopusdeploy.com/swaggerui/index.html)

[Octopus Deploy API](https://octopus.com/docs/api-and-integration/api)

## Example API Calls

Hit these URLs in your browser to test access to the Octopus Deploy server.

[Server Information](https://demo.octopusdeploy.com/api?apikey=API-GUEST)

[Get Projects](https://demo.octopusdeploy.com/api/projects?apikey=API-GUEST)

[Get Environments](https://demo.octopusdeploy.com/api/environments?apikey=API-GUEST)