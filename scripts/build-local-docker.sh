CGO_ENABLED=0 go build
docker build -t micro -f Dockerfile-local .
