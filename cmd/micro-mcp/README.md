# Micro MCP

A little experiment that exposes an MCP server to call Micro services

The server is essentially just exposing the ability to pass service/endpoint/request in a `client.Call` request.

To do it more intuitively we'd have to read the registry and register every service as a tool

Includes `call`, `describe` and `services` commands. Runs as the MCP stdio server
