[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/golang-migrate/migrate/ci.yaml?branch=master)](https://github.com/golang-migrate/migrate/actions/workflows/ci.yaml?query=branch%3Amaster)
[![GoDoc](https://pkg.go.dev/badge/github.com/golang-migrate/migrate)](https://pkg.go.dev/github.com/golang-migrate/migrate/v4)
[![Coverage Status](https://img.shields.io/coveralls/github/golang-migrate/migrate/master.svg)](https://coveralls.io/github/golang-migrate/migrate?branch=master)
[![packagecloud.io](https://img.shields.io/badge/deb-packagecloud.io-844fec.svg)](https://packagecloud.io/golang-migrate/migrate?filter=debs)
[![Docker Pulls](https://img.shields.io/docker/pulls/migrate/migrate.svg)](https://hub.docker.com/r/migrate/migrate/)
![Supported Go Versions](https://img.shields.io/badge/Go-1.24%2C%201.25-lightgrey.svg)
[![GitHub Release](https://img.shields.io/github/release/golang-migrate/migrate.svg)](https://github.com/golang-migrate/migrate/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/golang-migrate/migrate/v4)](https://goreportcard.com/report/github.com/golang-migrate/migrate/v4)

# # Golang Fiber Boilerplate

This project using Go Fiber, Ent ORM, Postgres, Air, and golang-migrate for development, so make sure to install it first. This boilerplate also have auth service and controller for fast development.

## Installation

Make sure you have golang installed in your machine and set the environment based on you systems. Then run this installation in your terminal :

```bash
git clone https://github.com/tfajar123/go-boilerplate.git
cd go-boilerplate
go mod tidy
cp .env.example .env
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
go install github.com/air-verse/air@latest
```

## Usage

To run this app, you only have to run `air` in the command prompt or terminal. If you want to make a migration, we recommend that you use `make` for easy development, you can take a look `makefile` for the config. Thus you can run this prompt for migration :

```bash
# making new migrations
make migrate-create

# migrate to database
make migrate-up

# rollback migrations by 1 version
make migrate-down

# view latest migration version
make migrate-version

# force migration version
make migrate-force {{version}}

# migration reset
make migrate-reset
```

## Contributing

For major changes, please open an issue first
to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

[MIT](https://choosealicense.com/licenses/mit/)
