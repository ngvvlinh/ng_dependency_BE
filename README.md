# eTop Backend

[![pipeline status](http://code.eyeteam.vn/etop-backend/backend/badges/master/pipeline.svg)](http://code.eyeteam.vn/etop-backend/backend/commits/master) [![coverage report](http://code.eyeteam.vn/etop-backend/backend/badges/master/coverage.svg)](http://code.eyeteam.vn/etop-backend/backend/commits/master)

## Quick Start

- Install [Go](https://golang.org/dl/)
- Install [dep](https://github.com/golang/dep)
- Put the project at `$GOPATH/src/etop.vn/backend`
- Add the following line to `/etc/hosts`:

```
127.0.0.1	postgres redis kafka
```

### Build and start

```bash
$ export PATH=$GOPATH/bin:$PATH
$ cd $GOPATH/src/etop.vn/backend
$ dep ensure
$ go install ./...

$ docker volume create --name=etop_redis_data
$ docker volume create --name=etop_postgres_data
$ docker-compose up -d
$ etop-server
```

You should see the following message:

    HTTP server listening at :8080

To verify that the server is working, execute:

    $ curl http://localhost:8080/api/admin.Misc/VersionInfo -H 'Content-Type: application/json' -d {}

It should output:

    {"service":"etop.Admin","version":"0.1"}

To view API documentation, open these URLs in browser:

- [http://localhost:8080/doc/admin](http://localhost:8080/doc/admin)
- [http://localhost:8080/doc/shop](http://localhost:8080/doc/shop)

## Development

- Install [protobuf](https://developers.google.com/protocol-buffers/docs/downloads)
- Install [jq](https://stedolan.github.io/jq/)
- Run the following script:

```
$ scripts/install-tools.sh
```

### Install dependencies

    $ scripts/install-tools.sh

### Generate protobuf

    $ scripts/protobuf-gen.sh

## Deployment

```
$ etop-server -help

Usage of etop-server:
  -config-file string
        Path to config file
  -example
        Print example config then exit
  -test
        Start services with default config for testing
```
