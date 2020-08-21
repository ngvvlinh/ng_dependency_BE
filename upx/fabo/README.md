# Faboshop

## Services

|Name|Description|
|---|---|
| [fabo-server](./o.o/backend/cmd/fabo-server) | Main API & Webhook Server |
| [uploader](./o.o/backend/cmd/uploader) | Upload photos |
| [event-handler](./o.o/backend/cmd/event-handler) | Handle events, send 

## Development

### Tools

- Install [Go](https://golang.org/dl/) (1.15 or above)

### Build

```
cd o.o
go install ./...
```

### Sample Config

```
fabo-server -example > fabo-server.yaml
```

## Production

```
fabo-server -config-file fabo-server.yaml
```
