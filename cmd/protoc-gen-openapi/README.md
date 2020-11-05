protoc-gen-openapi
==================

Generate an OpenAPI 3 spec for the Micro API from your Micro service proto.


Todo
----

- [x] Read proto files
- [x] Output OpenAPI3 JSON spec
- [ ] Build a map of "Schemas" (one for each "message")
    - [ ] Possible to use references? Not for v1
- [x] Add a "Path" for each proto "service" method
    - [x] Summary is proto filename and service.method
    - [ ] Description is comments from code
    - [ ] app/json RequestBody payload from message "input" schema
    - [ ] app/json Response payload from message "output" schema
- [ ] Add a "Server" for each Micro platform API endpoint (dev / prod etc)
- [ ] Parameters
    - [ ] Namespace (is this a path component in the Micro API?)
