set -x
sed_expression="s/: micro\/platform/: localhost:5000\/micro/g"
sed -e "$sed_expression" -i.bak cmd/platform/kubernetes/service/*.yaml
