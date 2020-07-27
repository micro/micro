DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"


docker run -d -p 5000:5000 --restart=always --name kind-registry -v /tmp/docker-registry:/var/lib/registry registry:2
$DIR/kind-build-micro.sh

kind create cluster --config $DIR/kind-config.yaml 
docker network connect "kind" "kind-registry"

for node in $(kind get nodes); 
do 
  kubectl annotate node "${node}" "kind.x-k8s.io/registry=localhost:5000"
done

sed_expression="s/: micro\/micro/: localhost:5000\/micro/g"
sed -e "$sed_expression" -i.bak $DIR/../platform/kubernetes/network/*.yaml

