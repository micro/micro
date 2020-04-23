# Warning! Dangerous integration tests

**These tests nuke your local micro store!**
Use it at your own risk.
It's mainly intended to run in CI and not as part of your local workflow.

The tests in this folder can be ran with `go test --tags=integration`.
It's not being triggered by `go test`.

Reasons why you should not run this locally:
* it nukes the micro store
* most of file manipulation commands assume Linux
* it creates a foobar directory which although is reverted in a defer, defers don't seem to work too well in tests
* it executes go gets from micro run output which might or might not modify your go.mod