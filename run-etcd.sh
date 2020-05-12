
#!/bin/sh

ETCD_CMD="/bin/etcd -data-dir=/data"
echo -e "Running '$ETCD_CMD'\nBEGIN ETCD OUTPUT\n"

exec $ETCD_CMD &
sleep 5
/micro $*
