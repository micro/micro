DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
# safety first
kubectl config use-context kind-kind
pushd $DIR/../cmd/platform/kubernetes
./uninstall.sh
popd
# delete all the namespaces we've added
namespaces=$(kubectl get namespaces -o name | sed 's/namespace\///g')
for ns in $namespaces 
do
    if [[ $ns == "kube-system" || $ns == "kube-node-lease" || $ns == "default" || $ns == "kube-public" || $ns == "local-path-storage" || $ns == "default" ]]; then
        continue
    fi
    kubectl delete namespace $ns
done

$DIR/./kind-launch.sh
