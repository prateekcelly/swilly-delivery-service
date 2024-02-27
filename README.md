# Swilly Delivery Service
This service processes files containing userIDs placed by the support team in a dedicated folder and
triggers webhook api to deliver messaged to those users. 

For design decisions, refer the [ADR](docs/architecture/decisions/0002-high-volume-delivery-service.md)

### Prerequisite

**Setup GO**
- On OSX run `brew install go` (> go version 1.20)
- Make sure that the executable **go** is in your shell's path.
- Add the following in your .zshrc or .bashrc
```
> GOPATH=<workspace_dir where code will be checked out>
> export GOPATH
> PATH="${PATH}:${GOPATH}/bin"
> export PATH
```
Make sure to update the `DIRECTORY_PATH` in application.yml before starting the server/worker. This is a required config without which application won't turn up.
Make sure to create the `processed` subdirectory inside the directory path as well.

**Setup all dependencies at once using docker-compose (recommended)**
1. Install docker
2. Run on project root
```sh
$ make docker.run
```
3. To stop the container
```sh
$ make docker.stop
```

### Commands
- copy configs from sample yml to env yml `make cp-config`
- to run tests `make test`
- run the server `make start-server`
- run the worker `make start-worker`
