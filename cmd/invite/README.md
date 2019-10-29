# Invite

Invite is a small slack invite server with google recaptcha

## Overview

Invite is a net/http server which handles recaptcha via google and redirects to a slack invite url

- Runs on port :8090 and serves `/join`
- Requires a slack invite url
- Requires a google recaptcha secret

## Usage

```
export SLACK_INVITE_URL=https://somecruft.slack/join/blabla
export GOOGLE_RECAPTCHA_SECRET=allthesecrets

go run main.go
```
