# Platform Runbook

This directory is dedicated to the platform operations runbook.

## Overview

A [runbook](https://wa.aws.amazon.com/wat.concept.runbook.en.html) is a well documented set of procedures 
to achieve specific outcomes. Often it coincides with a playbook for incident management. Here we propose 
to establish a clear list of operations, known failure modes and success cases we expect. By doing so 
we make it easy for anyone to operate the platform.
 
### Updating Micro

Each build of micro is tagged with a snapshot, e.g. `micro/micro:20200810104423b10609`. To update the platform
to use a new tag (with zero downtime), run the following command: `kubectl set image deployments micro=micro/micro:20200810104423b10609 -l micro=runtime`. The -l flag indicates we only want to do this to deployments with
the micro=runtime label. The `micro=` part of the argument indicates the container name.