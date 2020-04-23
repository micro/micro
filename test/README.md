# Warning! Dangerous integration tests

**These tests nuke your local micro store!**
Use it at your own risk.
It's mainly intended to run in CI and not as part of your local workflow.

The tests in this folder can be ran with `go test --tags=integration`.
It's not being triggered by `go test`.