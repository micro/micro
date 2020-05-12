
#!/bin/sh
# Check for $CLIENT_URLS
if [ -z ${CLIENT_URLS+x} ]; then
        CLIENT_URLS="http://0.0.0.0:4001,http://0.0.0.0:2379"
        echo "Using default CLIENT_URLS ($CLIENT_URLS)"
else
        echo "Detected new CLIENT_URLS value of $CLIENT_URLS"
fi

# Check for $PEER_URLS
if [ -z ${PEER_URLS+x} ]; then
        PEER_URLS="http://0.0.0.0:7001,http://0.0.0.0:2380"
        echo "Using default PEER_URLS ($PEER_URLS)"
else
        echo "Detected new PEER_URLS value of $PEER_URLS"
fi

ETCD_CMD="/bin/etcd -data-dir=/data -listen-peer-urls=${PEER_URLS} -listen-client-urls=${CLIENT_URLS}"
echo -e "Running '$ETCD_CMD'\nBEGIN ETCD OUTPUT\n"

exec $ETCD_CMD &
sleep 5
/micro $*
