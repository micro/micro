# Integration tests

Use these at your own risk.

Reasons why you should be careful with running this locally:
* it creates a foobar directory which although is reverted in a defer, defers don't seem to work too well in tests
* it executes go gets from micro run output which might or might not modify your go.mod

The tests in this folder can be ran with `go test --tags=integration`.
It's not being triggered by `go test`.

## Architecture

Key points:
- tests are being run in parallel, with different micro servers running in different containers
- local `micro run` commands will be executed with different env flags, eg. `micro -env=testConfigReadFromService run .` to connect to the above different servers.

## Working with these tests

Although the tests run in docker, the containers and envs are named so you can easily interact with them. Some useful tricks:

First, we have to build a local docker image:
```
bash scripts/test-docker.sh
```

To start a test, cd into the `test` folder and then:

```
go clean -testcache && go test --tags=integration  -failfast -v -run TestServerAuth$
```

```
$ docker ps
CONTAINER ID        IMAGE                  COMMAND                  CREATED             STATUS              PORTS                                                        NAMES
1e6a3003ea94        micro                  "sh /bin/run.sh serv…"   4 seconds ago       Up 1 second         2379-2380/tcp, 4001/tcp, 7001/tcp, 0.0.0.0:14081->8081/tcp   testServerAuth
```

As it can be seen the container name is the same as the test name.
The server output can be seen with `docker logs -f testServerAuth`.

The tests also add the env into the micro config file:

```
$ micro env
  local               none
* server              127.0.0.1:8081
  platform            proxy.m3o.com
  testServerAuth      127.0.0.1:14081
```

This means we can also interact with the server running in the container in the following way:

```
$ micro -env=testServerAuth status
```

The loop script can be used to test for flakiness:
```
cd test; bash loop.sh
```

or to run all tests once:

```
go clean -testcache && go test --tags=integration -v ./...
```

## K8s integration tests

We can run a number of integration tests against a k8s cluster rather than the default runtime implementation. We use a Kind (https://kind.sigs.k8s.io/) cluster to run a local cluster and then run the platform install scripts (with a few minor modifications).

### Running locally

#### Pre-reqs
To run the k8s integration tests locally you need to first install the pre-reqs:
- Kind, https://kind.sigs.k8s.io/
- Helm, https://helm.sh/docs/intro/install/
- cfssl, https://github.com/cloudflare/cfssl
- yq, https://github.com/mikefarah/yq

#### Running the tests
The tests can then be run:
1. `kind create cluster` - create the cluster
2. `./scripts/kind-launch.sh` - install micro in to the cluster
3. `cd tests && go clean -testcache && IN_TRAVIS_CI=yes go test --tags=integration,kind -v ./...` - run the tests

#### Adding more tests
Not all integration tests use a server so only a subset of the tests need to run against our Kind cluster. New tests should be defined in the usual way and then added to the `testFilter` slice defined near the top of [kind.go](kind.go). This is the list of all tests to be run against Kind. 

#### Running a local registry
If you prefer not having to push your images to docker hub for them to be pulled down by your Kind cluster, you can run a local registry and build and push your images to it. We have some handy scripts to get it working.
1. `./scripts/kind-local-reg.sh` - install and run a local registry, set up and launch the cluster to use it
2. `./scripts/kind-build-micro.sh` - build and push micro to the local registry
3. `./scripts/kind-launch.sh` - install micro in to the cluster

When you make any changes you can build and push using `kind-build-micro.sh` and then bounce all the micro pods `kubectl delete po -l micro=runtime` to pick up the new version.