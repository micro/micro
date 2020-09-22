cp go.mod go.mod.bak
go run . init --profile=ci --output=profile.go
go mod edit -replace github.com/micro/micro/profile/ci/v3=./profile/ci
go mod edit -replace google.golang.org/grpc=google.golang.org/grpc@v1.26.0
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build
docker build -t micro -f test/Dockerfile .
rm profile.go go.mod
cp go.mod.bak go.mod
rm go.mod.bak
go mod tidy
