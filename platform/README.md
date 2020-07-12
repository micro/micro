# Platform

This directory serves as platform bootstrapping for Micro.

## Overview

The platform provides **Micro as a Service** as a fully managed solution. The platform is 
bootstrapped onto Kubernetes on any major cloud provider, including load balancing and 
dns management. This repository serves as the entrypoint and single location for all bootstrapping
related source code and documentation.

## Contents

- [kubernetes](kubernetes) - include the config to deploy to k8s
- [runbook](runbook) - a directory dedicated to platform operations

## TODO

- Add config - post deployment bootstrapping config
- Add command - turn into a `micro env {create, update, delete, list}` command
- ...

