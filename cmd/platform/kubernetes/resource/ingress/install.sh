#!/bin/bash
# Reference: https://cert-manager.io/docs/tutorials/acme/ingress/

# REQUIRED MICRO ENV CF_API_KEY

# Install nginx using helm
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm install nginx ingress-nginx/ingress-nginx
kubectl wait deployment/nginx-ingress-nginx-controller --for=condition=available --timeout=120s

# Replace m3o.com with m3o.dev in staging
if [ $MICRO_ENV == "staging" ]; then
  sed -i '' 's/\*.m3o.app/\*.m3o.dev/g' ingress.yaml
  sed -i '' 's/m3o.com/m3o.dev/g' ingress.yaml
fi

# Install the ingress
kubectl apply -f ingress.yaml

# replace back
if [ $MICRO_ENV == "staging" ]; then
  sed -i '' 's/\*.m3o.dev/\*.m3o.app/g' ingress.yaml
  sed -i '' 's/m3o.dev/m3o.com/g' ingress.yaml
fi

# Don't use TLS locally
if [ "$MICRO_ENV" == "dev" ]; then
  exit 0
fi

# Install Cert Manager
kubectl create namespace cert-manager
helm repo add jetstack https://charts.jetstack.io
helm install \
  cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --version v1.0.3 \
  --set installCRDs=true
kubectl wait deployment --all --for=condition=available -n cert-manager --timeout=120s

echo "Waiting for a Public IP to be assigned to the nginx ingress..."
while true; do
  grpcIP=$(kubectl get ingress grpc-ingress -o jsonpath="{.status.loadBalancer.ingress[0].ip}")
  httpIP=$(kubectl get ingress http-ingress -o jsonpath="{.status.loadBalancer.ingress[0].ip}")
  if [ "$grpcIP" == "" ] | [ "$httpIP" == "" ]; then
    sleep 1
  else 
    break
  fi
done


echo "Please set the following DNS entries and then press [y] to continue:"
kubectl get ingress grpc-ingress -o jsonpath="{.spec.rules[*].host}" | xargs -n 1 -I{} printf "%-10s %-30s %-20s\n" "A" {} $grpcIP
kubectl get ingress http-ingress -o jsonpath="{.spec.rules[*].host}" | xargs -n 1 -I{} printf "%-10s %-30s %-20s\n" "A" {} $httpIP 

while true; do
  read -r ans
  if [ "$ans" == "y" ]; then
    break
  else 
    echo "Invalid input, please press [y] to continue"
  fi
done

# Update the ingress to use letsencrypt
kubectl apply -f ./letsencrypt.yaml
kubectl annotate ingress grpc-ingress cert-manager.io/issuer="letsencrypt-prod" --overwrite
kubectl annotate ingress http-ingress cert-manager.io/issuer="letsencrypt-prod" --overwrite

echo "nginx ingress configured, it will take about 2-3 minutes for the TLS certificate to be issued"