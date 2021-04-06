# Usage

Usage captures usage telemetry from micro

## Overview

Usage is a small net/http server backed by boltdb which captures micro usage telemetry

- Runs on port :8091
- Receives requests to `/usage` in proto format
- Backed by boltdb for stateful storage

## Usage

```
go run main.go
```
