# Invite Service

Invite service is responsible for generating invite tokens for the Catalyst Experience App.

## REST API Docs

Browse the REST API docs on index: [localhost:8000](http://localhost:8000)

## Requirements

- Go 1.16 and above
- PostgreSQL 13
- Docker & Docker Compose

## Build

To build the server, execute the below command.

```console
$ make build
```

## Server configs

Envionment variables:
- `DSN` - postgres connection string

CLI flags:
- `host` - server host
- `port` - server port

## Testing

To run the tests, execute the commands below.

1. Run postgres on docker.

```console
$ make postgres
```

2. Run the tests.

```console
$ make test
```

## Deployment

1. Build the docker image.

```console
$ make build-image
```

2. Deploy via docker-compose.

```console
$ docker-compose -f invitesvc.yml up
```
