# Rakuten

## Requirement
1. Docker
2. Docker Compose
3. Go 1.18 (if running outside of container)

## How to run
1. `$ cd rakuten`
2. `$ go mod vendor` to install dependency into vendor folder
3. `$ docker-compose up -d --build` to build image and spin up all docker containers

## Test
1. `$ cd rakuten`
2. `$ docker-compose up -d db-test` to spin up test db
3. `$ go test -vet=off -race -timeout=10m $( go list -e ./...)` 