# Go Todo API server

This is Todo API service to get/store Todo Items. It used in-memory stroage to store items which is written in Go.

The service provides the following features right out of the box:

* RESTful endpoints
* CI/CD using Github actions
* Deployable to kubernetes Cluster (Currently it's using DO kubernetes cluster)
* Unit tests
* CDC Test using Pact
 
The service uses the following Go packages which can be easily replaced with your own favorite ones
since their usages are mostly localized and abstracted. 

* Routing: [gorilla-mux](https://github.com/gorilla/mux)
* Logging Requests: [gorilla-handlers](https://github.com/gorilla/handlers)
* Provider Testing: [pact-go](https://github.com/pact-foundation/pact-go)

## Getting Started

If this is your first time encountering Go, please follow [the instructions](https://golang.org/doc/install) to
install Go on your computer. This service requires **Go 1.17 or above**.

[Docker](https://www.docker.com/get-started) is also needed if you want to try the kit without setting up your
own database server. The kit requires **Docker 20 or higher** for the multi-stage build support.

After installing Go and Docker, run the following commands to start todo service:

```shell
$ git clone github.com/kk-r/todo-api
$ cd todo-api
$ go mod download
// Build server binaries
$ go build -o go-api
$ ./go-api
```

At this time, you have a API server running at `http://127.0.0.1:8080`. It provides the following endpoints:

* `GET /health`: a healthcheck service provided for health checking purpose (needed when implementing a server cluster)
* `POST /api/todos`: Creates new Todo item
* `GET /api/todos`: returns a list of Todo items

Try the URL `http://localhost:8080/health` in a browser, and you should see something like `"{"alive": true}"` displayed.

If you have `cURL` or some API client tools (e.g. [Postman](https://www.getpostman.com/)), you may try the following 
more complex scenarios:

```shell
# Create new todo item via: POST /api/todos
curl -X POST -H "Content-Type: application/json" -d '{ "title": "test-gh"}' http://localhost:8080/api/todos
# should return list of all todos: {"0":{"title":"test-gh"}}

# Get List of All todos, such as: GET /api/todos
curl -X GET http://localhost:8080/api/todos
# should return a list of todos in the JSON format
```

Since I'm using in-memory storage to keep the lists, it's *not recommend* to run the appliccation multiple clusters. It will have some inconsitency between the records.


## Project Layout

The starter kit uses the following project layout:
 
```
.
├── .github/workflows    Gihub action workflow
│   └── go.yaml          CI/CD config file
├── config               Kubernetes configuration files for different environments
```

## Common Development Tasks

### Run Tests

There are two types of tests covered in this project Unit test , [CDC Provider Test migration](https://pactflow.io/blog/the-curious-case-for-the-provider-driven-contract/). We used pact.io to test the contract between Api service and Frontend. 
Before run provider test please make sure install the pact binaries by using this [link](https://github.com/pact-foundation/pact-go#installation-on-nix).  The following commands are commonly used for tests:

```shell
# Run unit tests.
CI=true go test -v ./...

# Execute all tests in the project folder.
# Usually you should run this command for before deployment. 
PACT_BROKER_URL=pact_broker_url PACT_BROKER_TOKEN=pact_token go test -v ./...
```

## Deployment

The application can be run as a docker container. You can use `Dockerfile` to build the application 
into a docker image. The docker container needs few environment variables to run, please check the *Dockerfile* and pass the build params as your config.

```shell
docker build --build-arg PACT_BROKER_URL=https://localhost:8888/ --build-arg PACT_BROKER_TOKEN=token  -t  todo-api:v1 --no-cache .
```

For CI/CD we use Github actions to deploy the applications in the kubernetes cluster. Currently it's configured to deploy the cluster in DigitalOcean. If you are planning to use similar pipeline please be make sure to set secrects value based on your requirement.


## Todo

- [x] Cover pact proider test.
- [ ] Update docker-compose file for local development.
- [ ] Create new Dockerfile for running tests.
- [ ] Update Build pipeline to run pact test.
- [ ] Create detailed doc about architectural decisions in the README
- [ ] Make docs with GoDoc
- [x] Building a deployment process to deploy staging enviornment.
- [] Run tests on staging enviornment as part of deployment process.


## License

[MIT License](./LICENSE.md)