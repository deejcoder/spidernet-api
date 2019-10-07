# 🕸️ spidernet

Spidernet is a collection of searchable servers from across the world. It also monitors activity for all servers in the background.
This repository provides RESTful services to some client.

## Getting started with a breeze 😏
* Clone the project
* Copy `config.example.yaml` and rename the new version as `config.yaml`
* Configure `config.yaml` to your needs
* Get all of the dependencies by executing `go get ./...` in your project root directory & **RUN**!


## Roadmap
```
/api
    - api.go
    - routes.go
/handlers
    - users_*.go
    - servers_*.go
/helpers
    - auth.go
    - middleware.go
    - response.go
/storage
    - postgres.go
    - server_test.go
    - server.go
/util
    - ...
```

* `/api` is where initialization happens, and contains the routes
* `/handlers` is reserved only for handlers, handlers should be named for example, `servers_*.go`, based on their route
* `/helpers` contains mainly middleware, or things to help you out such as a wrapper to make some response for a client, you can also access the application's context here too
* `/helpers/auth` contains a simple way to restrict a route to admins only and manages authorization
* `/storage/postges.go` is the Postgres client
* `/storage/server.go` allows control over servers, use the manager to manage servers e.g `NewServerManager(db)` 

