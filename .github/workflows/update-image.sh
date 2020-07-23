set -x
branch_name=${GITHUB_REF#refs/heads/}
branch_name=${branch_name//\//-}
sed_expression="s/: micro\/micro/: localhost:5000\/micro/g"
sed -e "$sed_expression" -i.bak platform/kubernetes/network/*.yaml