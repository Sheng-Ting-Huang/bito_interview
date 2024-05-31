
## Project Layout

### Packages:
- `api` : consists of the HTTP router and API handlers
- `model` : core models such as person and his/her attributes.
- `storage` : storing personal information and executing the query for the matching. 

`main.go` is the entry point for the program and the http server is listening to 8080 port.

`go.mod` and `go.sum` are the dependency packages from other third party libraries.

`dockerfile` is for building the docker image.

```
├── api
│   ├── api.go
│   ├── api_test.go
│   ├── dto.go
│   └── init.go
├── dockerfile
├── go.mod
├── go.sum
├── main.go
├── model
│   └── core.go
└── storage
    ├── access.go
    ├── access_test.go
    ├── idGenerator.go
    └── init.go
```

