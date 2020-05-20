# This is mostly intended to be triggered by CI
# as it modifies the source code.

mydir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"
cd $mydir/../platform/kubernetes/micro-resource
export TF_VAR_resource_namespace=resource

# dial down replica amount
sed -i '/replicas/ s/3/1/g' nats.tf cockroachdb.tf etcd.tf

terraform init; terraform apply -auto-approve

ssh-keygen -f /tmp/sshkey -m pkcs8 -q -N ""
ssh-keygen -f /tmp/sshkey -e  -m pkcs8 > /tmp/sshkey.pub

cd ../micro-platform
rm bot.tf

# change version to github branch
GITHUB_BRANCH=${GITHUB_REF##*/}
sed -i "/latest/ s/latest/$GITHUB_BRANCH/g" micro.tf
sed -i "/MICRO_ENABLE_ACME/ s/true/false/g" api.tf proxy.tf web.tf


export TF_VAR_platform_namespace=platform
export TF_VAR_micro_auth_private=$(cat /tmp/sshkey | base64 -w0)
export TF_VAR_micro_auth_public=$(cat /tmp/sshkey.pub | base64 -w0)
terraform init; terraform apply -auto-approve

