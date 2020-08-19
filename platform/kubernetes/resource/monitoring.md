Monitoring Stack
================

1) Comprised of Prometheus and Grafana, installed via Helm:
    - `prometheus/install.sh`
    - `grafana/install.sh`
2) Log in to Grafana (login is in the "grafana" secret)
3) Configure a Prometheus datasource (server is "http://prometheus-server")
4) Import some dashboards from Grafana.com:
    - 10000 is "Cluster Monitoring for Kubernetes"
    - 2279 is "NATS Server Dashboard"
5) Configure a Grafana alert plugin for OpsGenie
