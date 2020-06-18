# eTop Backend

## Quick Start

- Install [Go](https://golang.org/dl/)
- Set the environment variable `PROJECT_DIR` to an empty directory
- Clone the project to `$PROJECT_DIR/backend`
- Add the following line to `/etc/hosts`:

```
127.0.0.1	postgres redis kafka
```

### Build and start

```bash
$ export PATH=$GOPATH/bin:$PATH
$ cd $PROJECT_DIR/backend
$ go install ./...
$ docker-compose up -d
```

Create a database with name `etop_dev`, and initialize database schema:

```
$ go run ./scripts/init_testdb -dbname etop_dev -drop
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

- Install [jq](https://stedolan.github.io/jq/)
- Run the following script:

```
$ scripts/install-tools.sh
```

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
