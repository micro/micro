DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

$DIR/test-docker.sh
docker tag micro localhost:5000/micro
docker push localhost:5000/micro
