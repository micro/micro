# Explore

The explorer is a http server which aggregates all the github projects related to micro.

## Overview

The explorer aggregates github projects that make use of micro and stores them in elasticsearch making them queryable.

- Serves a net/http server on port :8089 and `/explore`
- Depends on a local elasticsearch instance for storage
- Requires github access token

## Usage

- Start elasticsearch
- Start explorer with github access token 

```
GITHUB_ACCESS_TOKEN=gxrykfsfpsfpsfpsf go run main.go
```
