set -x
sed_expression="s/: micro\/micro/: localhost:5000\/micro/g"
sed -e "$sed_expression" -i.bak platform/kubernetes/network/*.yaml