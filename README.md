# umbrella

[![Go Reference](https://pkg.go.dev/badge/github.com/go-phings/umbrella.svg)](https://pkg.go.dev/github.com/go-phings/umbrella) [![Go Report Card](https://goreportcard.com/badge/github.com/go-phings/umbrella)](https://goreportcard.com/report/github.com/go-phings/umbrella)

Package umbrella provides a simple authentication mechanism for an HTTP endpoint. With it, you can wrap any endpoint that should have its access restricted. In addition, it provides additional its own handler for registering new user, activating it and, naturally, signing in and out.

> ⚠️ The project is in beta, under heavy construction, and may introduce breaking changes in releases before `v1.0.0`.

## Table of Contents

1. [Sample code](#sample-code)
2. [Database connection](#database-connection)
3. [User model](#user-model)
4. [Features + Roadmap](#features)
5. [Motivation](#motivation)

## Sample code
A working application can be found in the `cmd/example1`. Type `make run-example1` to start an HTTP server and check the endpoints as shown below. jq is used to parse out the token from JSON output, however, it can be done manually as well.

```bash
# run the application
% make run-example1
# ...some output...

# sign in to get a token
% UMB_TOKEN=$(curl -s -X POST -d "email=admin@example.com&password=admin" http://localhost:8001/umbrella/login | jq -r '.data.token')

# call restricted endpoint without the token
% curl http://localhost:8001/secret_stuff/ 
YouHaveToBeLoggedIn

# call it again with token
% curl -H "Authorization: Bearer $UMB_TOKEN" http://localhost:8001/secret_stuff/
SecretStuffOnlyForAdmin%

# remove temporary postgresql docker
make clean
```

## Database connection
The module needs to store users and sessions in the database. If not attached otherwise,  [struct-db-postgres](https://github.com/go-phings/struct-db-postgres) will be used as an ORM by default.

To attach a custom ORM, it needs to implement the `ORM` interface. In the `orm.go` file, there is an example on how the previously mentioned DB module is wrapped in a struct that has all the methods required by `ORM` interface.
Pass `&UmbrellaConfig` instance with `ORM` field to the `NewUmbrella` constructor to attach your object.

## User model
Umbrella comes with its own `User` and `Session` structs. However, there might be a need to use a different user model containing more fields, with a separate ORM. Hence, similarily to previous paragraph, an interface called `UserInterface` has been defined. A custom user struct must implement that interface's methods.

To do the above:
1. set `NoUserConstructor` to true in the `&UmbrellaConfig` argument when calling `NewUmbrella`
2. create new `&umbrella.Interfaces` object with `User` field and attach it to `Interfaces` field of umbrella controller.

## Features
- [X] Wrapper support for any HTTP handler  
- [X] Data storage in PostgreSQL database by default  
- [X] Customisable database driver and ORM  
- [X] Flexible User model  
- [X] Optional endpoints for sign-in (creating session objects with access tokens) and sign-out (deactivating sessions and tokens)  
- [X] Optional endpoints for user registration and activation  
- [X] Hooks triggered after successful actions like registration or sign-in  
- [X] Option to use cookies instead of the authorisation header  
- [X] Support for redirection headers after successful or failed sign-in attempts  
- [X] User struct validation during user registration  
- [X] Customisable tag names for field validation  

### Roadmap
- [ ] Simple permission system

## Motivation
While building a backend REST API for a colleague in retail, I needed a simple way to secure HTTP endpoints with basic authentication. The goal was straightforward: users would log in with an email and password, receive a token with an expiration time, and use it to interact with the backend API. A frontend application handled this flow.

A few months later, I was approached with a similar request, this time for an internal company application that required user registration and activation.

More recently, as I began developing a platform for prototyping where I used the code, I realised that this small yet essential piece of code could be valuable to others. And so, I decided to share it here.
