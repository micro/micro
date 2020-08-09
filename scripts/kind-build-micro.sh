DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build
docker build -t micro -f test/Dockerfile-kind .

docker tag micro localhost:5000/micro
docker push localhost:5000/micro
