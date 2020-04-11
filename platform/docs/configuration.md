# Configuration

This document serves as the place for extended configuration

## Overview

The platform is automated through terraform and requires certain environmental config before it 
can be used, including the configuration for the underlying services and their access to resources 
like github and cloudflare.

## Dependencies

- Terraform
- Github
- ...

## Environment

A few things we need

- FRONTEND_ADDRESS - the URL the dashboard is served on
- GITHUB_TEAM_ID - The team which has access
- GITHUB_OAUTH_CLIENT_ID - github oauth client id
- GITHUB_OAUTH_CLIENT_SECRET - github oauth client secret
- GITHUB_OAUTH_REDIRECT_URL - github oauth redirect url
- MICRO_AUTH - the type of auth, e.g. jwt
- MICRO_AUTH_PUBLIC_KEY - the base64 encoded public jwt key
- MICRO_AUTH_PRIVATE_KEY - the base64 encoded private jwt key

Image Pull Credentials: The default serviceaccount needs "Image pull secrets" set to a GitHub token.

## Usage

Coming soon...
