Monitoring Stack
================

1) Comprised of Prometheus and Grafana, installed via Helm:
    - `prometheus/install.sh`
    - `grafana/install.sh`
2) Log in to Grafana (get the "admin" password from the K8S secret, then portforward)
    - `kubectl get secret --namespace monitoring grafana -o jsonpath="{.data.admin-password}" | base64 --decode ; echo`
    - `kctl port-forward service/grafana 8080:80`
3) Configure a Prometheus datasource (server is "http://prometheus-server")
4) Import some dashboards from Grafana.com:
    - 10000 is "Cluster Monitoring for Kubernetes"
    - 2279 is "NATS Server Dashboard"
    - 11074 is "Node Exporter for Prometheus Dashboard EN"
5) Configure a Grafana alert plugin for OpsGenie
