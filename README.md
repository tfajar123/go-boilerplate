[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/golang-migrate/migrate/ci.yaml?branch=master)](https://github.com/golang-migrate/migrate/actions/workflows/ci.yaml?query=branch%3Amaster)
[![GoDoc](https://pkg.go.dev/badge/github.com/golang-migrate/migrate)](https://pkg.go.dev/github.com/golang-migrate/migrate/v4)
[![Coverage Status](https://img.shields.io/coveralls/github/golang-migrate/migrate/master.svg)](https://coveralls.io/github/golang-migrate/migrate?branch=master)
[![packagecloud.io](https://img.shields.io/badge/deb-packagecloud.io-844fec.svg)](https://packagecloud.io/golang-migrate/migrate?filter=debs)
[![Docker Pulls](https://img.shields.io/docker/pulls/migrate/migrate.svg)](https://hub.docker.com/r/migrate/migrate/)
![Supported Go Versions](https://img.shields.io/badge/Go-1.24%2C%201.25-lightgrey.svg)
[![GitHub Release](https://img.shields.io/github/release/golang-migrate/migrate.svg)](https://github.com/golang-migrate/migrate/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/golang-migrate/migrate/v4)](https://goreportcard.com/report/github.com/golang-migrate/migrate/v4)

# # Golang Fiber Boilerplate

This project using Go Fiber, Ent ORM, Postgres, Air, and atlasgo for development, so make sure to install it first. This boilerplate also have auth service and controller for fast development.

## Installation

Make sure you have golang installed in your machine and set the environment based on you systems. Then run this installation in your terminal :

```bash
git clone https://github.com/tfajar123/go-boilerplate.git
cd go-boilerplate
go mod download
go mod tidy
cp .env.example .env
curl -sSf https://atlasgo.sh | sh //skip this if you already have installed atlas
go install github.com/air-verse/air@latest
```

## Usage

To run this app, you only have to run `air` in the command prompt or terminal. If you want to make a migration, you must install `atlas` first, so we recommend you to use linux environment to run this app. Thus you can run this prompt for migration :

```bash
set -a
source .env
set +a

go generate ./ent
atlas migrate diff init --env local
atlas migrate apply --env local
```

You can change env based on atlas.hcl, for example `--env staging` if you want apply migration to staging, just make sure.

## Contributing

For major changes, please open an issue first
to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

[MIT](https://choosealicense.com/licenses/mit/)
