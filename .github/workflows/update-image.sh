set -x
branch_name=${GITHUB_REF#refs/heads/}
branch_name=${branch_name//\//-}
sed_expression="s/: micro\/micro/: micro\/integbuilds:$branch_name/g"
sed -e "$sed_expression" -i.bak platform/kubernetes/network/*.yaml
cat platform/kubernetes/network/proxy.yaml