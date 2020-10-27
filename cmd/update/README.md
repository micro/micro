# Update

Update is a small go http server and script which maintains a local server upto date.

## Overview

- main.go - is a small net/http server which handles updates
  * Runs on port :9000
  * Returns latest commit, release and image at `GET /update`
  * Processes webhook updates at `POST /update`

## Usage

Setup webhooks in dockerhub and in github to point to your /update endpoint

Additionally set a github secret and pass via env var

```
GITHUB_WEBHOOK_SECRET=foobar go run main.go
```
