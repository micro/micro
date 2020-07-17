#!/bin/bash

# Run this script to uninstall the platform from a kubernetes cluster

# reomve micro secrets
kubectl delete secret micro-secrets

# uninstall the resources
cd ./resource/cockroachdb;
bash uninstall.sh;
cd ../etcd;
bash uninstall.sh;
cd ../nats;
bash uninstall.sh;

# delete the PVs and PVCs
for pvc in $(kubectl get pvc -o name | sed 's/persistentvolumeclaim\///g')
do 
    kubectl delete pvc $pvc
done
for pv in $(kubectl get pv -o name | sed 's/persistentvolume\///g')
do
    kubectl delete pv $pv
done

# move to the /kubernetes folder and apply the deployments
cd ../..;
kubectl delete -f network

# go back to the top level
cd ..;