set -x
sed_expression="s/: micro\/platform/: localhost:5000\/micro/g"
sed -e "$sed_expression" -i.bak platform/kubernetes/service/*.yaml
