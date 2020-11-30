# Clients

This repo contains clients for accessing Micro

## Overview

Clients are proto generated multi-language SDKs, CLI code and related for accessing Micro via gRPC (port 8081). SDKs are generated 
automatically based on the protos and the CLI is hand maintained code for specific commands built into the micro cli.

## Clients

- API - Open API specs and docs
- CLI - Command line interface
- SDK - gRPC generated code
  * Go
  * Java
  * Node
  * Python
  * Ruby
  * Rust

## Usage

The client/sdk directory contains gRPC generated clients for each language. Point your client at the micro proxy at `localhost:8081`.

The client/cli directory contains code that enables CLI commands to be built and used which also make use of the proxy.

## Contributing

- Modify the scripts and GitHub actions to include more languages
- Code exists in scripts/generate-clients.sh and cmd/protoc-gen-client

## TODO

- [ ] Publish client sdks to the relevant package managers
- [ ] Add additional clients for different forms of access
  * api http client
  * web typescript
