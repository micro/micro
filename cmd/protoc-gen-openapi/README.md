protoc-gen-openapi
==================

Generate an OpenAPI 3 spec for the Micro API from your Micro service proto.


Todo
----

- [x] Read proto files
- [x] Output OpenAPI3 JSON spec
    - [x] Info.Title comes from the proto package name
- [x] Build a map of "Schemas" (one for each "message")
    - [x] Field descriptions from code
    - [ ] Possible to use references? Not for v1
- [x] Add a "Path" for each proto "service" method
    - [x] Summary is proto filename and service.method
    - [ ] Description is comments from code
    - [x] app/json RequestBody payload from message "input" schema
    - [x] app/json Response payload from message "output" schema
        - [x] Default error payloads
- [x] Add a "Server" for each Micro platform API endpoint (dev / prod etc)
- [x] Auth with API token (according to docs)
- [ ] Parameters
    - [x] Namespace (as a header)
    - [ ] Namespace (as a hostname component, but this invites errors)
    - [x] Service name (as first part of the URL path)


References aren't working (request/responses should reference a schema)


Links
-----

- https://micro.mu/reference#api
- https://docs.m3o.com/getting-started/public-apis
