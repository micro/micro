# This file intended to be run manually to test for flakes

while true; do
    go clean -testcache && go test --tags=integration  -failfast -v ./...;
    if [ $? -ne 0 ]; then
        spd-say "The tests are flaky";
    fi;
done
