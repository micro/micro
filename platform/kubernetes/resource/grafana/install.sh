#!/bin/bash

# install Grafana using Helm:
helm install stable/grafana --set persistence.enabled=true
