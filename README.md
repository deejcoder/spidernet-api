# 🕸️ spidernet

Spidernet is a collection of searchable servers from across the world. It also monitors activity for all servers in the background.


## Getting started with a breeze 😏
* Clone the project
* Copy `config.example.yaml` and rename the new version as `config.yaml`
* Configure `config.yaml` to your needs
* Get all of the dependencies by executing `go get ./...` in your project root directory & **RUN**!

## Current project plan
### Endpoints
* GET /servers returns a server list/index, limited to 10 servers at a time, implement pagination.
* GET /servers?tags={list of tags} allows anyone to search a server by tag(s)
* GET /servers?query={search-string} allows anyone to search a server by search string (e.g IP)
* GET /servers/{id} returns a single server with extra information
* GET /servers/{id}/activity returns all activity associated with the server (downtime)
* PUT /servers/{id}/rating allows a user to submit a up/down vote for this server
* UPDATE /servers/{id} allows the person who originally submitted the server, to update it, or administrators.
* DELETE /servers/{id} allows the authorized to delete a server.
* POST /servers allows users to add servers, the server must be available (pingable) before it is added.
* POST /servers/{id}/comment allows a user to post a comment on this server
* UPDATE /users/{id} allows a user or administrator to edit a user account
* POST /register allows someone to create a new user account
* GET /users/{id} returns information about a user, including a list of servers they have added
* POST /login allows someone to login to their user account

### Project Structure
```
/api
    /helpers
        - auth.go
        - middleware.go
        - response.go
    /handlers
        - servers.go
        - auth.go
        - users.go
/storage
    - client.go
    - users.go
    /server
        - server.go
        - ratings.go
        - comments.go
        - activity.go
/tests
/workers
    - pinger.go
```

* **API** will contain RESTful services provided to the user through the HTTP protocol.
* **Storage** will provide database services to the **API**  and **workers**
* **Tests** will be the place where I perform tests
* **Workers** will provide workers which will be executed in different threads, aside of the **API**. These workers may in the future be distributed using Redis to multiple processes within the same host, or across multiple hosts. In particular, this is where we will keep track of servers statuses and activity.


