# This is mostly intended to be triggered by CI
# as it modifies the source code.

mydir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"
cd $mydir/../platform/kubernetes/micro-resource

terraform init; terraform apply -var-file ../kind.tfvars -auto-approve
# change version to github branch
GITHUB_BRANCH=${GITHUB_REF##*/}
sed -i "/latest/ s/micro:latest/micro:$GITHUB_BRANCH/g" ../kind.tfvars
terraform init; terraform apply -auto-approve -var-file=../kind.tfvars

