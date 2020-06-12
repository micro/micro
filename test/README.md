# Integration tests

Use these at your own risk.
It's mainly intended to run in CI and not as part of your local workflow.

The tests in this folder can be ran with `go test --tags=integration`.
It's not being triggered by `go test`.

Reasons why you should not run this locally:
* it creates a foobar directory which although is reverted in a defer, defers don't seem to work too well in tests
* it executes go gets from micro run output which might or might not modify your go.mod

## Working with these tests

Although the tests run in docker, the containers and envs are named so you can easily interact with them. Some useful tricks:

First, we have to build a local docker image:
```
bash scripts/build-local-docker.sh
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
  platform            proxy.micro.mu
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
