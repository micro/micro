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
    - 10000: Cluster Monitoring for Kubernetes
    - 2279: NATS Server
    - 11074: Node Exporter for Prometheus
    - 11465: Cockroach SQL
    - 11463: Cockroach Runtime
    - 11466: Cockroach Replicas
    - 11464: Cockroach Storage
5) Configure a Grafana alert plugin for OpsGenie
